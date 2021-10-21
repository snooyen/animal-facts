package main

import (
	"strconv"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/kustomize"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	bitnamiHelmRepo = `https://charts.bitnami.com/bitnami`
	redisChart      = `redis`
)

var (
	services = []string{
		"facts-api",
		"fact-scraper",
		"fact-publisher",
		"fact-admin",
	}
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Kubernetes Configs
		k8config := config.New(ctx, "kubernetes")
		namespace := k8config.Require("namespace")

		// Create Namespace
		_, err := corev1.NewNamespace(ctx, namespace, &corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaArgs{Name: pulumi.String(namespace)},
		})
		if err != nil {
			return err
		}

		// Deploy Redis
		redisConfig := config.New(ctx, "redis")
		redisPassword := string(redisConfig.Require("password"))
		replicas, err := strconv.Atoi(redisConfig.Require("replicas")) // FIXME: Unhandled error
		if err != nil {
			return err
		}

		err = DeployRedisChart(ctx, RedisValues{
			namespace: namespace,
			password:  redisPassword,
			replicas:  replicas,
		})
		if err != nil {
			return err
		}

		for _, service := range services {
			err = DeployServiceOverlay(ctx, service)
			if err != nil {
				return err
			}
		}

		// Create Twilio Auth Secret
		twilioConfig := config.New(ctx, "twilio")
		twilioAccountSID := string(twilioConfig.Require("accountsid"))
		twilioToken := string(twilioConfig.Require("token"))
		_, err = corev1.NewSecret(ctx, "twilio", &corev1.SecretArgs{
			Type: pulumi.String("opaque"),
			Metadata: &metav1.ObjectMetaArgs{
				Name:      pulumi.String("twilio"),
				Namespace: pulumi.String(namespace),
			},
			StringData: pulumi.StringMap{
				"accountsid": pulumi.String(twilioAccountSID),
				"token":      pulumi.String(twilioToken),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}

type RedisValues struct {
	namespace string
	password  string
	replicas  int
}

func DeployRedisChart(ctx *pulumi.Context, values RedisValues) error {
	_, err := helm.NewRelease(ctx, redisChart, &helm.ReleaseArgs{
		Chart: pulumi.String(redisChart),
		RepositoryOpts: helm.RepositoryOptsArgs{
			Repo: pulumi.String(bitnamiHelmRepo),
		},
		Name:            pulumi.String(redisChart),
		Namespace:       pulumi.String(values.namespace),
		CreateNamespace: pulumi.Bool(true),
		Values: pulumi.Map{
			"auth": pulumi.Map{
				"password": pulumi.String(values.password),
			},
			"replica": pulumi.Map{
				"replicaCount": pulumi.Int(values.replicas),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func DeployServiceOverlay(ctx *pulumi.Context, service string) (err error) {
	config := config.New(ctx, service)
	overlay := config.Require("overlay")
	_, err = kustomize.NewDirectory(ctx, service,
		kustomize.DirectoryArgs{
			Directory: pulumi.String(overlay),
		},
	)
	return
}

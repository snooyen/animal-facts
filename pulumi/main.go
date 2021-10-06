package main

import (
	"strconv"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	bitnamiHelmRepo = `https://charts.bitnami.com/bitnami`
	redisChart      = `redis`
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

		return nil
	})
}

type RedisValues struct {
	namespace string
	password  string
	replicas  int
}

func DeployRedisChart(ctx *pulumi.Context, values RedisValues) error {
	_, err := helm.NewChart(ctx, redisChart, helm.ChartArgs{
		Chart:     pulumi.String(redisChart),
		Namespace: pulumi.String(values.namespace),
		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String(bitnamiHelmRepo),
		},
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

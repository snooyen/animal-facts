package main

import (
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	bitnamiHelmRepo = `https://charts.bitnami.com/bitnami`
	redisChart      = `redis`
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Stack Configs
		k8config := config.New(ctx, "kubernetes")
		namespace := k8config.Require("namespace")

		err := DeployRedisChart(ctx, RedisValues{namespace: namespace})
		if err != nil {
			return err
		}

		return nil
	})
}

type RedisValues struct {
	namespace string
}

func DeployRedisChart(ctx *pulumi.Context, values RedisValues) error {
	_, err := helm.NewChart(ctx, redisChart, helm.ChartArgs{
		Chart:     pulumi.String(redisChart),
		Namespace: pulumi.String(values.namespace),
		FetchArgs: helm.FetchArgs{
			Repo: pulumi.String(bitnamiHelmRepo),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

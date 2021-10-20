package openapi

import "context"

type Client interface {
	ListApps(ctx context.Context, appIds []string) ([]AppInfo, error)

	ListEnvClusters(ctx context.Context, appId string) ([]EnvClusters, error)
	CreateCluster(ctx context.Context, name, appId string) (*ClusterInfo, error)

	ListNamespaces(ctx context.Context, appId, env, clusterName string) ([]NamespaceInfo, error)
	CreateNamespace(
		ctx context.Context, name, appId string, format Format, isPublic bool, comment string) (*NamespaceInfo, error)
	GetNamespaceItem(ctx context.Context, key, value string)
	CreateNamespaceItem(ctx context.Context, key, value string)
	UpdateNamespaceItem(ctx context.Context, key, value string)
	DeleteNamespaceItem(ctx context.Context)
	PublishNamespace(ctx context.Context)
}

type Config struct {
	Token         string
	PortalAddress string
	// Account userId in sso system.
	Account string
}

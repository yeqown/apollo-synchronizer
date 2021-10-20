package openapi

import "context"

type Client interface {
	ListApps(ctx context.Context, appIds []string) ([]AppInfo, error)

	ListEnvClusters(ctx context.Context, appId string) ([]EnvClusters, error)
	CreateCluster(ctx context.Context, name, appId string) (*ClusterInfo, error)

	ListNamespaces(ctx context.Context, appId, env, cluster string) ([]NamespaceInfo, error)
	CreateNamespace(
		ctx context.Context, name, appId string, format Format, isPublic bool, comment string) (*NamespaceInfo, error)
	PublishNamespace(ctx context.Context, appId, env, cluster, namespace string) (*PublishNamespaceResult, error)
	DeleteNamespace(ctx context.Context, appId, env, cluster, namespace string) error

	GetNamespaceItem(ctx context.Context, appId, env, cluster, namespace, key string) (*NamespaceItem, error)
	CreateNamespaceItem(ctx context.Context, appId, env, cluster, namespace, key, value string) (*NamespaceItem, error)
	UpdateNamespaceItem(ctx context.Context, appId, env, cluster, namespace, key, value string) (*NamespaceItem, error)
	DeleteNamespaceItem(ctx context.Context, appId, env, cluster, namespace, key string) error
}

type Config struct {
	Token         string
	PortalAddress string
	// Account userId in sso system.
	Account string
}

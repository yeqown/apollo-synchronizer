package openapi

import (
	"context"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var (
	_ Client = new(openapiClient)
)

type openapiClient struct {
	config *Config
	cc     *resty.Client
}

func New(c *Config) Client {
	client := resty.
		New().
		SetHeader("Authorization", c.Token).
		SetHeader("Content-Type", "application/json;charset=UTF-8")

	return &openapiClient{
		config: c,
		cc:     client,
	}
}

func (o openapiClient) ListApps(ctx context.Context, appIds []string) ([]AppInfo, error) {
	// expand uri
	uri := expand("${portalAddress}/openapi/v1/apps", map[string]string{
		"portalAddress": o.config.PortalAddress,
	})

	result := make([]AppInfo, 0, 8)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetQueryParam("appIds", strings.Join(appIds, ",")).
		SetResult(&result).
		Get(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListApps in openapiClient.ListApps")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) ListEnvClusters(ctx context.Context, appId string) ([]EnvClusters, error) {
	panic("implement me")
}

func (o openapiClient) CreateCluster(ctx context.Context, name, appId string) (*ClusterInfo, error) {
	panic("implement me")
}

func (o openapiClient) ListNamespaces(ctx context.Context, appId, env, clusterName string) ([]NamespaceInfo, error) {
	uri := expand(
		"${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/namespaces",
		map[string]string{
			"portalAddress": o.config.PortalAddress,
			"env":           env,
			"appId":         appId,
			"clusterName":   clusterName,
		},
	)

	result := make([]NamespaceInfo, 0, 8)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListNamespaces in openapiClient.ListNamespaces")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) CreateNamespace(
	ctx context.Context, name, appId string, format Format, isPublic bool, comment string) (*NamespaceInfo, error) {
	uri := expand("${portalAddress}/openapi/v1/apps/${appId}/appnamespaces", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"appId":         appId,
	})

	body := struct {
		Name                string `json:"name"`
		AppID               string `json:"appId"`
		Format              Format `json:"format"`
		IsPublic            bool   `json:"isPublic"`
		Comment             string `json:"comment"`
		DataChangeCreatedBy string `json:"dataChangeCreatedBy"`
	}{
		Name:                name,
		AppID:               appId,
		Format:              format,
		IsPublic:            isPublic,
		Comment:             comment,
		DataChangeCreatedBy: o.config.Account,
	}

	result := new(NamespaceInfo)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		Post(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListNamespaces in openapiClient.ListNamespaces")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) PublishNamespace(
	ctx context.Context, appId, env, cluster, namespace string) (*PublishNamespaceResult, error) {
	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/"+
		"namespaces/${namespaceName}/releases", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
	})

	body := struct {
		ReleaseTitle   string `json:"releaseTitle"`
		ReleaseComment string `json:"releaseComment"`
		ReleasedBy     string `json:"releasedBy"`
	}{
		ReleaseTitle:   time.Now().Format(time.RFC3339) + "-release",
		ReleaseComment: "apollo-synchronizer automatically publish",
		ReleasedBy:     o.config.Account,
	}

	result := new(PublishNamespaceResult)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		Post(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to PublishNamespace in openapiClient.PublishNamespace")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil

}

func (o openapiClient) DeleteNamespace(ctx context.Context, appId, env, cluster, namespace string) error {
	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/"+
		"namespaces/${namespaceName}", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
	})

	r, err := o.cc.
		R().
		SetContext(ctx).
		Delete(uri)
	if err != nil {
		return errors.Wrap(err, "failed to DeleteNamespace in openapiClient.DeleteNamespace")
	}

	if err = handleResponseError(r); err != nil {
		return err
	}

	return nil
}

func (o openapiClient) GetNamespaceItem(
	ctx context.Context, appId, env, cluster, namespace, key string) (*NamespaceItem, error) {

	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/"+
		"${clusterName}/namespaces/${namespaceName}/items/${key}", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
		"key":           key,
	})

	result := new(NamespaceItem)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetResult(&result).
		Get(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetNamespaceItem in openapiClient.GetNamespaceItem")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) CreateNamespaceItem(
	ctx context.Context, appId, env, cluster, namespace, key, value string) (*NamespaceItem, error) {

	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/"+
		"namespaces/${namespaceName}/items", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
	})

	body := struct {
		Key                 string `json:"key"`
		Value               string `json:"value"`
		Comment             string `json:"comment"`
		DataChangeCreatedBy string `json:"dataChangeCreatedBy"`
	}{
		Key:                 key,
		Value:               value,
		Comment:             "apollo-synchronizer auto created",
		DataChangeCreatedBy: o.config.Account,
	}

	result := new(NamespaceItem)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetBody(body).
		SetResult(&result).
		Post(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to GetNamespaceItem in openapiClient.GetNamespaceItem")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateNamespaceItem
// https://github.com/apolloconfig/apollo/wiki/Apollo%E5%BC%80%E6%94%BE%E5%B9%B3%E5%8F%B0
// #3211-%E4%BF%AE%E6%94%B9%E9%85%8D%E7%BD%AE%E6%8E%A5%E5%8F%A3
func (o openapiClient) UpdateNamespaceItem(
	ctx context.Context, appId, env, cluster, namespace, key, value string) (*NamespaceItem, error) {

	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/"+
		"namespaces/${namespaceName}/items/${key}", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
		"key":           key,
	})

	body := struct {
		Key                      string `json:"key"`
		Value                    string `json:"value"`
		Comment                  string `json:"comment"`
		DataChangeLastModifiedBy string `json:"dataChangeLastModifiedBy"`
		DataChangeCreatedBy      string `json:"dataChangeCreatedBy"`
	}{
		Key:                      key,
		Value:                    value,
		Comment:                  "apollo-synchronizer auto modified",
		DataChangeLastModifiedBy: o.config.Account,
		DataChangeCreatedBy:      o.config.Account,
	}

	result := new(NamespaceItem)
	r, err := o.cc.
		R().
		SetContext(ctx).
		SetQueryParam("createIfNotExists", "true").
		SetBody(body).
		SetResult(&result).
		Put(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to UpdateNamespaceItem in openapiClient.UpdateNamespaceItem")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) DeleteNamespaceItem(ctx context.Context, appId, env, cluster, namespace, key string) error {
	uri := expand("${portalAddress}/openapi/v1/envs/${env}/apps/${appId}/clusters/${clusterName}/"+
		"namespaces/${namespaceName}/items/${key}?operator=${operator}", map[string]string{
		"portalAddress": o.config.PortalAddress,
		"env":           env,
		"appId":         appId,
		"clusterName":   cluster,
		"namespaceName": namespace,
		"key":           key,
		"operator":      o.config.Account,
	})

	r, err := o.cc.
		R().
		SetContext(ctx).
		Delete(uri)
	if err != nil {
		return errors.Wrap(err, "failed to DeleteNamespaceItem in openapiClient.DeleteNamespaceItem")
	}

	if err = handleResponseError(r); err != nil {
		return err
	}

	return nil
}

package openapi

import (
	"context"
	"strings"

	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
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
	uri := expand("${portalAddress}/openapi/v1/apps/{appId}/appnamespaces", map[string]string{
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
		Get(uri)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ListNamespaces in openapiClient.ListNamespaces")
	}

	if err = handleResponseError(r); err != nil {
		return nil, err
	}

	return result, nil
}

func (o openapiClient) GetNamespaceItem(ctx context.Context, key, value string) {
	panic("implement me")
}

func (o openapiClient) CreateNamespaceItem(ctx context.Context, key, value string) {
	panic("implement me")
}

func (o openapiClient) UpdateNamespaceItem(ctx context.Context, key, value string) {
	panic("implement me")
}

func (o openapiClient) DeleteNamespaceItem(ctx context.Context) {
	panic("implement me")
}

func (o openapiClient) PublishNamespace(ctx context.Context) {
	panic("implement me")
}

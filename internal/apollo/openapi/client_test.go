package openapi_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yeqown/apollo-synchronizer/internal/apollo/openapi"
)

type openapiClientTestSuite struct {
	suite.Suite

	apollo openapi.Client
}

func (o *openapiClientTestSuite) SetupSuite() {
	o.apollo = openapi.New(&openapi.Config{
		Token:         "82a95a5722ae8649f64ca5859a13032acab4b2a3",
		PortalAddress: "http://localhost:8070",
		Account:       "apollo",
	})
}

func (o openapiClientTestSuite) Test_ListApps() {
	appInfos, err := o.apollo.ListApps(context.Background(), []string{"demo"})
	o.NoError(err)
	o.NotEmpty(appInfos)

	o.T().Logf("appInfos: %+v \n", appInfos)
}

func (o openapiClientTestSuite) Test_ListNamespaces() {
	namespaceInfos, err := o.apollo.ListNamespaces(context.Background(), "demo", "DEV", "default")
	o.NoError(err)
	o.NotEmpty(namespaceInfos)

	o.T().Logf("namespaceInfos: %+v \n", namespaceInfos)
}

func (o openapiClientTestSuite) Test_CreateNamespace() {
	namespace, err := o.apollo.CreateNamespace(context.Background(),
		"Test_CreateNamespaces", "demo", openapi.Format_JSON, false, "this is comment")
	o.NoError(err)
	o.NotEmpty(namespace)

	o.T().Logf("namespace: %+v \n", namespace)
}

func (o openapiClientTestSuite) Test_DeleteNamespace() {
	err := o.apollo.DeleteNamespace(context.Background(),
		"demo", "DEV", "default", "application1.yaml")
	o.NoError(err)
}

func (o openapiClientTestSuite) Test_PublishNamespace() {
	publishResult, err := o.apollo.PublishNamespace(context.Background(),
		"demo", "DEV", "default", "application1.yaml")
	o.NoError(err)
	o.NotEmpty(publishResult)

	o.T().Logf("namespace: %+v \n", publishResult)
}

func (o openapiClientTestSuite) Test_DeleteNamespaceConfig() {
	err := o.apollo.DeleteNamespaceItem(context.Background(),
		"demo", "DEV", "default", "application1.yaml", "content")
	o.NoError(err)
}

func (o openapiClientTestSuite) Test_GetNamespaceConfig() {
	err := o.apollo.DeleteNamespaceItem(context.Background(),
		"demo", "DEV", "default", "application1.yaml", "content")
	o.NoError(err)
}

func (o openapiClientTestSuite) Test_CreateNamespaceConfig() {
	value := `key: "value"`
	item, err := o.apollo.CreateNamespaceItem(context.Background(),
		"demo", "DEV", "default", "application1.yaml", "content", value)
	o.NoError(err)
	o.T().Logf("%+v\n", item)
}

func (o openapiClientTestSuite) Test_UpdateNamespaceConfig() {
	update := `key: "value"
key2: "value2"
`
	item, err := o.apollo.UpdateNamespaceItem(context.Background(),
		"demo", "DEV", "default", "application1.yaml", "content", update)
	o.NoError(err)
	o.T().Logf("%+v\n", item)
}

func Test_openapi(t *testing.T) {
	suite.Run(t, new(openapiClientTestSuite))
}

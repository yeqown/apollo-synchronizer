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
	appInfos, err := o.apollo.ListNamespaces(context.Background(), "demo", "DEV", "default")
	o.NoError(err)
	o.NotEmpty(appInfos)

	o.T().Logf("appInfos: %+v \n", appInfos)
}

func Test_openapi(t *testing.T) {
	suite.Run(t, new(openapiClientTestSuite))
}

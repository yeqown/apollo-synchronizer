package openapi

// AppInfo
// generated from https://mholt.github.io/json-to-go/
type AppInfo struct {
	Name                       string `json:"name"`
	AppID                      string `json:"appId"`
	OrgID                      string `json:"orgId"`
	OrgName                    string `json:"orgName"`
	OwnerName                  string `json:"ownerName"`
	OwnerEmail                 string `json:"ownerEmail"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

// EnvClusters
// generated from https://mholt.github.io/json-to-go/
type EnvClusters struct {
	Env      string   `json:"env"`
	Clusters []string `json:"clusters"`
}

// ClusterInfo
// generated from https://mholt.github.io/json-to-go/
type ClusterInfo struct {
	Name                       string `json:"name"`
	AppID                      string `json:"appId"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

// NamespaceInfo
// generated from https://mholt.github.io/json-to-go/
type NamespaceInfo struct {
	Name                       string `json:"namespaceName"`
	AppID                      string `json:"appId"`
	Format                     Format `json:"format"`
	IsPublic                   bool   `json:"isPublic"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

// NamespaceItem represents one item(key) in namespace
type NamespaceItem struct {
	Key                        string `json:"key"`
	Value                      string `json:"value"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

// PublishNamespaceResult
// generated from https://mholt.github.io/json-to-go/
type PublishNamespaceResult struct {
	AppID                      string `json:"appId"`
	ClusterName                string `json:"clusterName"`
	NamespaceName              string `json:"namespaceName"`
	Name                       string `json:"name"`
	Comment                    string `json:"comment"`
	DataChangeCreatedBy        string `json:"dataChangeCreatedBy"`
	DataChangeLastModifiedBy   string `json:"dataChangeLastModifiedBy"`
	DataChangeCreatedTime      string `json:"dataChangeCreatedTime"`
	DataChangeLastModifiedTime string `json:"dataChangeLastModifiedTime"`
}

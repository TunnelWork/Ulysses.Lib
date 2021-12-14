package server

type Credential struct {
	CredentialName  string      `json:"credential_name"`
	CredentialValue interface{} `json:"credential_value"` // string, int, bool, []string, []int only
}

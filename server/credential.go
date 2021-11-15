package server

type Credential struct {
	CredentialName  string
	CredentialValue interface{} // string, int, bool, []string, []int only
}

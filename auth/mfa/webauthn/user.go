package webauthn

import (
	"encoding/binary"
	"encoding/json"
	"strings"

	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/duo-labs/webauthn/protocol"
	duo "github.com/duo-labs/webauthn/webauthn"
)

type User struct {
	ID       uint64
	UserName string
	// DisplayName  string // Use UserName instead
	IconURL          string           // Baked in by WebAuthn struct
	AuthnCredentials []duo.Credential // Load from DB (if any)
	// SessionMap       map[string]duo.SessionData // Load from DB (if any)
}

func CreateNewUser(id uint64, name, iconURL string) (*User, error) {
	user := User{
		ID:               id,
		UserName:         name,
		IconURL:          iconURL,
		AuthnCredentials: []duo.Credential{},
		// SessionMap:       map[string]duo.SessionData{},
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	err = auth.InitMFA(id, "webauthn", string(userJson))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func LoadUser(id uint64) (*User, error) {
	userJson, err := auth.CheckoutMFA(id, "webauthn")
	if err != nil {
		return nil, err
	}
	var user User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

func (u *User) WebAuthnName() string {
	return u.UserName
}

func (u *User) WebAuthnDisplayName() string {
	return strings.Split(u.UserName, "@")[0]
}

func (u *User) WebAuthnIcon() string {
	return u.IconURL
}

func (u *User) WebAuthnCredentials() []duo.Credential {
	return u.AuthnCredentials
}

func (u *User) UpdateDatabase() error {
	userJson, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = auth.UpdateMFA(u.ID, "webauthn", string(userJson))
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CredentialExcludeList() []protocol.CredentialDescriptor {
	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.AuthnCredentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}

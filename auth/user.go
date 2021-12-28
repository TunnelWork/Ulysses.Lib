package auth

import (
	"encoding/base64"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

var (
	signer jwt.SigningMethod = &jwt.SigningMethodEd25519{}
)

type User struct {
	// Core Functional Info
	id            uint64
	Email         string `json:"email"`
	PublicKey     string `json:"public_key"` // ed25519.PublicKey in BASE64 representation
	Role          Role   `json:"role"`
	AffiliationID uint64 `json:"affiliation"`

	// Internal Fields for Signing/Verifying
	pubKey interface{}
}

type UserInfo struct {
	// Personal Info - Mandatory
	FirstName string `json:"first_name"` // Preferred First Name
	LastName  string `json:"last_name"`  // Preferred Last Name

	// Billing Info - Optional
	StreetAddress string `json:"street_address"`
	Suite         string `json:"suite"`
	City          string `json:"city"`
	State         string `json:"state"`
	CountryISO    string `json:"country_iso"`
	ZipCode       string `json:"zip_code"`
}

// GetUserByID should be called only after
// the user has been authenticated (Token validated)
func GetUserByID(id uint64) (*User, error) {
	user, err := getUserByID(id)
	if err != nil {
		return nil, err
	}

	// Base64-decode the Public Key
	user.pubKey, err = base64.StdEncoding.DecodeString(user.PublicKey)
	return user, err
}

// GetUserByEmail should be called for user login
// return nil, err when error/mismatch
func GetUserByEmail(email string) (*User, error) {
	user, err := getUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Base64-decode the Public Key
	user.pubKey, err = base64.StdEncoding.DecodeString(user.PublicKey)
	return user, err
}

func GetUsersByAffiliationID(affiliationID uint64) ([]*User, error) {
	users, err := getUsersByAffiliationID(affiliationID)
	if err != nil {
		return nil, err
	}

	var goodUsers []*User = make([]*User, 0)
	for _, user := range users {
		// Base64-decode the Public Key
		user.pubKey, err = base64.StdEncoding.DecodeString(user.PublicKey)
		if err == nil {
			// append
			goodUsers = append(goodUsers, user)
		}
	}

	return goodUsers, nil
}

func ListUserID() ([]uint64, error) {
	return listUserID()
}

func ListUserIDByAffiliationID(affiliationID uint64) ([]uint64, error) {
	return listUserIDByAffiliationID(affiliationID)
}

func (user *User) ID() uint64 {
	return user.id
}

// CreateUser should be called when registering a new user
func (user *User) Create() error {
	exist, err := user.EmailExists()
	if err != nil {
		return err
	}
	if exist {
		return errors.New("auth: email already exists")
	}

	// Check if all fields are valid
	if user.Email == "" {
		return errors.New("auth: email must not be empty")
	}

	return newUser(user)
}

// UserEmailExists should be called before submitting user creation form.
func (user *User) EmailExists() (bool, error) {
	return emailExists(user.Email)
}

// UpdateUser
func (user *User) Update() error {
	return updateUser(user)
}

// Wipe User Data
func (user *User) Wipe() error {
	return wipeUserData(user)
}

func (user *User) CreateInfo(info *UserInfo) error {
	return newUserInfo(user, info)
}

func (user *User) Info() (*UserInfo, error) {
	return getUserInfo(user.id)
}

func (user *User) UpdateInfo(info *UserInfo) error {
	return updateUserInfo(user, info)
}

func (user *User) Verify(msg, signature string) error {
	return signer.Verify(msg, signature, user.pubKey)
}

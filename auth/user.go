package auth

import (
	"errors"
)

type User struct {
	// Core Functional Info
	id            uint64
	Email         string `json:"email"`
	Password      string `json:"password"` // HMAC-Hashed
	Role          Role   `json:"role"`
	AffiliationID uint64 `json:"affiliation"`
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
	return getUserByID(id)
}

// GetUserByEmail should be called for user login
// return nil, err when error/mismatch
func GetUserByEmailPassword(email, password string) (*User, error) {
	return getUserByEmailPassword(email, password)
}

func GetUsersByAffiliationID(affiliationID uint64) ([]*User, error) {
	return getUsersByAffiliationID(affiliationID)
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
	if user.Email == "" || user.Password == "" {
		return errors.New("auth: invalid user data")
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

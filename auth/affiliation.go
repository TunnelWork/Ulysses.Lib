package auth

import (
	"database/sql"
	"errors"
)

var (
	ErrAffiliationNameEmpty           = errors.New("auth: affiliation name is empty")
	ErrAffiliationOwnerUserIDEmpty    = errors.New("auth: affiliation owner user id is empty")
	ErrAffiliationSharedWalletIDEmpty = errors.New("auth: affiliation shared wallet id is empty")
	ErrAffiliationStreetAddressEmpty  = errors.New("auth: affiliation street address is empty")
	ErrAffiliationCityEmpty           = errors.New("auth: affiliation city is empty")
	ErrAffiliationStateEmpty          = errors.New("auth: affiliation state is empty")
	ErrAffiliationCountryISOEmpty     = errors.New("auth: affiliation country iso is empty")
	ErrAffiliationZipCodeEmpty        = errors.New("auth: affiliation zip code is empty")
	ErrAffiliationContactEmailEmpty   = errors.New("auth: affiliation contact email is empty")
)

type Affiliation struct {
	id             uint64 // not to be touched by the user
	Name           string
	ParentID       uint64
	OwnerUserID    uint64 // must be a valid user id with a wallet (to be shared among users with permission)
	SharedWalletID uint64 // must be a valid wallet id
	StreetAddress  string
	Suite          string
	City           string
	State          string
	CountryISO     string
	ZipCode        string
	ContactEmail   string
}

func GetAffiliationByID(id uint64) (*Affiliation, error) {
	return getAffiliationByID(id)
}

func CreateAffiliation(affiliation *Affiliation) error {
	// Check if all fields are valid
	if affiliation.Name == "" {
		return ErrAffiliationNameEmpty
	}
	if affiliation.OwnerUserID == 0 {
		return ErrAffiliationOwnerUserIDEmpty
	}
	if affiliation.SharedWalletID == 0 {
		return ErrAffiliationSharedWalletIDEmpty
	}
	if affiliation.StreetAddress == "" {
		return ErrAffiliationStreetAddressEmpty
	}
	if affiliation.City == "" {
		return ErrAffiliationCityEmpty
	}
	if affiliation.State == "" {
		return ErrAffiliationStateEmpty
	}
	if affiliation.CountryISO == "" {
		return ErrAffiliationCountryISOEmpty
	}
	if affiliation.ZipCode == "" {
		return ErrAffiliationZipCodeEmpty
	}
	if affiliation.ContactEmail == "" {
		return ErrAffiliationContactEmailEmpty
	}

	return newAffiliation(affiliation)
}

func (affiliation *Affiliation) UpdateAffiliation() error {
	// Check if all fields are valid
	if affiliation.Name == "" {
		return ErrAffiliationNameEmpty
	}
	if affiliation.OwnerUserID == 0 {
		return ErrAffiliationOwnerUserIDEmpty
	}
	if affiliation.SharedWalletID == 0 {
		return ErrAffiliationSharedWalletIDEmpty
	}
	if affiliation.StreetAddress == "" {
		return ErrAffiliationStreetAddressEmpty
	}
	if affiliation.City == "" {
		return ErrAffiliationCityEmpty
	}
	if affiliation.State == "" {
		return ErrAffiliationStateEmpty
	}
	if affiliation.CountryISO == "" {
		return ErrAffiliationCountryISOEmpty
	}
	if affiliation.ZipCode == "" {
		return ErrAffiliationZipCodeEmpty
	}
	if affiliation.ContactEmail == "" {
		return ErrAffiliationContactEmailEmpty
	}

	return updateAffiliation(affiliation)
}

func (affiliation *Affiliation) Parent() (*Affiliation, error) {
	affiliation, err := getAffiliationByID(affiliation.ParentID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return affiliation, err
}

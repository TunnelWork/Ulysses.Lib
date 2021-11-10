package utotp

import (
	"errors"

	"github.com/TunnelWork/Ulysses.Lib/auth"
	"github.com/pquerna/otp/totp"
)

type TOTP struct {
	issuer string
}

func NewTOTP(conf map[string]string) *TOTP {
	if issuer, ok := conf["issuer"]; ok {
		return &TOTP{
			issuer: issuer,
		}
	}
	return &TOTP{
		issuer: "Ulysses Unknown Issuer",
	}
}

func (t *TOTP) Registered(userID uint64) bool {
	enabled, err := auth.MFAEnabled(userID, "utotp")
	if err != nil {
		return false
	}
	return enabled
}

func (t *TOTP) InitSignUp(userID uint64, username string) (map[string]interface{}, error) {
	if t.Registered(userID) {
		return nil, errors.New("utotp: user already registered")
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      t.issuer,
		AccountName: username,
	})
	if err != nil {
		return nil, err
	}

	t.Remove(userID)

	err = auth.InitMFA(userID, "utotp", key.Secret())
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"secret": key.Secret(),
		"url":    key.URL(),
	}, nil
}

func (t *TOTP) CompleteSignUp(userID uint64, mfaConf map[string]string) error {
	// Verify if all required params are present
	secret, ok := mfaConf["secret"]
	if !ok {
		return errors.New("utotp: expecting secret in sign up form")
	}

	code, ok := mfaConf["code"]
	if !ok {
		return errors.New("utotp: expecting code in sign up form")
	}

	secretOnRecord, err := auth.CheckoutMFA(userID, "utotp")
	if err != nil {
		return err
	}

	if secret != secretOnRecord {
		return errors.New("utotp: secret does not match")
	}
	// Verify that the code is valid
	if !totp.Validate(code, secret) {
		return errors.New("utotp: code is invalid")
	}

	err = auth.ConfirmMFA(userID, "utotp")
	return err
}

func (t *TOTP) NewChallenge(userID uint64) (map[string]interface{}, error) {
	if t.Registered(userID) {
		return map[string]interface{}{
			"timeout": 30,
		}, nil
	}
	return nil, errors.New("utotp: user not registered")
}

func (t *TOTP) SubmitChallenge(userID uint64, challengeResponse map[string]string) error {
	if !t.Registered(userID) {
		return errors.New("utotp: user not registered")
	}

	code, ok := challengeResponse["code"]
	if !ok {
		return errors.New("utotp: expecting code in sign up form")
	}
	secretOnRecord, err := auth.CheckoutMFA(userID, "utotp")
	if err != nil {
		return err
	}
	if !totp.Validate(code, secretOnRecord) {
		// Verify that the code is valid
		return errors.New("utotp: code is invalid")
	}

	return nil
}

func (t *TOTP) Remove(userID uint64) error {
	return auth.ClearMFA(userID, "utotp")
}

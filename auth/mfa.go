package auth

import "errors"

/**************** Interface ****************/
type MultiFactorAuthentication interface {
	Registered(userID uint64) bool

	// Register associate a MFA credential to user
	InitSignUp(userID uint64, username string) (map[string]interface{}, error)
	CompleteSignUp(userID uint64, mfaConf map[string]string) error

	// Challenge is called when user try to verify identity using the selected MFA.
	NewChallenge(userID uint64) (map[string]interface{}, error)
	SubmitChallenge(userID uint64, challengeResponse map[string]string) error

	// Remove the MFA credential from the database
	Remove(userID uint64) error
}

/**************** Helper Func ****************/

func EnabledMFA(userID uint64) ([]string, error) {
	return checkEnabledMFA(userID)
}

/**************** Aggregator ****************/
var (
	mfaInstanceRegistry = map[string]MultiFactorAuthentication{} // Note: this registry is not thread-safe and is not to be updated by MFA implementation's init() func

	ErrMFAInstanceUnknown = errors.New("auth: Unknown MFA instance")
)

func RegMFAInstance(MFAType string, instance MultiFactorAuthentication) {
	mfaInstanceRegistry[MFAType] = instance
}

func MFARegistered(MFAType string, userID uint64) bool {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.Registered(userID)
	}
	return false
}

func MFAInisSignUp(MFAType string, userID uint64, username string) (map[string]interface{}, error) {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.InitSignUp(userID, username)
	}
	return nil, ErrMFAInstanceUnknown
}

func MFACompleteSignUp(MFAType string, userID uint64, mfaConf map[string]string) error {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.CompleteSignUp(userID, mfaConf)
	}
	return ErrMFAInstanceUnknown
}

func MFANewChallenge(MFAType string, userID uint64) (map[string]interface{}, error) {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.NewChallenge(userID)
	}
	return nil, ErrMFAInstanceUnknown
}

func MFASubmitChallenge(MFAType string, userID uint64, challengeResponse map[string]string) error {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.SubmitChallenge(userID, challengeResponse)
	}
	return ErrMFAInstanceUnknown
}

func MFARemove(MFAType string, userID uint64) error {
	if instance, ok := mfaInstanceRegistry[MFAType]; ok {
		return instance.Remove(userID)
	}
	return ErrMFAInstanceUnknown
}

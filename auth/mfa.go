package auth

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

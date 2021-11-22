package security

func SetupCipher(passwdCipher Cipher) {
	if passwordCipher == nil && passwdCipher != nil {
		passwordCipher = passwdCipher
	} else if passwordCipher != nil {
		panic("SetupCipher: cipher already set")
	} else {
		panic("SetupCipher: input is nil")
	}
}

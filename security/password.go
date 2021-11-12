package security

import harpocrates "github.com/TunnelWork/Harpocrates"

var (
	passwordCipher Cipher = nil
)

func PoorPassword() string {
	passwd, err := harpocrates.GetNewWeakPassword(8)
	if err != nil {
		return ""
	}
	return passwd
}

func StrongPassword() string {
	passwd, err := harpocrates.GetNewStrongPassword(16)
	if err != nil {
		return ""
	}
	return passwd
}

func EncryptPassword(src string) string {
	if passwordCipher == nil {
		return src
	}
	dst, err := passwordCipher.HexDigestEncrypt(src)
	if err != nil {
		return src
	}
	return dst
}

func DecryptPassword(src string) string {
	if passwordCipher == nil {
		return src
	}
	dst, err := passwordCipher.HexDigestDecrypt(src)
	if err != nil {
		return src
	}
	return dst
}

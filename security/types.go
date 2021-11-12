package security

type Cipher interface {
	Encrypt(src []byte) ([]byte, error)
	Decrypt(src []byte) ([]byte, error)
	HexDigestEncrypt(str string) (string, error)
	HexDigestDecrypt(hexstr string) (string, error)
}

package cryptox

type Mode string

const (
	CBC Mode = "CBC" // aes加密-CBC模式
	CFB Mode = "CFB" // aes加密-CFB模式
	ECB Mode = "ECB" // aes加密-ECB模式
	GCM Mode = "GCM" // aes加密-GCM模式

	PKCS1 Mode = "PKCS1" // rsa加密-适用于私钥和公钥
	PKCS8 Mode = "PKCS8" // rsa加密-仅适用于私钥
	PKIX  Mode = "PKIX"  // rsa加密-仅适用于公钥
)

// Crypto 加解密接口
type Crypto interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)
}

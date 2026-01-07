package cryptox

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
)

const (
	PrivateKeyType = "PRIVATE KEY"
	PublicKeyType  = "PUBLIC KEY"
)

// RSA 生成RAS加密对象, 默认密钥长度为1024位
func RSA() (Crypto, error) {
	crypto, err := NewRsaCrypto(1024)
	if err != nil {
		return nil, errorx.Wrap(err, "new rsa crypto error")
	}
	return crypto, nil
}

// NewRsaCrypto 生成RAS加密对象
func NewRsaCrypto(bits int) (*RsaCrypto, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, errorx.Wrap(err, "generate rsa key error")
	}
	return &RsaCrypto{PrivateKey: privateKey}, nil
}

// ParseRsaCrypto 从现有的密钥解析出RAS加密对象
func ParseRsaCrypto(data []byte, mode Mode) (*RsaCrypto, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errorx.New("decode private block error")
	}
	switch mode {
	case PKCS1:
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 private key error")
		}
		return &RsaCrypto{PrivateKey: privateKey}, nil
	case PKCS8:
		privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS8 private key error")
		}
		return &RsaCrypto{PrivateKey: privateKey.(*rsa.PrivateKey)}, nil
	default:
		return nil, errorx.Sprintf("unsupported private mode: %s", mode)
	}
}

// RsaCrypto RSA加密器
type RsaCrypto struct {
	PrivateKey *rsa.PrivateKey // rsa私钥
}

// Encrypt 公钥加密
func (c *RsaCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &c.PrivateKey.PublicKey, plaintext)
}

// Decrypt 私钥解密
func (c *RsaCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, c.PrivateKey, ciphertext)
}

// SavePrivateKey 密钥持久化
func (c *RsaCrypto) SavePrivateKey(path string, mode Mode) error {
	switch mode {
	case PKCS1:
		return writePemBlock(path, &pem.Block{
			Type:  PrivateKeyType,
			Bytes: x509.MarshalPKCS1PrivateKey(c.PrivateKey),
		})
	case PKCS8:
		data, err := x509.MarshalPKCS8PrivateKey(c.PrivateKey)
		if err != nil {
			return errorx.Wrap(err, "marshal PKCS8 private key error")
		}
		return writePemBlock(path, &pem.Block{
			Type:  PrivateKeyType,
			Bytes: data,
		})
	default:
		return errorx.Sprintf("unsupported private key mode: %s", mode)
	}
}

// SavePublicKey 公钥持久化
func (c *RsaCrypto) SavePublicKey(path string, mode Mode) error {
	switch mode {
	case PKCS1:
		return writePemBlock(path, &pem.Block{
			Type:  PublicKeyType,
			Bytes: x509.MarshalPKCS1PublicKey(&c.PrivateKey.PublicKey),
		})
	case PKIX:
		data, err := x509.MarshalPKIXPublicKey(&c.PrivateKey.PublicKey)
		if err != nil {
			return errorx.Wrap(err, "marshal PKIX public key error")
		}
		return writePemBlock(path, &pem.Block{
			Type:  PublicKeyType,
			Bytes: data,
		})
	default:
		return errorx.Sprintf("unsupported public key mode: %s", mode)
	}
}

// RsaEncrypt 公钥加密
func RsaEncrypt(plaintext, publicKey []byte, mode Mode) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errorx.New("decode public key error")
	}
	switch mode {
	case PKCS1:
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 public key error")
		}
		return rsa.EncryptPKCS1v15(rand.Reader, key, plaintext)
	case PKIX:
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKIX public key error")
		}
		return rsa.EncryptPKCS1v15(rand.Reader, key.(*rsa.PublicKey), plaintext)
	default:
		return nil, errorx.Sprintf("unsupported public key mode: %s", mode)
	}
}

// ParseRsaPrivateKey 从私钥block解析私钥
func ParseRsaPrivateKey(privateKey []byte, mode Mode) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errorx.New("decode public key error")
	}
	switch mode {
	case PKCS1:
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 private key error")
		}
		return key, nil
	case PKCS8:
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS8 private key error")
		}
		return key.(*rsa.PrivateKey), nil
	default:
		return nil, errorx.Sprintf("unsupported private key mode: %s", mode)
	}
}

// ParseRsaPublicKey 从公钥block解析公钥
func ParseRsaPublicKey(publicKey []byte, mode Mode) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errorx.New("decode public key error")
	}
	switch mode {
	case PKCS1:
		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKCS1 public key error")
		}
		return key, nil
	case PKIX:
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errorx.Wrap(err, "parse PKIX public key error")
		}
		return key.(*rsa.PublicKey), nil
	default:
		return nil, errorx.Sprintf("unsupported public key mode: %s", mode)
	}
}

// writePemBlock 写入pem文件
func writePemBlock(path string, block *pem.Block) error {
	var buf bytes.Buffer
	if err := pem.Encode(&buf, block); err != nil {
		return errorx.Wrap(err, "pem encode error")
	} else if err = filex.WriteFile(path, buf.Bytes()); err != nil {
		return errorx.Wrap(err, "pem write error")
	}
	return nil
}

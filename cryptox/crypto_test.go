package cryptox

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-xuan/utilx/encodingx"
	"github.com/go-xuan/utilx/filex"
)

func TestAES(t *testing.T) {
	key := "fd6d1c1fb333a9ab4e17ce683c988d75"
	iv := "b43f9be55644daca"
	aesCrypto, err := NewAesCrypto(key, iv, GCM)
	if err != nil {
		fmt.Println(err)
	}

	bytes, _ := json.Marshal(struct {
		AppID     string `json:"app_id"`
		Timestamp int64  `json:"timestamp"`
	}{
		AppID:     "1234567890",
		Timestamp: time.Now().Unix(),
	})

	var ciphertext, plaintext []byte
	if ciphertext, err = aesCrypto.Encrypt(bytes); err != nil {
		panic(err)
	}
	fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

	if plaintext, err = aesCrypto.Decrypt(ciphertext); err != nil {
		panic(err)
	}
	fmt.Println("解密：", string(plaintext))
}

func TestRsa(t *testing.T) {
	t.Run("Rsa", func(t *testing.T) {
		crypto, err := NewRsaCrypto(1024)
		if err != nil {
			fmt.Println(err)
			return
		}
		plaintext := []byte("hello world")
		var ciphertext []byte
		if ciphertext, err = crypto.Encrypt(plaintext); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))

		if plaintext, err = crypto.Decrypt(ciphertext); err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("解密：", string(plaintext))

		if err = crypto.SavePrivateKey("./rsa/private.pem", PKCS1); err != nil {
			fmt.Println(err)
			return
		}
		if err = crypto.SavePublicKey("./rsa/public.pem", PKCS1); err != nil {
			fmt.Println(err)
			return
		}
	})

	t.Run("public key", func(t *testing.T) {
		publicKey, err := filex.ReadFile("./rsa/public.pem")
		if err != nil {
			fmt.Println(err)
			return
		}
		plaintext := []byte("hello world")
		ciphertext, err := RsaEncrypt(plaintext, publicKey, PKCS1)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("加密：", encodingx.Base64().Encode(ciphertext))
	})
}

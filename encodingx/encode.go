package encodingx

// Encoder 加密器
type Encoder interface {
	Encode(plaintext []byte) string
}

// Decoder 解密器
type Decoder interface {
	Decode(ciphertext string) ([]byte, error)
}

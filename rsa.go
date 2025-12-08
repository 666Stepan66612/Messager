package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
	"time"
)

type rsaEnv struct {
	n *big.Int
	e *big.Int
	d *big.Int
}

type publicRSAkey struct {
	e *big.Int
	n *big.Int
}

type privateRSAkey struct {
	d *big.Int
	n *big.Int
}

func (r *rsaEnv) SetRSAEnv() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	r.n = key.N
	r.e = big.NewInt(int64(key.E))
	r.d = key.D

	return nil
}

func (p *publicRSAkey) SetPublicRSAkeyEnv(rsa rsaEnv) {
	p.e = rsa.e
	p.n = rsa.n
}

func (p *privateRSAkey) SetPrivateRSAkeyEnv(rsa rsaEnv) {
	p.d = rsa.d
	p.n = rsa.n
}

func (p *publicRSAkey) GetPublicRSAkeyEnv() (*big.Int, *big.Int) {
	return p.e, p.n
}

func (p *privateRSAkey) GetPrivateRSAkeyEnv() (*big.Int, *big.Int) {
	return p.d, p.n
}

func AES256Encrypt(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func AES256Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	start := time.Now()
	message := "Hardest hello world!"
	rsa := rsaEnv{}
	var publicRSAKey publicRSAkey
	var privateRSAKey privateRSAkey

	rsa.SetRSAEnv()
	publicRSAKey.SetPublicRSAkeyEnv(rsa)
	publicRSAKey.GetPublicRSAkeyEnv()
	privateRSAKey.SetPrivateRSAkeyEnv(rsa)
	privateRSAKey.GetPrivateRSAkeyEnv()

	aesKey := make([]byte, 32)
	rand.Read(aesKey)

	publicTextUnderAES := AES256Encrypt(aesKey, []byte(message))

	aesKeyBig := new(big.Int).SetBytes(aesKey)
	encryptedAESKey := new(big.Int).Exp(aesKeyBig, publicRSAKey.e, publicRSAKey.n)

	decryptedAESKey := new(big.Int).Exp(encryptedAESKey, privateRSAKey.d, privateRSAKey.n)

	decrypt, err := AES256Decrypt(decryptedAESKey.Bytes(), publicTextUnderAES)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(decrypt))   // Hardest hello world!
	fmt.Println(time.Since(start)) // ~110 ms
}
package ecdh

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/ecdh"
    "crypto/rand"
    "crypto/sha256"
    "fmt"
    "io"
    "time"
    "golang.org/x/crypto/hkdf"
)

func AES256Encrypt(key, plaintext []byte) []byte {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func AES256Decrypt(key, ciphertext []byte) string {
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
    return string(plaintext)
}

func DeriveAESKey(sharedSecret []byte) []byte {
    salt := make([]byte, 32)
    h := hkdf.New(sha256.New, sharedSecret, salt, []byte("AES256 key for ECDH"))
    key := make([]byte, 32)
    h.Read(key)
    return key
}

func CreateKeyPair(curve ecdh.Curve) (*ecdh.PrivateKey, *ecdh.PublicKey) {
	priv, _ := curve.GenerateKey(rand.Reader)
    pub := priv.PublicKey()
	
	return priv, pub
}

func ComputeSharedSecret(priv *ecdh.PrivateKey, pub *ecdh.PublicKey) []byte {
    secret, _ := priv.ECDH(pub)
    return secret
}

func Ecdh(){
	start := time.Now()
	message := "Hardest hello world!"
	curve := ecdh.P256()

	bobPriv, bobPub := CreateKeyPair(curve)
	alicePriv, alicePub := CreateKeyPair(curve)

	sharedSecretAlice := ComputeSharedSecret(alicePriv, bobPub)
	sharedSecretBob := ComputeSharedSecret(bobPriv, alicePub)

	aesKeyAlice := DeriveAESKey(sharedSecretAlice)
	aesKeyBob := DeriveAESKey(sharedSecretBob)

	ciphertext := AES256Encrypt(aesKeyAlice, []byte(message))
	decrypted := AES256Decrypt(aesKeyBob, ciphertext)

	fmt.Println(decrypted)
	fmt.Println(time.Since(start))
}
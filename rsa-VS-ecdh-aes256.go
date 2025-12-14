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
	"messager/ecdh"
)

//This is rsa-aes cycle: When a user writes a message, they send a request to another user to generate two RSA keys.
//  After generation, the sender of the message receives the public key from the intended recipient,
//  which does not need to be hidden from a third-party observer.
//  Next, the sender applies the AES256 algorithm to their message,
//  resulting in encrypted content, which also does not need to be hidden, and sends it to the recipient.
//  Then, using the public RSA key received from the recipient, the sender encrypts the temporary key used for encrypting their message.
//  The resulting ciphertext from applying the RSA public key to the temporary key also does not need to be hidden and is sent to the recipient.
//  The recipient receives the encrypted message (which was encrypted using the public AES256 key)
//  and the AES256 key in encrypted form (encrypted with the public RSA key).
//  They decrypt the AES256 key using their private RSA key and apply it to the encrypted message.

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

func (r *rsaEnv) SetRSAEnv() error { // create new rsa env
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	r.n = key.N
	r.e = big.NewInt(int64(key.E))
	r.d = key.D

	return nil
}

func (p *publicRSAkey) SetPublicRSAkeyEnv(rsa rsaEnv) { // set public key env from rsa struct
	p.e = rsa.e
	p.n = rsa.n
}

func (p *privateRSAkey) SetPrivateRSAkeyEnv(rsa rsaEnv) { // set private key env from rsa struct
	p.d = rsa.d
	p.n = rsa.n
}

func (p *publicRSAkey) GetPublicRSAkeyEnv() (*big.Int, *big.Int) { // get public key env for rsa
	return p.e, p.n
}

func (p *privateRSAkey) GetPrivateRSAkeyEnv() (*big.Int, *big.Int) { // get private key env for rsa
	return p.d, p.n
}

func AES256Encrypt(key, plaintext []byte) []byte { // aes encrypts the message that the user wants to send with the key that aes encryption will be encrypted with the public key rsa
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, plaintext, nil)
}

func AES256Decrypt(key, ciphertext []byte) ([]byte, error) { // after the recipient receives the encrypted message, aes decrypts the encrypted message using the same key
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func main() {
	start := time.Now()
	message := "Hardest hello world!" // create message
	rsa := rsaEnv{}
	var publicRSAKey publicRSAkey
	var privateRSAKey privateRSAkey

	rsa.SetRSAEnv()
	publicRSAKey.SetPublicRSAkeyEnv(rsa)
	publicRSAKey.GetPublicRSAkeyEnv()
	privateRSAKey.SetPrivateRSAkeyEnv(rsa)
	privateRSAKey.GetPrivateRSAkeyEnv()

	aesKey := make([]byte, 32) // create aes key
	rand.Read(aesKey)

	publicTextUnderAES := AES256Encrypt(aesKey, []byte(message)) // 1. encrypt message by aes key; 2. send to message recipient an encrypted message to the recipient, it is not secret

	aesKeyBig := new(big.Int).SetBytes(aesKey)
	encryptedAESKey := new(big.Int).Exp(aesKeyBig, publicRSAKey.e, publicRSAKey.n) // 3. apply a public rsa key to the aes key; 4. send the aes key in encrypted form, it is not secret

	decryptedAESKey := new(big.Int).Exp(encryptedAESKey, privateRSAKey.d, privateRSAKey.n) // 5. decrypt the aes key we obtained in step 4 using the private rsa key

	decrypt, err := AES256Decrypt(decryptedAESKey.Bytes(), publicTextUnderAES) // 6. decrypt the encrypted message from step 2 using the decrypted rsa key
	if err != nil {
		panic(err)
	}

	fmt.Println("Using rsa")
	fmt.Println(string(decrypt))   // Hardest hello world!
	fmt.Println(time.Since(start)) // ~400ms - 4sec

	fmt.Println("Using ecdh")
	ecdh.Ecdh() // very very very fast ~200-900µs sometimes module time write 0s
}
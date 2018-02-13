// enrypt and decrypt using AES
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"testing"
)

func ExAes() {
	plainText := []byte("dummy,2018-05-05")
	//key := []byte("passw0rdpassw0rdpassw0rdpassw0rd")
	key := []byte("datasectionmlpkg")

	// Create new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	fmt.Println(aead.NonceSize())
	nonce := make([]byte, aead.NonceSize())
	/* randomaize */
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}

	cipherText := aead.Seal(nil, nonce, plainText, nil)
	fmt.Printf("Cipher text: %x\n", cipherText)

	plainText_, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return
	}
	fmt.Printf("Decrypted text: %s\n", string(plainText_))
}

func main() {
	ExAes()
}

func (t *Token) NewAES() cipher.AEAD {
	// Create new AES cipher block
	block, err := aes.NewCipher([]byte(token.SecretKey))
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	return aead

}

func (t *Token) EncryptAES(plainText string) ([]byte, error) {
	// Create new AES cipher block
	block, err := aes.NewCipher([]byte(token.SecretKey))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return []byte(""), err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return []byte(""), err
	}

	nonce := make([]byte, aead.NonceSize())
	/* randomaize
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Printf("err: %s\n", err)
		return []byte(""), err
	}
	*/

	cipherText := aead.Seal(nil, nonce, []byte(plainText), nil)
	return cipherText, nil
}

func (t *Token) DecryptAES(cipherText []byte) ([]byte, error) {
	aead := token.NewAES()
	nonce := make([]byte, aead.NonceSize())
	plainText, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		fmt.Printf("err: %s\n", err)
		return []byte(""), err
	}
	return plainText, nil
}

func TestEncryptAES(t *testing.T) {
	table := []struct {
		plainText string
	}{
		{
			plainText: "dummy1,1900-05-05",
		},
		{
			plainText: "dummy,2505-05-05",
		},
	}

	for _, row := range table {
		cipherText, err := token.EncryptAES(row.plainText)
		if err != nil {
			t.Fatalf("got error in Encrypt: %s", err)
		}
		fmt.Printf("encrypted AES!: %x\n", cipherText)
	}

}

func TestDecryptAES(t *testing.T) {
	table := []struct {
		cipherText string
	}{
		{
			cipherText: "0e8bf6f9a91d689180d8ef1c471297b1c760875be1be651f8d1d2bf03ff12fd02b",
		},
		{
			cipherText: "0e8bf6f9a900769589ddf201420a8ab46c30064b32771b55f30d146815c4fa5f",
		},
	}

	for _, row := range table {
		// 16進数 -> byte
		decHex, _ := hex.DecodeString(row.cipherText)
		plainText, err := token.DecryptAES(decHex)
		if err != nil {
			t.Fatalf("got error in Decrypt: %s", err)
		}
		fmt.Printf("decrypted AES!: %v\n", string(plainText))
	}

}

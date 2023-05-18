package tokens

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _key []byte
var _iv []byte

func Init(passphrase string) {
	_key, _iv = makeKeyIv(passphrase)
}

func ValidateToken(r *http.Request) error {
	// for name, values := range r.Header {
	// 	// Loop over all values for the name.
	// 	for _, value := range values {
	// 		fmt.Println(name, value)
	// 	}
	// }

	token := r.Header.Get("X-AUTH-TOKEN") + ""
	if token == "" {
		return fmt.Errorf("invalid token or token missing")
	}

	// try to decrypt token
	dec, err := decrypt(token)
	if err != nil {
		return err
	}
	//fmt.Println(dec)

	// make sure token is valid - should be generated within past 2 minutes
	// token is time as UnixNano
	//
	// convert decrypted token to int64 (unixnano)
	un, err := strconv.ParseInt(dec, 10, 64)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid token")
	}

	t1 := time.Now()
	t2 := time.Unix(0, un)
	if t1.Sub(t2).Seconds() > 120 {
		return fmt.Errorf("token expired")
	}
	return nil
}

func decrypt(ciphertext string) (string, error) {
	// create cipher
	block, err := aes.NewCipher(_key)
	if err != nil {
		return "", err
	}

	// allocate space for ciphered data
	encbytes, err := hex.DecodeString(ciphertext)
	out := make([]byte, len(encbytes))

	aesDecrypter := cipher.NewCFBDecrypter(block, _iv)
	if err != nil {
		return "", err
	}
	aesDecrypter.XORKeyStream(out, encbytes)

	// return string and strip-out line-feed chars i.e. ascii(10)
	idx := bytes.IndexByte(out, '\r')
	if idx < 1 {
		idx = len(out)
	}
	return string(out[:idx]), nil
}

func makeKeyIv(passphrase string) ([]byte, []byte) {
	// make key of 32 bit long from passphrase
	ln := len(passphrase)
	if ln < 32 {
		passphrase = passphrase + strings.Repeat("~", 32-ln)
	}
	key := []byte(passphrase[0:32])
	iv := []byte(passphrase[0:16])
	return key, iv
}

// // -------------
// // Sample to generate token at client
// //
// func generateToken() (string, error) {
// 	plain_token := strconv.FormatInt(time.Now().UnixNano(), 10)
// 	// create cipher
// 	block, err := aes.NewCipher(_key)
// 	if err != nil {
// 		return "", err
// 	}
// 	// allocate space for ciphered data
// 	out := make([]byte, len(plain_token))

// 	aesEncrypter := cipher.NewCFBEncrypter(block, _iv)
// 	aesEncrypter.XORKeyStream(out, []byte(plain_token))

// 	// return hex string
// 	return hex.EncodeToString(out), nil
// }

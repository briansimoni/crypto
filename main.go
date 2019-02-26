package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var key []byte

func main() {

	http.HandleFunc("/encrypt", http.HandlerFunc(encryptionHandler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listening on " + port)
	err := http.ListenAndServe(":"+port, nil)
	log.Println(err.Error())
}

func encryptionHandler(w http.ResponseWriter, r *http.Request) {
	// read the plaintext from the body
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cipherText, err := encrypt(key, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(cipherText))
}

// encrypt will aes256 encrypt provided plaintext and return base64 encoded ciphertext
func encrypt(key []byte, plaintext []byte) (string, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := make([]byte, 0)
	cipherText = gcm.Seal(cipherText, nonce, plaintext, nil)
	encodedCipherText := base64.StdEncoding.EncodeToString(cipherText)
	return encodedCipherText, nil
}

func init() {
	key = make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err.Error())
	}
}

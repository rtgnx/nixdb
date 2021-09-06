package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/Reverse-Labs/nixdb"
	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	User  nixdb.PasswdEntry
	Read  bool
	Write bool
	jwt.StandardClaims
}

type JWTSecretKey []byte

func NewJWTClaim(user nixdb.PasswdEntry, read, write bool) JWTClaim {
	return JWTClaim{
		user, true, false,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
}

func JWTToken(secret JWTSecretKey, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSigned, err := token.SignedString([]byte(secret))

	return tokenSigned, err
}

func NewJWTSecretKey(bitSize uint) (JWTSecretKey, error) {
	key := JWTSecretKey(make([]byte, bitSize/8))
	if n, err := rand.Read(key); n != int(bitSize)/8 || err != nil {
		return key, fmt.Errorf("not enough entropy to generate key of size %d", bitSize)
	}

	return key, nil
}

func (k JWTSecretKey) Write(r io.Writer) error {
	if n, err := r.Write(k); n != len(k) || err != nil {
		return fmt.Errorf("unable to write key file")
	}

	return nil
}

func (k *JWTSecretKey) Read(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)

	if err == nil {
		*k = JWTSecretKey(buf)
	}

	return err
}

func WriteNewKey(keyPath string, bitSize uint) error {
	fd, err := os.OpenFile(keyPath, os.O_CREATE, 0600)

	if err != nil {
		log.Fatalln(err.Error())
	}

	defer fd.Close()

	key, err := NewJWTSecretKey(bitSize)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return key.Write(fd)
}

func ReadJWTSecret(keyPath string) (JWTSecretKey, error) {
	key := make(JWTSecretKey, 0)

	fd, err := os.Open(keyPath)

	if err != nil {
		return key, err
	}

	return key, key.Read(fd)
}

func ReadOrGenerateJWTSecret(keyPath string) (JWTSecretKey, error) {
	key, err := ReadJWTSecret(keyPath)

	if err != nil {
		log.Printf("[-] Key is not present, generating...")
		if err := WriteNewKey(keyPath, 4096); err != nil {
			return key, err
		}

		return ReadJWTSecret(keyPath)
	}

	return key, err
}

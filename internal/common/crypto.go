// Package common contains functions, structs, wrappers and others that is reusable
// throughout entire codebase
package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sweet-go/stdlib/encryption"
	"golang.org/x/crypto/bcrypt"
)

// DefaultBlockSize is the default block size used for encryption/decryption
const DefaultBlockSize int = 16

// SharedCryptor instance contains common functionality relates to cryptograpic functions
// and is reusable throughout the entire codebase
type SharedCryptor struct {
	encryptionKey []byte
	iv            string
	blockSize     int
	hashCost      int
}

// SharedCryptorIface interface for SharedCryptor. Provided to ease the mocking process later
type SharedCryptorIface interface {
	Encrypt(plainText string) (string, error)
	Hash(data []byte) (string, error)
	CreateJWT(claims jwt.Claims) (string, error)
	ValidateJWT(token string, opts ValidateJWTOpts) (*jwt.Token, error)
	CompareHash(hashed []byte, plain []byte) error
}

// CreateCryptorOpts is the options used to create a new cryptor instance.
type CreateCryptorOpts struct {
	HashCost      int
	EncryptionKey []byte
	IV            string
	BlockSize     int
}

// NewSharedCryptor create a new instance of SharedCryptor
func NewSharedCryptor(opts *CreateCryptorOpts) *SharedCryptor {
	return &SharedCryptor{
		encryptionKey: encryption.SHA256Hash(opts.EncryptionKey), // better implement hkdf
		iv:            opts.IV,
		blockSize:     opts.BlockSize,
		hashCost:      opts.HashCost,
	}
}

// Encrypt takes plainText and returning the encrypted form of it
func (s *SharedCryptor) Encrypt(plainText string) (string, error) {
	ivKey, err := hex.DecodeString(s.iv)
	if err != nil {
		return "", err
	}

	bPlaintext := s.pkcs5Padding([]byte(plainText), s.blockSize, len(plainText))

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(bPlaintext))

	mode := cipher.NewCBCEncrypter(block, ivKey)
	mode.CryptBlocks(ciphertext, bPlaintext)

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt takes the cipherText and returning the decrypted value of it
func (s *SharedCryptor) Decrypt(cipherText string) (string, error) {
	ivKey, err := s.generateIVKey(s.iv)
	if err != nil {
		return "", err
	}

	cipherTextDecoded, err := hex.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, ivKey)
	mode.CryptBlocks(cipherTextDecoded, cipherTextDecoded)

	return string(s.pkcs5Unpadding(cipherTextDecoded)), nil
}

// Hash generates hashed value utilizing bcrypt of data in form of base64 encoded string
func (s *SharedCryptor) Hash(data []byte) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(data, s.hashCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hashed), nil
}

// CompareHash check whether hashed is equal from plain
func (s *SharedCryptor) CompareHash(hashed []byte, plain []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, plain)
}

// CreateJWT create a jwt token
func (s *SharedCryptor) CreateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(s.encryptionKey)
	if err != nil {
		return "", err
	}

	return signed, nil
}

// ValidateJWTOpts options to validate the JWT token
type ValidateJWTOpts struct {
	Issuer  string
	Subject string
}

// ValidateJWT validate JWT with several extra checks provided by the jwt module, such as:
// WithExpirationRequired, WithIssuer, WithSubject, WithValidMethods
func (s *SharedCryptor) ValidateJWT(token string, opts ValidateJWTOpts) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, errors.New("invalid token alg")
			}

			return s.encryptionKey, nil
		},
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(opts.Issuer),
		jwt.WithSubject(opts.Subject),
		jwt.WithIssuedAt(),
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Name,
		}),
	)

	if err != nil {
		return nil, err
	}

	// safety check
	if !parsedToken.Valid {
		return nil, errors.New("jwt token is not a valid token")
	}

	return parsedToken, nil
}

func (s *SharedCryptor) pkcs5Unpadding(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}

	length := len(src)
	unpadding := int(src[length-1])
	cutLen := (length - unpadding)
	// check boundaries
	if cutLen < 0 || cutLen > length {
		return src
	}

	return src[:cutLen]
}

func (s *SharedCryptor) pkcs5Padding(ciphertext []byte, blockSize int, _ int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func (s *SharedCryptor) generateIVKey(iv string) ([]byte, error) {
	if len(iv) > 0 {
		ivKey, err := hex.DecodeString(iv)
		if err != nil {
			return nil, errors.New("unable to hex decode iv")
		}

		return ivKey, nil
	}

	ivKey, err := generateRandomIVKey(s.blockSize)
	if err != nil {
		return nil, errors.New("unable to generate random iv key")
	}

	return hex.DecodeString(ivKey)
}

// GenerateRandomIVKey generate random IV value
func generateRandomIVKey(blockSize int) (string, error) {
	bytes := make([]byte, blockSize)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

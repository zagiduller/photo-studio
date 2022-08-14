package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"photostudio/components"
)

// @project photo-studio
// @created 10.08.2022

type Service struct {
	db            *gorm.DB
	pk            *ecdsa.PrivateKey
	pub           crypto.PublicKey
	signingMethod *jwt.SigningMethodECDSA
}

func New() (*Service, error) {
	return &Service{}, nil
}

func (s *Service) Configure() error {
	pkPath := viper.GetString("components.auth.privateKey")
	if pkPath == "" {
		return fmt.Errorf("auth.Configure: privateKey is empty ")
	}
	db := components.GetDB()
	if db == nil {
		return fmt.Errorf("auth.Configure: %w ", components.ErrorCodeDbIsNil)
	}
	s.db = db

	pk, err := configurePrivateKey(pkPath)
	if err != nil {
		return fmt.Errorf("auth.Configure: %w ", err)
	}

	s.pk, s.pub = pk, pk.Public()
	s.signingMethod = jwt.SigningMethodES256

	return nil
}

// configurePrivateKey init Parse file and parse to ECDSA key
func configurePrivateKey(pkPath string) (*ecdsa.PrivateKey, error) {
	file, err := os.OpenFile(pkPath, os.O_RDONLY, os.ModeType)
	if err != nil {
		return nil, fmt.Errorf("auth.ConfigurePrivateKey: %w ", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("auth.ConfigurePrivateKey: %w ", err)
	}
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("auth.ConfigurePrivateKey: block is nil ")
	}
	signKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("auth.ConfigurePrivateKey: %w ", err)
	}
	return signKey, nil
}

func (s *Service) CreateToken(custom jwt.Claims) (string, error) {
	t := jwt.New(s.signingMethod)

	t.Claims = custom
	return t.SignedString(s.pk)
}

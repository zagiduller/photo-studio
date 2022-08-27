package access

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/zagiduller/photo-studio/components"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
)

// @project photo-studio
// @created 10.08.2022

type Service struct {
	components.Default
	db            *gorm.DB
	pk            *ecdsa.PrivateKey
	pub           crypto.PublicKey
	signingMethod *jwt.SigningMethodECDSA
}

func New() *Service {
	return &Service{
		Default: components.DefaultComponent("access"),
	}
}

func (s *Service) Configure(ctx context.Context) error {
	s.Default.Ctx = ctx
	pkPath := viper.GetString("components.access.privateKey")
	if pkPath == "" {
		return fmt.Errorf("access.Configure: privateKey is empty ")
	}
	s.db = components.GetDB()
	if s.db == nil {
		return fmt.Errorf("access.Configure: %w ", components.ErrorCodeDbIsNil)
	}
	// migrate model
	if err := s.db.AutoMigrate(&Access{}); err != nil {
		return fmt.Errorf("access.Configure: %w ", err)
	}

	pk, err := configurePrivateKey(pkPath)
	if err != nil {
		return fmt.Errorf("access.Configure: %w ", err)
	}

	s.pk, s.pub = pk, pk.Public()
	s.signingMethod = jwt.SigningMethodES256

	return nil
}

// configurePrivateKey init Parse file and parse to ECDSA key
func configurePrivateKey(pkPath string) (*ecdsa.PrivateKey, error) {
	path, err := filepath.Abs(pkPath)
	if err != nil {
		return nil, fmt.Errorf("access.configurePrivateKey: %w ", err)
	}
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModeType)
	if err != nil {
		return nil, fmt.Errorf("access.ConfigurePrivateKey: %w ", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("access.ConfigurePrivateKey: %w ", err)
	}
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("access.ConfigurePrivateKey: block is nil ")
	}
	signKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("access.ConfigurePrivateKey: %w ", err)
	}
	return signKey, nil
}

func (s *Service) CreateToken(custom jwt.Claims) (string, error) {
	t := jwt.New(s.signingMethod)

	t.Claims = custom
	return t.SignedString(s.pk)
}
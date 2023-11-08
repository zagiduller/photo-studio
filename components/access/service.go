package access

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/zagiduller/photo-studio/components"
	"github.com/zagiduller/photo-studio/components/users"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// @project photo-studio
// @created 10.08.2022
// @arthur

type Service struct {
	components.Default

	ctx context.Context

	db  *gorm.DB
	pk  *ecdsa.PrivateKey
	pub crypto.PublicKey

	users *users.Service
}

func (s *Service) ConfigureDependencies(components []components.Component) {
	for _, c := range components {
		if srv, ok := c.(*users.Service); ok {
			s.users = srv
			s.GetLogger().
				WithField("dependency", srv.GetName()).
				Infof("ConfigureDependencies")
			break
		}
	}
}

var signingMethod = jwt.SigningMethodES256

func New(ctx context.Context) *Service {
	return &Service{
		ctx:     ctx,
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
	if err := s.db.AutoMigrate(&Access{}, &Password{}, &Login{}); err != nil {
		return fmt.Errorf("access.Configure: %w ", err)
	}

	pk, err := configurePrivateKey(pkPath)
	if err != nil {
		return fmt.Errorf("access.Configure: %w ", err)
	}

	s.pk, s.pub = pk, pk.Public()

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
	t := jwt.New(signingMethod)

	t.Claims = custom
	return t.SignedString(s.pk)
}

func (s *Service) CreateTokenByLogin(db *gorm.DB, login *Login) (*Access, error) {
	if login == nil || login.GetUser() == nil {
		return nil, fmt.Errorf("CreateTokenByUser: [%w] ", components.ErrModelIsNil)
	}
	user := login.GetUser()
	sess := uuid.NewString()
	tokenString, err := s.CreateToken(Claims{
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    s.GetName(),
			Audience:  nil,
			NotBefore: nil,
			Subject:   fmt.Sprintf("%d", login.ID),
			ExpiresAt: &jwt.NumericDate{Time: time.Now().AddDate(0, 1, 0)},
			IssuedAt:  &jwt.NumericDate{Time: time.Now()},
			ID:        sess,
		},
		User: user.GetClaims(),
	})
	if err != nil {
		return nil, fmt.Errorf("CreateTokenByUser: [%w]", err)
	}
	token := NewAccess(user.ID, tokenString)
	token.SessionID = sess
	token.SetDB(db)
	if _, err := s.ValidateToken(tokenString); err != nil {
		return nil, fmt.Errorf("CreateTokenByUser: [%w]", err)
	}

	if err := token.Save(); err != nil {
		return nil, fmt.Errorf("CreateTokenByUser: [%w]", err)
	}
	return token, nil
}

// ValidateToken validates the token and returns the claims
func (s *Service) ValidateToken(token string) (*jwt.Token, error) {
	t, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if signingMethod != token.Method {
			return nil, fmt.Errorf("ValidateToken: Unexpected signing method: %v ", token.Header["alg"])
		}
		return s.pub, nil
	})
	if err != nil {
		return nil, fmt.Errorf("ValidateToken: %s ", err)
	}
	if !t.Valid {
		return nil, fmt.Errorf("ValidateToken: [%w] ", ErrInvalidAccessToken)
	}
	return t, nil
}

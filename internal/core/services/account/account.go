package account

import (
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/stewie1520/ecommerce-backend/internal/core/services/hash"
	"math/rand"
	"time"

	"github.com/stewie1520/ecommerce-backend/internal/core"
)

type IAccountService interface {
	RequestOtpToVerify(rs core.RequestScope, email string) (string, time.Time, error)
	VerifyOTP(rs core.RequestScope, email string, sessionId string, otp string) error
	CreateRegisteringSessionId(rs core.RequestScope, email string) (string, time.Time, error)
	CreateUserByEmail(rs core.RequestScope, email string, password string, name string, birthday string, verify bool) (string, error)
	VerifyRegisteringSessionId(rs core.RequestScope, registeringSessionId string, email string) error
}

type AccountService struct {
	cache *redis.Client
	hash  hash.IHashService
}

func NewAccountService(cache *redis.Client, hash hash.IHashService) *AccountService {
	return &AccountService{
		cache: cache,
		hash:  hash,
	}
}

func (as AccountService) generateOtpValue(email string) (otpVal string, otp string) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	otp = fmt.Sprintf("%d", r.Intn(9000)+1000)
	return as.formatOTP(email, otp), otp
}

func (as AccountService) formatOTP(email string, otp string) string {
	return fmt.Sprintf("%s:%s", email, otp)
}

func (as AccountService) generateOtpKey() string {
	return fmt.Sprintf("otp.%s", uuid.NewString())
}

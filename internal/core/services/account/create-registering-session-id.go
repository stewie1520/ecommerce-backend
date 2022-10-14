package account

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stewie1520/ecommerce-backend/internal/core"
)

func (as *AccountService) CreateRegisteringSessionId(rs core.RequestScope, email string) (string, time.Time, error) {
	sessionId := fmt.Sprintf("registering.%s", uuid.NewString())
	expiredIn := time.Hour * 24
	err := as.cache.Set(rs.GetContext(), sessionId, email, expiredIn).Err()
	if err != nil {
		return "", time.Time{}, err
	}

	return sessionId, time.Now().Add(expiredIn), nil
}

package account

import (
	"fmt"

	"github.com/stewie1520/ecommerce-backend/internal/core"
)

func (as *AccountService) VerifyRegisteringSessionId(rs core.RequestScope, registeringSessionId string, email string) error {
	val, err := as.cache.Get(rs.GetContext(), registeringSessionId).Result()
	if err != nil {
		return err
	}

	if val != email {
		return fmt.Errorf("invalid registeringSessionId\n")
	}

	_, err = as.cache.Del(rs.GetContext(), registeringSessionId).Result()
	if err != nil {
		return err
	}

	return nil
}

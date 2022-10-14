package core

import (
	"context"

	"github.com/pocketbase/pocketbase"
	pb_models "github.com/pocketbase/pocketbase/models"
)

type RequestScope interface {
	// RequestID returns ID identifying one or multiple correlated HTTP requests
	RequestID() string
	// UserID returns the user ID of revoke current request
	UserID() string
	// Context embedded for context.Context interface, implemented by routing.Context.Request
	// Be compatible with other context-aware process (grpc, external HTTP calls)
	context.Context
	WithContext(context.Context) RequestScope
	// return context if needed fo an other service
	GetContext() context.Context

	PbApp() *pocketbase.PocketBase
}

type requestScope struct {
	requestID string
	user pb_models.User
	context.Context
	app *pocketbase.PocketBase
}

func (rs *requestScope) UserID() string {
	return rs.user.Id
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func (rs *requestScope) WithContext(ctx context.Context) RequestScope {
	rs2 := new(requestScope)
	*rs2 = *rs
	rs2.Context = ctx
	return rs2
}

func (rs *requestScope) GetContext() context.Context {
	return rs.Context
}

func (rs *requestScope) PbApp() *pocketbase.PocketBase {
	return rs.app
}

func NewRequestScope(requestID string, ctx context.Context, pbApp *pocketbase.PocketBase) RequestScope {
	return &requestScope{
		requestID: requestID,
		Context:   ctx,
		app: pbApp,
	}
}


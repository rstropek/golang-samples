package customersvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/google/uuid"
)

// CustomerMiddleware wraps a CustomerService to surround business logic methods
// with additional functionalities like e.g. logging
type CustomerMiddleware func(CustomerService) CustomerService

// CustomerLoggingMiddleware returns a factory for a logging middleware for the customer service.
func CustomerLoggingMiddleware(logger log.Logger) CustomerMiddleware {
	return func(next CustomerService) CustomerService {
		return &customerLoggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type customerLoggingMiddleware struct {
	next   CustomerService
	logger log.Logger
}

func (mw customerLoggingMiddleware) GetCustomers(ctx context.Context, orderBy string) (c []Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetCustomers", "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetCustomers(ctx, orderBy)
}

func (mw customerLoggingMiddleware) GetCustomer(ctx context.Context, cid uuid.UUID) (c Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetCustomer(ctx, cid)
}

func (mw customerLoggingMiddleware) AddCustomer(ctx context.Context, c Customer) (cust Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "AddCustomer", "id", c.CustomerID.String(), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.AddCustomer(ctx, c)
}

func (mw customerLoggingMiddleware) DeleteCustomer(ctx context.Context, cid uuid.UUID) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteCustomer(ctx, cid)
}

func (mw customerLoggingMiddleware) PatchCustomer(ctx context.Context, cid uuid.UUID, c Customer) (cust Customer, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchCustomer", "id", cid.String(), "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchCustomer(ctx, cid, c)
}

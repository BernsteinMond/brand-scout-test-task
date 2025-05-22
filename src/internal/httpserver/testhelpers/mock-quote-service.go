package testhelpers

import (
	"context"
	"github.com/BernsteinMond/brand-scout-test-task/src/internal/httpserver"
	"github.com/BernsteinMond/brand-scout-test-task/src/internal/service"
	"github.com/google/uuid"
)

type MockQuoteService struct {
	retError error
}

var _ httpserver.QuoteService = (*MockQuoteService)(nil)

func (m *MockQuoteService) CreateNewQuote(context.Context, string, string) error {
	return m.retError
}

func (m *MockQuoteService) GetQuotesWithFilter(context.Context, string) ([]service.Quote, error) {
	if m.retError != nil {
		return nil, m.retError
	}
	return quotesArrayFixture, nil
}

func (m *MockQuoteService) GetRandomQuote(context.Context) (*service.Quote, error) {
	if m.retError != nil {
		return nil, m.retError
	}

	return &quotesArrayFixture[0], nil
}

func (m *MockQuoteService) DeleteQuoteByID(ctx context.Context, id uuid.UUID) error {
	return m.retError
}

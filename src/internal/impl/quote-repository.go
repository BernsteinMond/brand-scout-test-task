package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMond/brand-scout-test-task/src/internal/service"
	"github.com/google/uuid"
)

type QuoteRepository struct {
	db *sql.DB
}

var _ service.QuoteRepository = (*QuoteRepository)(nil)

func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (q *QuoteRepository) CreateNewQuote(ctx context.Context, quote *service.Quote) error {
	const query = `INSERT INTO quote.quotes (id, author, quote) VALUES ($1, $2, $3)`

	res, err := q.db.ExecContext(ctx, query, quote.ID, quote.Author, quote.Quote)
	if err != nil {
		return fmt.Errorf("run sql query: %w", err)
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return service.ErrRepoAlreadyExists
	}

	return nil
}

func (q *QuoteRepository) DeleteQuoteByID(ctx context.Context, id uuid.UUID) error {
	const query = `DELETE FROM quote.quotes WHERE id = $1`

	_, err := q.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("run sql query: %w", err)
	}

	return nil
}

func (q *QuoteRepository) GetQuotesWithFilter(ctx context.Context, authorFilter string) (_ []service.Quote, err error) {
	var query = `SELECT id, author, quote FROM quote.quotes`

	if authorFilter != "" {
		query += " WHERE author = $1"
	}

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}
	defer func() {
		err = errors.Join(err, rows.Close())
	}()

	var (
		ret   = make([]service.Quote, 0)
		quote service.Quote
	)
	for rows.Next() {
		err = rows.Scan(&quote.ID, &quote.Author, &quote.Quote)
		if err != nil {
			return nil, fmt.Errorf("scan into row: %w", err)
		}

		ret = append(ret, quote)
	}

	return ret, nil
}

func (q *QuoteRepository) GetRandomQuote(ctx context.Context) (*service.Quote, error) {
	const query = `SELECT id, author, quote FROM quote.quotes ORDER BY random() LIMIT 1`

	var ret service.Quote

	err := q.db.QueryRowContext(ctx, query).Scan(&ret.ID, &ret.Author, &ret.Quote)
	if err != nil {
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	return &ret, nil
}

func (q *QuoteRepository) GetQuoteByAuthor(ctx context.Context, author string) (*service.Quote, error) {
	const query = `SELECT id,quote FROM quote.quotes WHERE author = $1`

	ret := service.Quote{
		Author: author,
	}

	err := q.db.QueryRowContext(ctx, query, author).Scan(&ret.ID, &ret.Quote)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, service.ErrRepoNotFound
		}
		return nil, fmt.Errorf("run sql query: %w", err)
	}

	return &ret, nil
}

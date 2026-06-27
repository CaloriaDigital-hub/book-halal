package reader

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReaderRepository struct {
	pool *pgxpool.Pool
}

func NewReaderRepository(pool *pgxpool.Pool) *ReaderRepository {
	return &ReaderRepository{pool: pool}
}

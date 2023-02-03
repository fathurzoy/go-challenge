package domain

import (
	"context"
	"time"
)

// Costumer is representing the Costumer data struct
type Costumer struct {
	ID        			int64     `json:"id"`
	CostumerName    string    `json:"costumer_name" validate:"required"`
	CostumerContNo  string    `json:"costumer_cont_no" validate:"required"`
	CostumerAddress string 		`json:"costumer_address"`
	TotalBuy 				string 		`json:"total_buy"`
	CreatorId 			string 		`json:"creator_id"`
	Date 						time.Time `json:"date"`
}

// CostumerUsecase represent the costumer's usecases
type CostumerUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]Costumer, string, error)
	GetByID(ctx context.Context, id int64) (Costumer, error)
	Update(ctx context.Context, ar *Costumer) error
	GetByTitle(ctx context.Context, title string) (Costumer, error)
	Store(context.Context, *Costumer) error
	Delete(ctx context.Context, id int64) error
}

// CostumerRepository represent the costumer's repository contract
type CostumerRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []Costumer, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (Costumer, error)
	GetByTitle(ctx context.Context, title string) (Costumer, error)
	Update(ctx context.Context, ar *Costumer) error
	Store(ctx context.Context, a *Costumer) error
	Delete(ctx context.Context, id int64) error
}

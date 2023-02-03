package usecase

import (
	"context"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type costumerUsecase struct {
	costumerRepo    domain.CostumerRepository
	contextTimeout time.Duration
}

// NewCostumerUsecase will create new an costumerUsecase object representation of domain.CostumerUsecase interface
func NewCostumerUsecase(a domain.CostumerRepository, timeout time.Duration) domain.CostumerUsecase {
	return &costumerUsecase{
		costumerRepo:    a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
func (a *costumerUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.Costumer, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.costumerRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (a *costumerUsecase) GetByID(c context.Context, id int64) (res domain.Costumer, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.costumerRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (a *costumerUsecase) Update(c context.Context, ar *domain.Costumer) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.Date = time.Now()
	return a.costumerRepo.Update(ctx, ar)
}

func (a *costumerUsecase) GetByTitle(c context.Context, title string) (res domain.Costumer, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.costumerRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	return
}

func (a *costumerUsecase) Store(c context.Context, m *domain.Costumer) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedCostumer, _ := a.GetByTitle(ctx, m.CostumerName)
	if existedCostumer != (domain.Costumer{}) {
		return domain.ErrConflict
	}

	err = a.costumerRepo.Store(ctx, m)
	return
}

func (a *costumerUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedCostumer, err := a.costumerRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedCostumer == (domain.Costumer{}) {
		return domain.ErrNotFound
	}
	return a.costumerRepo.Delete(ctx, id)
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/costumer/repository"
	"github.com/bxcodec/go-clean-arch/domain"
)

type postgresCostumerRepository struct {
	Conn *sql.DB
}

// NewPostgresCostumerRepository will create an object that represent the costumer.Repository interface
func NewPostgresCostumerRepository(conn *sql.DB) domain.CostumerRepository {
	return &postgresCostumerRepository{conn}
}

func (m *postgresCostumerRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Costumer, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Costumer, 0)
	for rows.Next() {
		t := domain.Costumer{}
		err = rows.Scan(
			&t.ID,
			&t.CostumerName,
			&t.CostumerContNo,
			&t.CostumerAddress,
			&t.TotalBuy,
			&t.CreatorId,
			&t.Date,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *postgresCostumerRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Costumer, nextCursor string, err error) {
	query := `SELECT id, costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date FROM costumer WHERE date > $1 ORDER BY date LIMIT $2 `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].Date)
	}

	return
}
func (m *postgresCostumerRepository) GetByID(ctx context.Context, id int64) (res domain.Costumer, err error) {
	query := `SELECT id, costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date FROM costumer WHERE ID = $1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Costumer{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *postgresCostumerRepository) GetByTitle(ctx context.Context, title string) (res domain.Costumer, err error) {
	query := `SELECT id, costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date FROM costumer WHERE costumer_name = $1`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *postgresCostumerRepository) Store(ctx context.Context, a *domain.Costumer) (err error) {
	query := `INSERT INTO costumer(costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}
	var id int
	err = stmt.QueryRowContext(ctx, a.CostumerName, a.CostumerContNo, a.CostumerAddress, a.TotalBuy, a.CreatorId, time.Now()).Scan(&id)
	if err != nil {
		return
	}
	a.ID = int64(id)
	return
}

func (m *postgresCostumerRepository) Delete(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM costumer WHERE ID = $1`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *postgresCostumerRepository) Update(ctx context.Context, ar *domain.Costumer) (err error) {
	// query := `INSERT INTO costumer(costumer_name, costumer_cont_no, costumer_address, total_buy, creator_id, date) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	query := `UPDATE costumer set costumer_name=$1, costumer_cont_no=$2, costumer_address=$3, total_buy=$4, creator_id=$5, date=$6 WHERE ID = 7$`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.CostumerName, ar.CostumerContNo, ar.CostumerAddress,  ar.TotalBuy, ar.CreatorId, ar.Date, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

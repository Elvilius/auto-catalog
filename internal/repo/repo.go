package repo

import (
	"context"
	"log/slog"

	"github.com/Elvilius/auto-catalog/domain"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RepoInterface interface {
	CreateCars(ctx context.Context, cars []domain.Car) error
	GetCars(ctx context.Context, filter CarFilter) ([]domain.Car, error)
	DeleteCar(ctx context.Context, ID int) error
	UpdateCar(ctx context.Context, updateCar UpdateCar) error
}

type CarFilter struct {
	RegNum          string
	Mark            string
	Model           string
	YearFrom        int
	YearTo          int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
	Page            int
	PageSize        int
}

type UpdateCar struct {
	ID              int
	RegNum          string
	Mark            string
	Model           string
	Year            int
	OwnerName       string
	OwnerSurname    string
	OwnerPatronymic string
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{db: db}
}

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) CreateCars(ctx context.Context, cars []domain.Car) error {
	conn, err := r.db.Connx(ctx)
	if err != nil {
		slog.Error("[ERROR] failed to connect to db: %v", err)
		return err
	}
	defer conn.Close()

	sl := squirrel.StatementBuilder.Insert("cars").Columns("reg_num", "mark", "model", "year", "owner_name", "owner_surname", "owner_patronymic").PlaceholderFormat(squirrel.Dollar).RunWith(conn)

	for _, car := range cars {
		sl = sl.Values(car.RegNum, car.Mark, car.Model, car.Year, car.OwnerName, car.OwnerSurname, car.OwnerPatronymic)
	}

	sl = sl.Suffix("ON CONFLICT (reg_num) DO NOTHING")
	_, err = sl.ExecContext(ctx)
	if err != nil {
		slog.Error("[ERROR] failed create: %v", err)
		return err
	}
	return nil
}

func (r *Repo) GetCars(ctx context.Context, filter CarFilter) ([]domain.Car, error) {
	cars := make([]domain.Car, 0, filter.PageSize)

	conn, err := r.db.Connx(ctx)
	if err != nil {
		slog.Error("[ERROR] failed to connect to db: %v", err)
		return cars, err
	}
	defer conn.Close()

	carsSl := squirrel.StatementBuilder.
		Select("id", "reg_num", "mark", "model", "year", "owner_name", "owner_surname", "owner_patronymic").From("cars").
		PlaceholderFormat(squirrel.Dollar).
		RunWith(conn)

	carsSl = r.getFilterList(carsSl, filter)

	carsSl = carsSl.Offset(uint64((filter.Page - 1) * filter.PageSize)).Limit(uint64(filter.PageSize))

	rows, err := carsSl.QueryContext(ctx)
	if err != nil {
		slog.Error("[ERROR] failed query: %v", err)
		return cars, err
	}

	for rows.Next() {
		var car domain.Car
		err := rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.OwnerName, &car.OwnerSurname, &car.OwnerPatronymic)
		if err != nil {
			slog.Error("[ERROR] failed find: %v", err)
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

func (r *Repo) DeleteCar(ctx context.Context, ID int) error {
	conn, err := r.db.Connx(ctx)
	if err != nil {
		slog.Error("[ERROR] failed to connect to db: %v", err)
		return err
	}
	defer conn.Close()

	sl := squirrel.StatementBuilder.Delete("cars").PlaceholderFormat(squirrel.Dollar).RunWith(conn).Where(squirrel.Eq{"id": ID})
	_, err = sl.ExecContext(ctx)
	if err != nil {
		slog.Error("[ERROR] failed delete: %v", err)
		return err
	}
	return nil
}

func (r *Repo) UpdateCar(ctx context.Context, updateCar UpdateCar) error {
	conn, err := r.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	sl := squirrel.StatementBuilder.
		Update("cars").
		Set("reg_num", updateCar.RegNum).
		Set("mark", updateCar.Mark).
		Set("model", updateCar.Model).
		Set("year", updateCar.Year).
		Set("owner_name", updateCar.OwnerName).
		Set("owner_surname", updateCar.OwnerSurname).
		Set("owner_patronymic", updateCar.OwnerPatronymic).
		PlaceholderFormat(squirrel.Dollar).
		RunWith(conn).
		Where(squirrel.Eq{"id": updateCar.ID})

	_, err = sl.ExecContext(ctx)
	if err != nil {
		slog.Error("[ERROR] failed update: %v", err)
		return err
	}
	return nil
}

func (r *Repo) getFilterList(sl squirrel.SelectBuilder, filter CarFilter) squirrel.SelectBuilder {
	if filter.RegNum != "" {
		sl = sl.Where(squirrel.Eq{"reg_num": filter.RegNum})
	}

	if filter.Model != "" {
		sl = sl.Where(squirrel.Eq{"model": filter.Model})
	}

	if filter.Mark != "" {
		sl = sl.Where(squirrel.Eq{"mark": filter.Mark})
	}

	if filter.OwnerName != "" {
		sl = sl.Where(squirrel.Eq{"owner_name": filter.OwnerName})
	}

	if filter.OwnerPatronymic != "" {
		sl = sl.Where(squirrel.Eq{"owner_patronymic": filter.OwnerName})
	}

	if filter.OwnerSurname != "" {
		sl = sl.Where(squirrel.Eq{"owner_surname": filter.OwnerSurname})
	}

	if filter.YearFrom != 0 {
		sl = sl.Where(squirrel.GtOrEq{"year": filter.YearFrom})
	}

	if filter.YearTo != 0 {
		sl = sl.Where(squirrel.LtOrEq{"year": filter.YearTo})
	}

	return sl
}

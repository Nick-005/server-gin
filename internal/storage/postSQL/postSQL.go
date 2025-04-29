package sqlite

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	structs "main.go/internal/api/Struct"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func GetAllStatus(storage *sqlx.DB) ([]structs.GetStatus, error) {
	var result []structs.GetStatus
	const op = "storage.sqlite.Get.Status"
	query, args, err := psql.Select("*").From("status").ToSql()
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST TO DB!", op)
		return result, fmt.Errorf("error in creating request to DB")
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, err
	}
	return result, nil
}

func PostNewStatus(storage *sqlx.DB, name string) error {

	query, args, err := psql.Insert("status").Columns("name").Values(name).ToSql()
	if err != nil {
		return err
	}
	_, err = storage.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

// func GetVacancies(storage *sqlx.Tx, limit int, last_id int) ([]Vacancies, error) {
// 	const op = "storage.postgres.Get.Vacancies"
// 	var result []Vacancies
// 	query, args, err := psql.Select("v.id as vacancy_id", "v.emp_id as employee_id", "e.name_organization", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "ex.name as experience", "v.about_work", "v.is_visible as visible", "v.created_at", "v.updated_at").
// 	From("vacancy v").
// 	InnerJoin("employer e ON e.id").
// }

func PostNewExperience(storage *sqlx.DB, name string) error {
	query, args, err := psql.Insert("experience").Columns("name").Values(name).ToSql()
	if err != nil {
		return err
	}
	_, err = storage.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func PostNewEmployer(storage *sqlx.DB, body structs.RequestEmployee) (structs.SuccessEmployer, error) {
	var result structs.SuccessEmployer

	query, args, err := psql.Select("id").From("status").Where(sq.Eq{"name": body.Status}).ToSql()
	if err != nil {
		return result, err
	}
	var status_id int
	err = storage.Get(&status_id, query, args...)
	if err != nil {
		return result, err
	}

	queryMain, argsMain, err := psql.Insert("employer").Columns("name_organization", "phone_number", "email", "inn", "status_id").
		Values(body.NameOrganization, body.PhoneNumber, body.Email, body.INN, status_id).Suffix("RETURNING *").ToSql()
	if err != nil {
		return result, err
	}
	err = storage.Get(&result, queryMain, argsMain...)
	if err != nil {
		return result, err
	}
	return result, nil
}

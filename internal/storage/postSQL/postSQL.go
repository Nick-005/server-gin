package sqlite

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	structs "main.go/internal/api/Struct"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func GetAllCandidates(storage *sqlx.DB) ([]structs.InfoCandidate, error) {
	var result []structs.InfoCandidate

	query, args, err := psql.Select("*").From("candidates").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании запроса на получения данных из таблицы. error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil

}

func PostNewCandidate(storage *sqlx.DB, req structs.RequestCandidate) (structs.InfoCandidate, error) {
	var result structs.InfoCandidate

	query, args, err := psql.Select("id").From("status").Where(sq.Eq{"name": req.UserStatus}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании запроса на получения данных из таблицы. error: %s", err.Error())
	}
	var Ids int

	err = storage.Get(&Ids, query, args...)

	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}

	query, args, err = psql.Insert("candidates").
		Columns("name", "phone_number", "email", "password", "status_id").
		Values(req.Name, req.PhoneNumber, req.Email, req.Password, Ids).
		Suffix("RETURNING *").
		ToSql()

	if err != nil {
		return result, fmt.Errorf("ошибка в формировании запроса на добавление новых данных в таблицу. error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге добавленных данных. error: %s", err.Error())
	}
	return result, nil
}

func PostNewResume(storage *sqlx.DB, req structs.RequestResume) error {

	query, args, err := psql.Select("id").From("experience").Where(sq.Eq{"name": req.Experience}).ToSql()
	if err != nil {
		return err
	}

	var expId int

	err = storage.Get(&expId, query, args...)
	if err != nil {
		return fmt.Errorf("неправильно выбрали опыт. Такого нету в БД. error: %s", err.Error())
	}

	var userID int
	query, args, err = psql.Select("id").From("candidates").Where(sq.Eq{"email": req.UserEmail}).ToSql()
	if err != nil {
		return err
	}

	err = storage.Get(&userID, query, args...)
	if err != nil {
		return fmt.Errorf("неправильно выбрали опыт. Такого нету в БД. error: %s", err.Error())
	}

	MainQuery, MainArgs, err := psql.Insert("resume").Columns("candidate_id", "experience_id", "description").Values(userID, expId, req.Description).ToSql()
	if err != nil {
		return fmt.Errorf("неполучилось сформировать sql скрипты для добавления в БД. error: %s", err.Error())
	}

	_, err = storage.Exec(MainQuery, MainArgs...)
	if err != nil {
		return fmt.Errorf("неполучилось выполнить добавление в БД. error: %s", err.Error())
	}
	return nil
}

func GetAllStatus(storage *sqlx.DB) ([]structs.GetStatus, error) {
	var result []structs.GetStatus
	const op = "storage.sqlite.Get.Status"
	query, args, err := psql.Select("*").From("status").ToSql()
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST TO DB!", op)
		return result, fmt.Errorf("error in creating request to DB. error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
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

func GetAllExperience(storage *sqlx.DB) ([]structs.GetStatus, error) {
	var result []structs.GetStatus
	const op = "storage.sqlite.Get.Experience"
	query, args, err := psql.Select("*").From("experience").ToSql()
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

func GetAllEmployee(storage *sqlx.DB) ([]structs.SuccessEmployer, error) {
	var result []structs.SuccessEmployer

	query, args, err := psql.Select("*").From("employer").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании скрипта запроса. error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, err
	}

	return result, nil
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

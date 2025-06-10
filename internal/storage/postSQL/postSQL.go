package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func GetStatusByName(storage *sqlx.Tx, name string) (s.GetStatus, error) {

	var result s.GetStatus
	query, args, err := psql.Select("*").From("status").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}
	return result, nil
}

func DeleteStatusByName(storage *sqlx.Tx, name string) error {

	query, args, err := psql.Delete("status").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено! Перепроверьте данные и попробуйте снова")
	}
	return nil
}

func GetStatusByID(storage *sqlx.Tx, id int) (s.GetStatus, error) {

	var result s.GetStatus
	query, args, err := psql.Select("*").From("status").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}
	return result, nil
}

func DeleteExperienceByName(storage *sqlx.Tx, name string) error {

	query, args, err := psql.Delete("experience").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}

	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено! Перепроверьте данные и попробуйте снова")
	}

	return nil
}

func GetEmployeeLogin(storage *sqlx.Tx, email, password string) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer
	query, args, err := psql.Select(
		"e.id", "e.name_organization", "e.phone_number", "e.email", "e.inn", "e.password", "e.created_at", "e.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).
		From("employer e").
		Join("status s ON e.status_id = s.id").
		Where(sq.Eq{"email": email, "password": password}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf("неверный логин или пароль. Такого работодателя нету в системе! error: %v", err)
	} else if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetEmployeeByID(storage *sqlx.Tx, emp_id int) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer
	query, args, err := psql.Select(
		"e.id", "e.name_organization", "e.phone_number", "e.email", "e.inn", "e.password", "e.created_at", "e.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).
		From("employer e").
		Join("status s ON e.status_id = s.id").
		Where(sq.Eq{"e.id": emp_id}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetEmployeeByEmail(storage *sqlx.Tx, email string) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer
	query, args, err := psql.Select(
		"e.id", "e.name_organization", "e.phone_number", "e.email", "e.inn", "e.password", "e.created_at", "e.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).
		From("employer e").
		Join("status s ON e.status_id = s.id").
		Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetNumberOfVacancies(storage *sqlx.Tx) (int, error) {
	var number int = -1

	query, args, err := psql.Select("count(id)").From("vacancy").Where(sq.Eq{"is_visible": true}).ToSql()

	if err != nil {
		return number, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Get(&number, query, args...)
	if err != nil {
		return number, fmt.Errorf("ошибка в маппинге данных вакансии ! error: %s", err.Error())
	}

	return number, nil
}

func GetVacancyByID(storage *sqlx.Tx, id int) (s.VacancyData, error) {
	var result s.VacancyData

	query, args, err := psql.Select(
		"v.id", "v.name", "v.email", "v.price", "v.phone_number", "v.location", "v.about_work", "v.is_visible", "v.created_at", "v.updated_at",
		"e.id as \"experience.id\"", "e.name as \"experience.name\"", "e.created_at as \"experience.created_at\"",
	).
		From("vacancy v").
		Join("experience e ON v.experience_id = e.id").Where(sq.Eq{
		"v.id": id,
	}).
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		return result, err
	} else if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных вакансии ! error: %s", err.Error())
	}
	return result, nil
}

func GetAllVacanciesByEmployee(storage *sqlx.Tx, emp_id int) ([]s.VacancyData, error) {
	var result []s.VacancyData

	query, args, err := psql.Select(
		"v.id", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "v.about_work", "v.is_visible", "v.created_at", "v.updated_at",
		"e.id as \"experience.id\"", "e.name as \"experience.name\"", "e.created_at as \"experience.created_at\"",
	).
		From("vacancy v").
		Join("experience e ON v.experience_id = e.id").Where(sq.Eq{
		"v.emp_id": emp_id,
	}).OrderBy("id ASC").
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных резюме ! error: %s", err.Error())
	}
	return result, nil
}

func GetVacancyInfoByID(storage *sqlx.Tx, vac_id int) (s.VacancyData_Limit, error) {
	var result s.VacancyData_Limit
	query, args, err := psql.Select(
		"v.id", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "v.about_work", "v.is_visible", "v.created_at", "v.updated_at",

		"e.id as \"experience.id\"", "e.name as \"experience.name\"", "e.created_at as \"experience.created_at\"",

		"em.id as \"employer.id\"", "em.name_organization as \"employer.name_organization\"",
		"em.phone_number as \"employer.phone_number\"", "em.email as \"employer.email\"",
		"em.inn as \"employer.inn\"", "em.password as \"employer.password\"",
		"em.created_at as \"employer.created_at\"", "em.updated_at as \"employer.updated_at\"",

		"s.id as \"employer.status.id\"", "s.name as \"employer.status.name\"", "s.created_at as \"employer.status.created_at\"",
	).From("vacancy v").
		Join("experience e ON v.experience_id = e.id").
		Join("employer em ON v.emp_id = em.id").
		Join("status s ON em.status_id = s.id").OrderBy("v.id ASC").
		Where(sq.Eq{"v.id": vac_id}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		return result, err
	} else if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных вакансий! error: %s", err.Error())
	}
	return result, nil
}

func GetVacancyLimitByTimes(storage *sqlx.Tx, limit int, time time.Time) ([]s.VacancyData_Limit, error) {
	var result []s.VacancyData_Limit
	query, args, err := psql.Select(
		"v.id", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "v.about_work", "v.is_visible", "v.created_at", "v.updated_at",

		"e.id as \"experience.id\"", "e.name as \"experience.name\"", "e.created_at as \"experience.created_at\"",

		"em.id as \"employer.id\"", "em.name_organization as \"employer.name_organization\"",
		"em.phone_number as \"employer.phone_number\"", "em.email as \"employer.email\"",
		"em.inn as \"employer.inn\"", "em.password as \"employer.password\"",
		"em.created_at as \"employer.created_at\"", "em.updated_at as \"employer.updated_at\"",

		"s.id as \"employer.status.id\"", "s.name as \"employer.status.name\"", "s.created_at as \"employer.status.created_at\"",
	).From("vacancy v").
		Join("experience e ON v.experience_id = e.id").
		Join("employer em ON v.emp_id = em.id").
		Join("status s ON em.status_id = s.id").OrderBy("v.id ASC").
		Where(sq.Gt{"v.created_at": time}).Limit(uint64(limit)).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных вакансий! error: %s", err.Error())
	}
	return result, nil
}

func GetVacancyLimit(storage *sqlx.Tx, page, perPage int) ([]s.VacancyData_Limit, error) {
	var result []s.VacancyData_Limit

	offset := (page - 1) * perPage

	queryBuilder := psql.Select(
		"v.id", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "v.about_work", "v.is_visible", "v.created_at", "v.updated_at",

		"e.id as \"experience.id\"", "e.name as \"experience.name\"", "e.created_at as \"experience.created_at\"",

		"em.id as \"employer.id\"", "em.name_organization as \"employer.name_organization\"",
		"em.phone_number as \"employer.phone_number\"", "em.email as \"employer.email\"",
		"em.inn as \"employer.inn\"", "em.password as \"employer.password\"",
		"em.created_at as \"employer.created_at\"", "em.updated_at as \"employer.updated_at\"",

		"s.id as \"employer.status.id\"", "s.name as \"employer.status.name\"", "s.created_at as \"employer.status.created_at\"",
	).From("vacancy v").
		Join("experience e ON v.experience_id = e.id").
		Join("employer em ON v.emp_id = em.id").
		Join("status s ON em.status_id = s.id").OrderBy("v.id ASC").
		Where(sq.Eq{"v.is_visible": true}).Limit(uint64(perPage)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()

	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных вакансий! error: %s", err.Error())
	}
	return result, nil
}

func PostNewVacancy(storage *sqlx.Tx, req s.ResponseVac, emp_id int) (s.VacancyData, error) {
	var result s.VacancyData

	query, args, err := psql.Insert("vacancy").
		Columns(
			"emp_id",
			"name",
			"price",
			"email",
			"phone_number",
			"location",
			"experience_id",
			"about_work",
			"is_visible",
		).Values(
		emp_id,
		req.VacancyName,
		req.Price,
		req.Email,
		req.PhoneNumber,
		req.Location,
		req.ExperienceId,
		req.About,
		req.IsVisible,
	).Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании запроса на добавление новых данных в таблицу. error: %s", err.Error())
	}
	var id int
	err = storage.Get(&id, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге добавленных данных. error: %s", err.Error())
	}

	result, err = GetVacancyByID(storage, id)
	if err != nil {
		return result, err
	}

	return result, nil
}

func UpdateVacancyInfo(storage *sqlx.Tx, req s.VacancyPut, uid int) error {

	query, args, err := psql.Update("vacancy").
		Set("name", req.VacancyName).
		Set("price", req.Price).
		Set("email", req.Email).
		Set("phone_number", req.PhoneNumber).
		Set("location", req.Location).
		Set("experience_id", req.ExperienceId).
		Set("about_work", req.About).
		Set("is_visible", req.IsVisible).
		Where(sq.Eq{"id": req.ID, "emp_id": uid}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для обновления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("данные не были обновлены, так как обновляемой вакансии не было найдено! Перепроверьте данные и попробуйте снова")
	}

	return nil
}

func UpdateCandidateInfo(storage *sqlx.Tx, req s.RequestCandidate, id int) error {
	var args []interface{}
	var query string
	var err error

	if len(req.Password) <= 3 {
		query, args, err = psql.Update("candidates").
			Set("name", req.Name).
			Set("phone_number", req.PhoneNumber).
			Set("email", req.Email).
			// Set("password", req.Password).
			Set("status_id", req.Status_id).
			Where(sq.Eq{"id": id}).
			ToSql()
	} else {
		query, args, err = psql.Update("candidates").
			Set("name", req.Name).
			Set("phone_number", req.PhoneNumber).
			Set("email", req.Email).
			Set("password", req.Password).
			Set("status_id", req.Status_id).
			Where(sq.Eq{"id": id}).
			ToSql()
	}

	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для обновления данных! error: %s", err.Error())
	}
	_, err = storage.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func GetResponseByCandidate(storage *sqlx.Tx, uid int) ([]s.ResponseByVac, error) {
	var result []s.ResponseByVac

	query, args, err := psql.Select(
		"r.id",
		"v.id as \"vacancy.id\"",
		"v.name as \"vacancy.name\"",
		"v.price as \"vacancy.price\"",
		"v.email as \"vacancy.email\"",
		"v.phone_number as \"vacancy.phone_number\"",
		"v.location as \"vacancy.location\"",
		"v.about_work as \"vacancy.about_work\"",
		"v.is_visible as \"vacancy.is_visible\"",
		"v.created_at as \"vacancy.created_at\"",
		"v.updated_at as \"vacancy.updated_at\"",
		"em.name_organization as \"vacancy.employee_name\"",
		"ex.id as \"vacancy.experience.id\"",
		"ex.name as \"vacancy.experience.name\"",
		"ex.created_at as \"vacancy.experience.created_at\"",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).
		From("response r").
		Join("vacancy v ON r.vacancy_id = v.id").
		Join("experience ex ON v.experience_id = ex.id").
		Join("status s ON r.status_id = s.id").
		Join("employer em ON v.emp_id = em.id").
		Where(sq.Eq{"r.candidates_id": uid}).OrderBy("r.id ASC").
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка при выполнении скрипта на получение данных. error: %s", err.Error())
	}
	return result, nil
}

func GetResponseOnVacancy(storage *sqlx.Tx, uid, vac_id int) (s.ResponseOnVacancy, error) {
	var result s.ResponseOnVacancy
	result.IsResponsed = false
	query, args, err := psql.Select(
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
		// "c.id as \"candidate.id\"", "c.name as \"candidate.name\"", "c.phone_number as \"candidate.phone_number\"", "c.email as \"candidate.email\"",
		// "c.password as \"candidate.password\"", "c.created_at as \"candidate.created_at\"", "c.updated_at as \"candidate.updated_at\"",
		// "s2.id as \"candidate.status.id\"", "s2.name as \"candidate.status.name\"", "s2.created_at as \"candidate.status.created_at\"",

	).
		From("response r").
		Join("status s ON r.status_id = s.id").
		Where(sq.Eq{"r.candidates_id": uid, "r.vacancy_id": vac_id}).
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}
	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		result.IsResponsed = false
		return result, nil
	} else if err != nil {
		return result, fmt.Errorf("ошибка при выполнении скрипта на получения данных. error: %s", err.Error())
	}
	result.IsResponsed = true
	return result, nil
}

func AuthorizationMethodForUsers(storage *sqlx.Tx, email, password string) {

}

func CheckEmailIsValid(storage *sqlx.Tx, email string) (bool, error) {
	var amnt_emp int = -1
	var amnt_cnd int = -1
	query, args, err := psql.Select("count(id)").From("employer").Where(sq.Eq{"email": email}).ToSql()

	if err != nil {
		return false, err
	}

	err = storage.Get(&amnt_emp, query, args...)
	if err == sql.ErrNoRows {
		amnt_emp = 0
	} else if err != nil {
		return false, fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}

	query, args, err = psql.Select("count(id)").From("candidates").Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return false, err
	}
	err = storage.Get(&amnt_cnd, query, args...)
	if err == sql.ErrNoRows {
		amnt_cnd = 0
	} else if err != nil {
		return false, fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}

	if amnt_cnd+amnt_emp != 0 {
		return false, fmt.Errorf("такой email уже используется! Выберите другой и попробуйте снова")
	}
	return true, nil
}

func CheckUserByEmailOnEmployer(storage *sqlx.Tx, email string) (bool, error) {
	var amnt_emp int = -1
	query, args, err := psql.Select("count(id)").From("employer").Where(sq.Eq{"email": email}).ToSql()

	if err != nil {
		return false, err
	}

	err = storage.Get(&amnt_emp, query, args...)
	if amnt_emp == 0 {
		return false, nil

	} else if err != nil {
		return false, fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}
	return true, nil
}

func GetResponseByVacancy(storage *sqlx.Tx, vac_id int) (s.SuccessResponse, error) {
	var result s.SuccessResponse

	query, args, err := psql.Select(
		"r.id", "r.created_at",
		"c.id as \"candidate.id\"", "c.name as \"candidate.name\"", "c.phone_number as \"candidate.phone_number\"", "c.email as \"candidate.email\"",
		"c.password as \"candidate.password\"", "c.created_at as \"candidate.created_at\"", "c.updated_at as \"candidate.updated_at\"",
		"s2.id as \"candidate.status.id\"", "s2.name as \"candidate.status.name\"", "s2.created_at as \"candidate.status.created_at\"",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).
		From("response r").
		Join("candidates c ON r.candidates_id = c.id").
		Join("status s2 ON c.status_id = s2.id").
		Join("status s ON r.status_id = s.id  ").
		Where(sq.Eq{"r.vacancy_id": vac_id}).OrderBy("r.created_at ASC").
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}
	err = storage.Select(&result.Responses, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}
	return result, nil
}

func DeleteVacancy(storage *sqlx.Tx, uid, id int, role string) error {

	if role == "ADMIN" {
		query, args, err := psql.Delete("vacancy").Where(sq.Eq{"id": id}).ToSql()
		if err != nil {
			return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
		}
		result, err := storage.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("таких записей не было найдено у данного пользователя! Перепроверьте данные и попробуйте снова")
		}
	} else {
		query, args, err := psql.Delete("vacancy").Where(sq.Eq{"id": id, "emp_id": uid}).ToSql()
		if err != nil {
			return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
		}
		result, err := storage.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
		}
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("таких записей не было найдено у данного пользователя! Перепроверьте данные и попробуйте снова")
		}
	}
	return nil
}

func PostResponse(storage *sqlx.Tx, id, vac_id int) (int, error) {
	var res_id int

	query, args, err := psql.Insert("response").Columns("candidates_id", "vacancy_id", "status_id").
		Values(id, vac_id, 3).Suffix("ON CONFLICT (candidates_id, vacancy_id) DO NOTHING RETURNING id").ToSql()
	if err != nil {
		return -1, fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}
	err = storage.Get(&res_id, query, args...)
	if err == sql.ErrNoRows {
		return -1, fmt.Errorf("вы уже откликались на эту вакансию! error: %v", err)
	} else if err != nil {
		return -1, fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}
	return res_id, nil
}

func PatchVisibilityVacancy(storage *sqlx.Tx, vacID, empID int, visible bool) error {

	query, args, err := psql.Update("vacancy").Set("is_visible", visible).Where(sq.Eq{"id": vacID, "emp_id": empID}).ToSql()
	if err != nil {
		return err
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("данные не были обновлены, так как обновляемой вакансии не было найдено! Перепроверьте данные и попробуйте снова")
	}
	return nil
}

func PatchResponse(storage *sqlx.Tx, req s.ResponsePatch) error {

	query, args, err := psql.Update("response").
		Set("status_id", req.Status_id).
		Where(sq.Eq{"id": req.Response_id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для обновления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("данные не были обновлены, так как обновляемого отклика не было найдено! Перепроверьте данные и попробуйте снова")
	}

	return nil
}

func GetCandidateById(storage *sqlx.Tx, id int) (s.InfoCandidate, error) {
	var result s.InfoCandidate

	query, args, err := psql.Select(
		"c.id", "c.name", "c.phone_number", "c.email", "c.password", "c.created_at", "c.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("candidates c").Join("status s ON c.status_id = s.id").
		Where(sq.Eq{"c.id": id}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf("такого пользователя не было найдено! error: %s", err.Error())
	} else if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetCandidateByLogin(storage *sqlx.Tx, email, password string) (s.InfoCandidate, error) {
	var result s.InfoCandidate

	query, args, err := psql.Select(
		"c.id", "c.name", "c.phone_number", "c.email", "c.password", "c.created_at", "c.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("candidates c").Join("status s ON c.status_id = s.id").
		Where(sq.Eq{"email": email, "password": password}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err == sql.ErrNoRows {
		return result, fmt.Errorf("Неверный логин или пароль. Такого соискателя нету в системе! error: %v", err)
	} else if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetCandidateByEmail(storage *sqlx.Tx, email string) (s.InfoCandidate, error) {
	var result s.InfoCandidate

	query, args, err := psql.Select(
		"c.id", "c.name", "c.phone_number", "c.email", "c.password", "c.created_at", "c.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("candidates c").Join("status s ON c.status_id = s.id").
		Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func DeleteResponse(storage *sqlx.Tx, id, uid int) error {
	query, args, err := psql.Delete("response").Where(sq.Eq{"vacancy_id": id, "candidates_id": uid}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено у данного пользователя! Перепроверьте данные и попробуйте снова")
	}
	return nil
}

func DeleteResume(storage *sqlx.Tx, id, uid int) error {

	query, args, err := psql.Delete("resume").Where(sq.Eq{"id": id, "candidate_id": uid}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено у данного пользователя! Перепроверьте данные и попробуйте снова")
	}
	return nil
}

func UpdateCandidateResume(storage *sqlx.Tx, req s.RequestResumeUpdate, uid int) error {

	query, args, err := psql.Update("resume").
		Set("experience_id", req.Experience).
		Set("description", req.Description).
		Where(sq.Eq{"id": req.Resume_id, "candidate_id": uid}).
		ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для обновления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("данные не были обновлены, так как обновляемого резюме не было найдено! Перепроверьте данные и попробуйте снова")
	}

	return nil
}

func GetAllResumeByCandidate(storage *sqlx.Tx, id int) (s.ResumeResult, error) {
	var result s.ResumeResult
	// ~ candidates
	query, args, err := psql.Select(
		"c.id",
		"c.name",
		"c.phone_number",
		"c.email ",
		"c.password",
		"c.created_at ",
		"c.updated_at ",
		// ! status
		"s.id as \"status.id\"",
		"s.name as \"status.name\"",
		"s.created_at as \"status.created_at\"",
	).From("candidates c").
		Join("status s ON c.status_id = s.id").
		Where(sq.Eq{"c.id": id}).OrderBy("c.id ASC").
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Get(&result.Candidate, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных кандидата! error: %s", err.Error())
	}

	// ~ resume
	query, args, err = psql.Select(
		"r.id ",
		"r.description ",
		"r.created_at ",
		"r.updated_at",

		"ex.id as \"experience.id\"",
		"ex.name as \"experience.name\"",
		"ex.created_at as \"experience.created_at\"",
	).
		From("resume r").
		Join("experience ex ON r.experience_id = ex.id").
		Where(sq.Eq{"r.candidate_id": id}).OrderBy("r.id ASC").
		ToSql()

	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Select(&result.Resumes, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных резюме ! error: %s", err.Error())
	}
	return result, nil
}

// нормализовал, всё норм
func GetAllCandidates(storage *sqlx.Tx) ([]s.InfoCandidate, error) {
	var result []s.InfoCandidate

	query, args, err := psql.Select(
		"c.id", "c.name", "c.phone_number", "c.email", "c.password", "c.created_at", "c.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("candidates c").Join("status s ON c.status_id = s.id").OrderBy("c.id ASC").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil

}

func DeleteCandidate(storage *sqlx.Tx, uid int) error {
	query, args, err := psql.Delete("candidates").Where(sq.Eq{"id": uid}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено! Перепроверьте данные и попробуйте снова")
	}
	return nil

}

func PostNewCandidate(storage *sqlx.Tx, req s.RequestCandidate) (s.InfoCandidate, error) {
	var result s.InfoCandidate

	query, args, err := psql.Insert("candidates").
		Columns("name", "phone_number", "email", "password", "status_id").
		Values(req.Name, req.PhoneNumber, req.Email, req.Password, req.Status_id).
		ToSql()

	if err != nil {
		return result, fmt.Errorf("ошибка в формировании запроса на добавление новых данных в таблицу. error: %s", err.Error())
	}

	_, err = storage.Exec(query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге добавленных данных. error: %s", err.Error())
	}

	result, err = GetCandidateByEmail(storage, req.Email)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении добавленных данных. error: %s", err.Error())
	}
	return result, nil
}

func PostNewResume(storage *sqlx.Tx, req s.RequestResume, userID int) error {

	MainQuery, MainArgs, err := psql.Insert("resume").Columns("candidate_id", "experience_id", "description").Values(userID, req.Experience, req.Description).ToSql()
	if err != nil {
		return fmt.Errorf("неполучилось сформировать sql скрипты для добавления в БД. error: %s", err.Error())
	}

	_, err = storage.Exec(MainQuery, MainArgs...)
	if err != nil {
		return fmt.Errorf("неполучилось выполнить добавление в БД. error: %s", err.Error())
	}
	return nil
}

func GetAllStatus(storage *sqlx.Tx) ([]s.GetStatus, error) {
	var result []s.GetStatus
	query, args, err := psql.Select("*").From("status").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil
}

func PostNewStatus(storage *sqlx.Tx, name string) error {

	query, args, err := psql.Insert("status").Columns("name").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}
	_, err = storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}
	return nil
}

func GetExperienceByName(storage *sqlx.Tx, name string) (s.GetStatus, error) {
	var result s.GetStatus

	query, args, err := psql.Select("*").From("experience").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil
}

func GetAllExperience(storage *sqlx.Tx) ([]s.GetStatus, error) {
	var result []s.GetStatus

	query, args, err := psql.Select("*").From("experience").ToSql()
	if err != nil {

		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil
}

func PostNewExperience(storage *sqlx.Tx, name string) error {
	query, args, err := psql.Insert("experience").Columns("name").Values(name).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}
	_, err = storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении скрипта на добавления данных. error: %s", err.Error())
	}
	return nil
}

func GetAllEmployee(storage *sqlx.Tx) ([]s.SuccessEmployer, error) {
	var result []s.SuccessEmployer

	query, args, err := psql.Select(
		"em.id", "em.name_organization", "em.phone_number", "em.password", "em.email", "em.inn", "em.created_at", "em.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("employer em").Join("status s ON em.status_id = s.id").OrderBy("em.id ASC").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании скрипта запроса. error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}

	return result, nil
}

func DeleteEmployee(storage *sqlx.Tx, uid int) error {
	query, args, err := psql.Delete("employer").Where(sq.Eq{"id": uid}).ToSql()
	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для удаления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка в исполнении SQL скрипта на удаление! error: %s", err.Error())
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("таких записей не было найдено! Перепроверьте данные и попробуйте снова")
	}
	return nil
}

func PostNewEmployer(storage *sqlx.Tx, body s.RequestEmployee) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer

	queryMain, argsMain, err := psql.Insert("employer").Columns("name_organization", "phone_number", "email", "inn", "password", "status_id").
		Values(body.NameOrganization, body.PhoneNumber, body.Email, body.INN, body.Password, body.Status_id).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании скрипта запроса. error: %s", err.Error())
	}
	_, err = storage.Exec(queryMain, argsMain...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}

	result, err = GetEmployeeByEmail(storage, body.Email)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении добавленных данных. error: %s", err.Error())
	}
	return result, nil
}

func UpdateEmployeeInfo(storage *sqlx.Tx, req s.RequestEmployer, uid int) error {
	var args []interface{}
	var query string
	var err error
	if len(req.Password) <= 3 {
		query, args, err = psql.Update("employer").
			Set("name_organization", req.NameOrganization).
			Set("phone_number", req.PhoneNumber).
			Set("email", req.Email).
			// Set("password", req.Password).
			Set("status_id", req.Status_id).
			Where(sq.Eq{"id": uid}).
			ToSql()
	} else {
		query, args, err = psql.Update("employer").
			Set("name_organization", req.NameOrganization).
			Set("phone_number", req.PhoneNumber).
			Set("email", req.Email).
			Set("password", req.Password).
			Set("status_id", req.Status_id).
			Where(sq.Eq{"id": uid}).
			ToSql()
	}

	if err != nil {
		return fmt.Errorf("ошибка в создании SQL скрипта для обновления данных! error: %s", err.Error())
	}
	result, err := storage.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("данные не были обновлены, так как работодатель не был найден! Перепроверьте данные и попробуйте снова")
	}

	return nil
}

func CreateAccessToken(claim *s.Claims) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	result, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_TOKEN_EMP")))
	if err != nil {
		return "error", err
	}
	result = strings.ReplaceAll(result, "+", "-")
	result = strings.ReplaceAll(result, "/", "_")
	return result, nil
}

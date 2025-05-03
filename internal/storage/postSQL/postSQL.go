package sqlite

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func GetAllVacanciesByEmployee(storage *sqlx.DB, id int) {

}

func GetStatusByName(storage *sqlx.DB, name string) (s.GetStatus, error) {

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

func GetEmployeeByEmail(storage *sqlx.DB, email string) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer
	query, args, err := psql.Select("*").From("employer").Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func PostNewVacancy(storage *sqlx.DB, req s.ResponseVac) (s.Vacancies, s.SuccessEmployer, s.GetStatus, error) {
	var vac s.Vacancies
	var emf s.SuccessEmployer

	exp, err := GetExperienceByName(storage, req.Experience)
	if err != nil {
		return vac, emf, exp, fmt.Errorf("ошибка в получении данных из таблицы experience! error: %s", err.Error())
	}

	employee, err := GetEmployeeByEmail(storage, req.Emp_Email)
	if err != nil {
		return vac, employee, exp, fmt.Errorf("ошибка в получении данных из таблицы employee! error: %s", err.Error())
	}

	// дописать
	query, args, err := psql.Insert("vacancy").
		Columns("emp_id", "name", "price", "email", "phone_number", "location", "experience_id", "about_work", "is_visible").
		Values(employee.ID, req.Vac_Name, req.Price, req.Email, req.PhoneNumber, req.Location, exp.ID, req.About, req.Is_visible).
		Suffix("RETURNING *").ToSql()
	if err != nil {
		return vac, employee, exp, fmt.Errorf("ошибка в создании SQL скрипта для добавления данных! error: %s", err.Error())
	}

	err = storage.Get(&vac, query, args...)
	if err != nil {
		return vac, employee, exp, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return vac, employee, exp, nil
}

func GetCandidateById(storage *sqlx.DB, id int) (s.InfoCandidate, error) {
	var result s.InfoCandidate

	query, args, err := psql.Select("*").From("candidates").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Get(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных! error: %s", err.Error())
	}

	return result, nil
}

func GetAllResumeByCandidate(storage *sqlx.DB, id int) (s.ResumeResult, error) {
	var result s.ResumeResult

	// query, args, err := psql.Select(
	// ! resume
	// 	"r.id", "r.description", "r.created_at", "r.updated_at",
	// ! candidates
	// "c.id as \"candidate.id\"",
	// "c.name as \"candidate.name\"",
	// "c.phone_number as \"candidate.phone_number\"",
	// "c.email as \"candidate.email\"",
	// "c.password as \"candidate.password\"",
	// "c.created_at as \"candidate.created_at\"",
	// "c.updated_at as \"candidate.updated_at\"",
	// ! status
	// "s.id as \"candidate.status.id\"",
	// "s.name as \"candidate.status.name\"",
	// "s.created_at as \"candidate.status.created_at\"",
	// ! experience
	// 	"ex.id as \"experience.id\"",
	// 	"ex.name as \"experience.name\"",
	// 	"ex.created_at as \"experience.created_at\"",
	// ).From("resume r").
	// 	Join("candidates c ON r.candidate_id = c.id").
	// 	Join("status s ON c.status_id = s.id").
	// 	Join("experience ex ON r.experience_id = ex.id").
	// 	Where(sq.Eq{"r.candidate_id": id}).ToSql()
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
		Where(sq.Eq{"c.id": id}).
		ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	// fmt.Println(query)
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
		Where(sq.Eq{"r.candidate_id": id}).
		ToSql()

	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}
	// fmt.Println(query)
	err = storage.Select(&result.Resumes, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в маппинге данных резюме ! error: %s", err.Error())
	}
	return result, nil
}

// нормализовал, всё норм
func GetAllCandidates(storage *sqlx.DB) ([]s.InfoCandidate, error) {
	var result []s.InfoCandidate

	query, args, err := psql.Select(
		"c.id", "c.name", "c.phone_number", "c.email", "c.password", "c.created_at", "c.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("candidates c").Join("status s ON c.status_id = s.id").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в создании SQL скрипта для получения данных! error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}
	return result, nil

}

func PostNewCandidate(storage *sqlx.DB, req s.RequestCandidate) (s.InfoCandidate, error) {
	var result s.InfoCandidate

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

func PostNewResume(storage *sqlx.DB, req s.RequestResume) error {

	exp, err := GetExperienceByName(storage, req.Experience)
	if err != nil {
		return err
	}
	var userID int
	query, args, err := psql.Select("id").From("candidates").Where(sq.Eq{"email": req.UserEmail}).ToSql()
	if err != nil {
		return err
	}

	err = storage.Get(&userID, query, args...)
	if err != nil {
		return fmt.Errorf("неправильно выбрали опыт. Такого нету в БД. error: %s", err.Error())
	}

	MainQuery, MainArgs, err := psql.Insert("resume").Columns("candidate_id", "experience_id", "description").Values(userID, exp.ID, req.Description).ToSql()
	if err != nil {
		return fmt.Errorf("неполучилось сформировать sql скрипты для добавления в БД. error: %s", err.Error())
	}

	_, err = storage.Exec(MainQuery, MainArgs...)
	if err != nil {
		return fmt.Errorf("неполучилось выполнить добавление в БД. error: %s", err.Error())
	}
	return nil
}

func GetAllStatus(storage *sqlx.DB) ([]s.GetStatus, error) {
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

func PostNewStatus(storage *sqlx.DB, name string) error {

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

// func GetVacancies(storage *sqlx.Tx, limit int, last_id int) ([]Vacancies, error) {
// 	const op = "storage.postgres.Get.Vacancies"
// 	var result []Vacancies
// 	query, args, err := psql.Select("v.id as vacancy_id", "v.emp_id as employee_id", "e.name_organization", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "ex.name as experience", "v.about_work", "v.is_visible as visible", "v.created_at", "v.updated_at").
// 	From("vacancy v").
// 	InnerJoin("employer e ON e.id").
// }

func GetExperienceByName(storage *sqlx.DB, name string) (s.GetStatus, error) {
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

func GetAllExperience(storage *sqlx.DB) ([]s.GetStatus, error) {
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

func PostNewExperience(storage *sqlx.DB, name string) error {
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

func GetAllEmployee(storage *sqlx.DB) ([]s.SuccessEmployer, error) {
	var result []s.SuccessEmployer

	query, args, err := psql.Select(
		"em.id", "em.name_organization", "em.phone_number", "em.email", "em.inn", "em.created_at", "em.updated_at",
		"s.id as \"status.id\"", "s.name as \"status.name\"", "s.created_at as \"status.created_at\"",
	).From("employer em").Join("status s ON em.status_id = s.id").ToSql()
	if err != nil {
		return result, fmt.Errorf("ошибка в формировании скрипта запроса. error: %s", err.Error())
	}

	err = storage.Select(&result, query, args...)
	if err != nil {
		return result, fmt.Errorf("ошибка в получении и маппинге данных. error: %s", err.Error())
	}

	return result, nil
}

func PostNewEmployer(storage *sqlx.DB, body s.RequestEmployee) (s.SuccessEmployer, error) {
	var result s.SuccessEmployer

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

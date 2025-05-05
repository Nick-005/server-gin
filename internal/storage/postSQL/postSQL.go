package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	s "main.go/internal/api/Struct"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func GetAllVacanciesByEmployee(storage *sqlx.Tx, id int) {

}

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
		return result, err
	} else if err != nil {
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

// TODO переделать нахуй это. Что за хуйня тут...
func PostNewVacancy(storage *sqlx.Tx, req s.ResponseVac) (s.Vacancies, s.SuccessEmployer, s.GetStatus, error) {
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

func GetCandidateById(storage *sqlx.Tx, id int) (s.InfoCandidate, error) {
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
		return result, err
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
		Where(sq.Eq{"c.id": id}).
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
		Where(sq.Eq{"r.candidate_id": id}).
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

// func GetVacancies(storage *sqlx.Tx, limit int, last_id int) ([]Vacancies, error) {
// 	const op = "storage.postgres.Get.Vacancies"
// 	var result []Vacancies
// 	query, args, err := psql.Select("v.id as vacancy_id", "v.emp_id as employee_id", "e.name_organization", "v.name", "v.price", "v.email", "v.phone_number", "v.location", "ex.name as experience", "v.about_work", "v.is_visible as visible", "v.created_at", "v.updated_at").
// 	From("vacancy v").
// 	InnerJoin("employer e ON e.id").
// }

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

// func DecodingToken(tokenString string) (string, error) {
// 	var
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte(secretKEY), nil
// 	})

// 	return token, nil
// }

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

// func CreateToken(email, user string) (string, error) {
// 	type Header struct {
// 		Alg string `json:"alg"` // Алгоритм подписи
// 		Typ string `json:"typ"` // Тип токена
// 	}

// 	type Payload struct {
// 		Iss string `json:"iss"`
// 		Sub string `json:"sub"` // Subject (обычно идентификатор пользователя)
// 		Iat int64  `json:"iat"` // Issued at - время в которое был выдан токен
// 		Exp int64  `json:"exp"` // Время истечения токена (в Unix timestamp)
// 	}
// 	var secretKEY string
// 	if user == "emp" {
// 		secretKEY = os.Getenv("JWT_SECRET_TOKEN_EMP")
// 	} else {
// 		secretKEY = os.Getenv("JWT_SECRET_TOKEN_USER")
// 	}

// 	var header Header
// 	header.Alg = "HS256"
// 	header.Typ = "JWT"

// 	var payload Payload
// 	payload.Iss = "Nick005-aka-monkeyZV"
// 	payload.Sub = email
// 	payload.Iat = time.Now().Unix()
// 	payload.Exp = time.Now().Add(time.Minute * 15).Unix()

// 	headerJSON, err := json.Marshal(header)
// 	if err != nil {
// 		return "error", fmt.Errorf("error in converting HEADER to JSON")
// 	}
// 	headerBASE64 := base64.RawURLEncoding.Strict().EncodeToString(headerJSON)

// 	payloadJSON, err := json.Marshal(payload)
// 	if err != nil {
// 		return "error", fmt.Errorf("error in converting PAYLOAD to JSON")
// 	}
// 	payloadBASE64 := base64.RawURLEncoding.Strict().EncodeToString(payloadJSON)

// 	// создаем подпись для JWTшки
// 	signaturePayAndHeader := fmt.Sprintf("%s.%s", headerBASE64, payloadBASE64)

// 	h := hmac.New(sha256.New, []byte(secretKEY))
// 	h.Write([]byte(signaturePayAndHeader))
// 	var signature string = base64.RawStdEncoding.EncodeToString(h.Sum(nil))

// 	var tokenJWT string = fmt.Sprintf("%s.%s.%s", headerBASE64, payloadBASE64, signature)

// 	return tokenJWT, nil
// }

// func CreateAccessToken(email, user string) (string, error) {
// 	const op = "sqlite.CreateAccessToken.User"
// 	token, err := CreateToken(email, user)
// 	if err != nil {
// 		return "Error", fmt.Errorf("%s: %w", op, err)
// 	}
// 	token = strings.ReplaceAll(token, "+", "-")
// 	token = strings.ReplaceAll(token, "/", "_")
// 	return token, nil
// }

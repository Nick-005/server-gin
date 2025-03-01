package sqlite

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type RequestEmployee struct {
	ID               int    `json:"ID"`
	Limit            int    `json:"limit"`
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	Geography        string `json:"geography"`
	About            string `json:"about"`
}

type ResponseVac struct {
	ID         int    `json:"ID"`
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}

func CreateTokenTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Token"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtEmp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS token(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		active_token TEXT NOT NULL, 
		is_active INTEGER NOT NULL CHECK (is_active IN (0,1))
		);
		CREATE INDEX IF NOT EXISTS about ON token(user_id);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtEmp.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func CreateEmployeeTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Emp"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtEmp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS employee(
		id INTEGER PRIMARY KEY,
		limitVac INTEGER,
		nameOrganization TEXT NOT NULL UNIQUE,
		phoneNumber TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE ,
		geography TEXT NOT NULL,
		about TEXT);
		CREATE INDEX IF NOT EXISTS about ON employee(about);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtEmp.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func CreateTableUser(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.User"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		name TEXT NOT NULL,
		phoneNumber TEXT NOT NULL UNIQUE
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = res.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func CreateVacancyTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmtVacancy, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vacancy(
		id INTEGER PRIMARY KEY,
		employee_id INTEGER,
		name TEXT NOT NULL,
		price INTEGER,
		location TEXT NOT NULL,
		experience TEXT);
		CREATE INDEX IF NOT EXISTS price ON vacancy(price);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtVacancy.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllVacancy() ([]ResponseVac, error) {

	const op = "storage.sqlite.Get.AllVacancy"
	_, err := s.db.Prepare("SELECT * FROM vacancy")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		log.Fatal(1)
	}
	result := []ResponseVac{}
	row, err := s.db.Query("SELECT * FROM vacancy")
	if err != nil {
		fmt.Println(err, "Error")
		return nil, nil
	}
	for row.Next() {
		r := ResponseVac{}
		err := row.Scan(&r.ID, &r.Emp_ID, &r.Vac_Name, &r.Price, &r.Location, &r.Experience)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// r.Status = resp.OK().Status
		result = append(result, r)
	}
	fmt.Println()
	return result, nil
}

func (s *Storage) GetAllVacsForEmployee(emp_id int) ([]ResponseVac, error) {
	const op = "storage.sqlite.Get.AllVacancy"
	_, err := s.db.Prepare("SELECT * FROM vacancy WHERE employee_id = ?")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, fmt.Errorf("ERROR IN CREATING REQUEST OT DB")
	}
	result := []ResponseVac{}
	row, err := s.db.Query("SELECT * FROM vacancy WHERE employee_id = ?", emp_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: ошибка в бд (xdd)", op)

		} else {
			return result, fmt.Errorf("%s: какая-то ошибка в получении работодателя по его id. Если вы это видите, то напишите разрабу и скажите что он даун xdd", op)
		}
	}
	for row.Next() {
		r := ResponseVac{}
		err := row.Scan(&r.ID, &r.Emp_ID, &r.Vac_Name, &r.Price, &r.Location, &r.Experience)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// r.Status = resp.OK().Status
		result = append(result, r)
	}
	return result, nil
}

func (s *Storage) GetEmployee(ID int) (RequestEmployee, error) {
	const op = "storage.sqlite.Get.EmployeeByIDs"
	var result RequestEmployee
	stmtVacancy, err := s.db.Prepare("SELECT * FROM employee WHERE id = ?")
	if err != nil {
		return result, fmt.Errorf("%s: ошибка в создании запроса к бд", op)
	}
	_ = stmtVacancy

	err = s.db.QueryRow("SELECT * FROM employee WHERE id = ?", ID).Scan(&result.ID, &result.Limit, &result.NameOrganization, &result.PhoneNumber, &result.Email, &result.Geography, &result.About)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: ошибка в бд (xdd)", op)

		} else {
			return result, fmt.Errorf("%s: какая-то ошибка в получении работодателя по его id. Если вы это видите, то напишите разрабу и скажите что он даун xdd", op)
		}
	}

	return result, nil
}

func (s *Storage) VacancyByID(ID int) (ResponseVac, error) {
	const op = "storage.sqlite.Get.VacancyByIDs"
	var result ResponseVac
	_, err := s.db.Prepare("SELECT * FROM vacancy WHERE id = ?")

	if err != nil {
		return result, fmt.Errorf("%s: ошибка в создании запроса к бд", op)
	}

	err = s.db.QueryRow("SELECT * FROM vacancy WHERE id = ?", ID).Scan(&result.ID, &result.Emp_ID, &result.Vac_Name, &result.Price, &result.Location, &result.Experience)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: ошибка в бд (xdd)", op)

		} else {
			return result, fmt.Errorf("%s: какая-то ошибка в получении вакансии по её id. Если вы это видите, то напишите разрабу и скажите что он даун xdd", op)
		}
	}

	return result, nil
}

func (s *Storage) AddEmployee(limitIsOver int, nameOrganization string, phoneNumber string, email string, geography string, about string) (int64, error) {
	const op = "storage.sqlite.Add.Emp"
	stmt, err := s.db.Prepare("INSERT INTO employee(limitVac ,nameOrganization,phoneNumber,email,geography,about) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(limitIsOver, nameOrganization, phoneNumber, email, geography, about)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: Произошла ошибка в добавлении данных в бд. Вероятно, такие данные уже пользуются", op)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) AddVacancy(employee_id int, name string, price int, location string, experience string) (int64, int64, error) {
	const op = "storage.sqlite.Add.Vacancy"

	stmtVacancy, err := s.db.Prepare("INSERT INTO vacancy(employee_id,name ,price,location,experience) VALUES (?,?,?,?,?)")

	if err != nil {
		return -1, -1, fmt.Errorf("%s: %w\n\t error in try to prepare sql request", op, err)
	}
	limit := s.GetLimit(employee_id)
	if limit == -1 {
		return -1, -1, fmt.Errorf("%s: The employer has reached the limit", op)
	}

	resultd, err := stmtVacancy.Exec(employee_id, name, price, location, experience)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, -1, fmt.Errorf("%s: error in try to get sqlite", op)
		}
		return -1, -1, fmt.Errorf("%s: %w", op, err)
	}

	vac_id, err := resultd.LastInsertId()
	if err != nil {
		return -1, -1, fmt.Errorf("%s: %w", op, err)
	}

	return vac_id, int64(limit), nil
}

func (s *Storage) GetLimit(ID int) int {

	_, err := s.db.Prepare("SELECT limitVac FROM employee WHERE id = ?")
	if err != nil {
		return -1
	}
	var count int
	row := s.db.QueryRow("SELECT limitVac FROM employee WHERE id = $1", ID)
	err = row.Scan(&count)
	if err != nil {
		return -1
	}
	if count >= 10 {
		return -1
	}
	update := count + 1
	_, err = s.db.Prepare("UPDATE employee SET limitVac = ? WHERE id = ?")
	if err != nil {
		return -1
	}
	_, err = s.db.Exec("UPDATE employee SET limitVac = $1 WHERE id = $2", update, ID)
	if err != nil {
		return -1
	}
	return update
}

func CreateToken(email string) (string, error) {
	type Header struct {
		Alg string `json:"alg"` // Алгоритм подписи
		Typ string `json:"typ"` // Тип токена
	}

	type Payload struct {
		Iss string `json:"iss"`
		Sub string `json:"sub"` // Subject (обычно идентификатор пользователя)
		Iat int64  `json:"iat"` // Issued at - время в которое был выдан токен
		Exp int64  `json:"exp"` // Время истечения токена (в Unix timestamp)
	}
	var secretKEY string = "ISP-7-21-borodinna"

	var header Header
	header.Alg = "HS256"
	header.Typ = "JWT"

	var payload Payload
	payload.Iss = "Nick005-aka-monkeyZV"
	payload.Sub = email
	payload.Iat = time.Now().Unix()
	payload.Exp = time.Now().Add(time.Second * 60).Unix()

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "error", fmt.Errorf("error in converting HEADER to JSON")
	}
	headerBASE64 := base64.RawURLEncoding.Strict().EncodeToString(headerJSON)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "error", fmt.Errorf("error in converting PAYLOAD to JSON")
	}
	payloadBASE64 := base64.RawURLEncoding.Strict().EncodeToString(payloadJSON)

	// создаем подпись для JWTшки
	signaturePayAndHeader := fmt.Sprintf("%s.%s", headerBASE64, payloadBASE64)

	h := hmac.New(sha256.New, []byte(secretKEY))
	h.Write([]byte(signaturePayAndHeader))
	var signature string = base64.RawStdEncoding.EncodeToString(h.Sum(nil))

	var tokenJWT string = fmt.Sprintf("%s.%s.%s", headerBASE64, payloadBASE64, signature)

	return tokenJWT, nil
}

func (s *Storage) CreateAccessToken(email string, uid int) (int, string, error) {
	const op = "sqlite.CreateAccessToken.user"
	token, err := CreateToken(email)
	if err != nil {
		return -1, "error", err
	}
	stmtUser, err := s.db.Prepare("INSERT INTO token(user_id, active_token, is_active) VALUES (?,?,?)")
	if err != nil {
		return -1, "error", fmt.Errorf("%s: %w", op, err)
	}
	indexd, err := stmtUser.Exec(uid, token, 1)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, "error", fmt.Errorf("%s: %w", op, err)
		}
		return -1, "error", fmt.Errorf("%s: %w", op, err)
	}
	user_id, err := indexd.LastInsertId()
	if err != nil {
		return -1, "error", fmt.Errorf("%s: %w", op, err)
	}
	return int(user_id), token, nil
}

func (s *Storage) AddUser(email string, password string, name string, phoneNumber string) (int, error) {
	const op = "storage.sqlite.Add.User"
	stmtUser, err := s.db.Prepare("INSERT INTO user(email, password, name , phoneNumber) VALUES (?,?,?,?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	indexd, err := stmtUser.Exec(email, password, name, phoneNumber)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %s", op, "такой пользователь уже существует!")
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	uid, err := indexd.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	return int(uid), nil
}

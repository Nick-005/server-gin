package sqlite

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

type RequestNewToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Storage struct {
	db *sql.DB
}

type RequestEmployee struct {
	ID               int    `json:"ID"`
	NameOrganization string `json:"nameOrg"`
	PhoneNumber      string `json:"phoneNumber"`
	Email            string `json:"email"`
	INN              string `json:"inn"`
	Status           string `json:"status"`
}

type ResponseVac struct {
	ID          int    `json:"ID"`
	Emp_ID      int    `json:"emp_id"`
	Vac_Name    string `json:"vac_name"`
	Price       int    `json:"price"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`
	Experience  string `json:"exp"`
	About       string `json:"about"`
	Is_visible  bool   `json:"is_visible"`
}

type ResponseSearchVac struct {
	ID         int    `json:"ID"`
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
}

// done
func CreateVacancyTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmtVacancy, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vacancy(
			id INTEGER PRIMARY KEY,
			emp_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			price REAL,
			email TEXT,
			phoneNumber TEXT,
			location TEXT,
			experience_id INTEGER NOT NULL,
			aboutWork TEXT,
			is_visible BOOLEAN DEFAULT TRUE,
			FOREIGN KEY (emp_id) REFERENCES employer(id) ON DELETE CASCADE,
			FOREIGN KEY (experience_id) REFERENCES experience(id)
	);
	
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

// done
func CreateResponeVacTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Response"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtResp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS response(
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL, 
		vacancy_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
		status_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE,
		FOREIGN KEY (vacancy_id) REFERENCES vacancy(id) ON DELETE CASCADE,
		FOREIGN KEY (status_id) REFERENCES status(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmtResp.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// done
func CreateEmployeeTable(storagPath string) (*Storage, error) {
	const op = "storage.sqlite.Emp"
	db, err := sql.Open("sqlite3", storagPath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	stmtEmp, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS employer(
		id INTEGER PRIMARY KEY,
		nameOrganization TEXT NOT NULL,
		phoneNumber TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE ,
		INN TEXT NOT NULL UNIQUE,
		status_id INTEGER NOT NULL,
		FOREIGN KEY (status_id) REFERENCES status(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS about ON employer(status_id);
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

// done
func CreateTableUser(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.User"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS user(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		phoneNumber TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		resume_id INTEGER UNIQUE,
		status_id INTEGER NOT NULL,
		FOREIGN KEY (status_id) REFERENCES status(id) ON DELETE CASCADE,
		FOREIGN KEY (resume_id) REFERENCES resume(id)
	);

	CREATE INDEX IF NOT EXISTS idx_user_status_id ON user(status_id);
	CREATE INDEX IF NOT EXISTS idx_user_resume_id ON user(resume_id);
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

// done
func CreateStatusTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.Status"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	// Создаем таблицу
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS status(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Добавляем начальные значения
	_, err = db.Exec(`
	INSERT INTO status (id, name) VALUES (1, 'Заблокирован')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO status (id, name) VALUES (2, 'Активен')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO status (id, name) VALUES (3, 'Требует активации')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO status (id, name) VALUES (4, 'Приглашение')
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO status (id, name) VALUES (5, 'Отказано')
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO status (id, name) VALUES (6, 'Ожидание')
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO status (id, name) VALUES (7, 'Активный поиск')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO status (id, name) VALUES (8, 'Рассматриваю предложения')
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO status (id, name) VALUES (9, 'Предложили работу, пока думаю')
	ON CONFLICT (id) DO NOTHING;

	INSERT INTO status (id, name) VALUES (10, 'Не ищу работу')
	ON CONFLICT (id) DO NOTHING;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// done
func CreateResumeTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.Resume"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS resume(
		id INTEGER PRIMARY KEY,
		experience_id INTEGER NOT NULL,
		description TEXT NOT NULL,
		FOREIGN KEY (experience_id) REFERENCES experience(id)
	);
	CREATE INDEX IF NOT EXISTS idx_resume_id ON user(experience_id);
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

// done
func CreateExperienceTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.Experience"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS experience(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Добавляем начальные значения
	_, err = db.Exec(`
	INSERT INTO experience (id, name) VALUES (1, 'Без опыта')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO experience (id, name) VALUES (2, '1-3 года')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO experience (id, name) VALUES (3, '3-5 лет')
	ON CONFLICT (id) DO NOTHING;
	
	INSERT INTO experience (id, name) VALUES (4, 'Более 5 лет')
	ON CONFLICT (id) DO NOTHING;
	
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetAllVacancy() ([]ResponseVac, error) {
	result := []ResponseVac{}
	const op = "storage.sqlite.Get.AllVacancy"
	_, err := s.db.Prepare("SELECT * FROM vacancy")
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST TO DB!", op)
		return result, fmt.Errorf("error in creating request to DB")
	}

	row, err := s.db.Query("SELECT * FROM vacancy")
	if err != nil {
		return nil, fmt.Errorf("error in exec sql script")
	}
	for row.Next() {
		r := ResponseVac{}
		err := row.Scan(&r.ID, &r.Emp_ID, &r.Vac_Name, &r.Price, &r.Email, &r.PhoneNumber, &r.Location, &r.Experience, &r.About, &r.Is_visible)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, r)
	}
	fmt.Println()
	return result, nil
}

func (s *Storage) GetAllVacsForEmployee(emp_id int) ([]ResponseVacancyByIDs, error) {
	const op = "storage.sqlite.Get.AllVacancy"
	_, err := s.db.Prepare(`select vacancy.id, vacancy.emp_id, employer.nameOrganization,vacancy.name,vacancy.price, vacancy.email, vacancy.phoneNumber, vacancy.location, experience.name, vacancy.aboutWork, vacancy.is_visible 
							from vacancy 
							INNER JOIN employer on vacancy.emp_id = employer.id 
							INNER JOIN experience on vacancy.experience_id = experience.id
							where vacancy.emp_id = ?`)
	if err != nil {
		fmt.Println("ERROR IN CREATING REQUEST OT DB!", op)
		return nil, err
	}
	result := []ResponseVacancyByIDs{}
	row, err := s.db.Query(`select vacancy.id, vacancy.emp_id, employer.nameOrganization,vacancy.name,vacancy.price, vacancy.email, vacancy.phoneNumber, vacancy.location, experience.name, vacancy.aboutWork, vacancy.is_visible 
							from vacancy 
							INNER JOIN employer on vacancy.emp_id = employer.id 
							INNER JOIN experience on vacancy.experience_id = experience.id
							where vacancy.emp_id = ?`, emp_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: ошибка в бд (xdd)", op)

		} else {
			return result, fmt.Errorf("%s: какая-то ошибка в получении работодателя по его id. Если вы это видите, то напишите разрабу и скажите что он даун xdd", op)
		}
	}
	for row.Next() {
		r := ResponseVacancyByIDs{}
		err := row.Scan(&r.ID, &r.Emp_ID, &r.Employer_Name, &r.Vac_Name, &r.Price, &r.Email, &r.PhoneNumber, &r.Location, &r.Experience, &r.About, &r.Is_visible)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// r.Status = resp.OK().Status
		result = append(result, r)
	}
	return result, nil
}

type AllResponse struct {
	Name             string `json:"name"`
	Price            int    `json:"price"`
	NameOrganization string `json:"nameOrg"`
	Location         string `json:"location"`
	Experience       string `json:"exp"`
	Status           string `json:"status"`
}

func (s *Storage) GetAllResponse(UID int) ([]AllResponse, error) {
	var result []AllResponse
	rows, err := s.db.Query(`SELECT vacancy.name, vacancy.price, employer.nameOrganization, vacancy.location, experience.name, status.name from response
							INNER JOIN user on response.user_id = user.id
							INNER JOIN vacancy on response.vacancy_id = vacancy.id
							INNER JOIN employer on vacancy.emp_id = employer.id
							INNER JOIN experience on vacancy.experience_id = experience.id
							INNER JOIN status on response.status_id = status.id
							where response.user_id = $1`, UID)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var response AllResponse
		err := rows.Scan(&response.Name, &response.Price, &response.NameOrganization, &response.Location, &response.Experience, &response.Status)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, response)
	}
	return result, nil
}

func (s *Storage) GetEmployee(ID int) (RequestEmployee, error) {
	const op = "storage.sqlite.Get.EmployeeByIDs"
	var result RequestEmployee
	stmtVacancy, err := s.db.Prepare(`SELECT employer.id, nameOrganization, phoneNumber, email, INN, status.name from employer
									INNER JOIN status on employer.status_id = status.id 
									where employer.id = $1`)
	if err != nil {
		return result, fmt.Errorf("%s: ошибка в создании запроса к бд", op)
	}
	_ = stmtVacancy

	err = s.db.QueryRow(`SELECT employer.id, nameOrganization, phoneNumber, email, INN, status.name from employer
						INNER JOIN status on employer.status_id = status.id
						where employer.id = $1`, ID).Scan(&result.ID, &result.NameOrganization, &result.PhoneNumber, &result.Email, &result.INN, &result.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: %w", op, err)

		} else {
			return result, fmt.Errorf("%s: %w", op, err)
		}
	}

	return result, nil
}

type ResponseVacancyByIDs struct {
	ID            int    `json:"ID"`
	Emp_ID        int    `json:"emp_id"`
	Employer_Name string `json:"emp_name"`
	Vac_Name      string `json:"vac_name"`
	Price         int    `json:"price"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phoneNumber"`
	Location      string `json:"location"`
	Experience    string `json:"exp"`
	About         string `json:"about"`
	Is_visible    bool   `json:"is_visible"`
}

func (s *Storage) VacancyByID(ID int) (ResponseVacancyByIDs, error) {
	const op = "storage.sqlite.Get.VacancyByIDs"
	var result ResponseVacancyByIDs
	_, err := s.db.Prepare(`select vacancy.id, emp_id, employer.nameOrganization, vacancy.name, vacancy.price, vacancy.email, vacancy.phoneNumber, vacancy.location, experience.name, vacancy.aboutWork, vacancy.is_visible
							from vacancy
							INNER JOIN employer on vacancy.emp_id = employer.id
							INNER JOIN experience on vacancy.experience_id = experience.id
							where vacancy.id = $1`)

	if err != nil {
		return result, fmt.Errorf("%s: ошибка в создании запроса к бд", op)
	}

	err = s.db.QueryRow(`select vacancy.id, emp_id, employer.nameOrganization, vacancy.name,
								vacancy.price, vacancy.email, vacancy.phoneNumber, vacancy.location,
								experience.name, vacancy.aboutWork, vacancy.is_visible
								from vacancy
								INNER JOIN employer on vacancy.emp_id = employer.id
								INNER JOIN experience on vacancy.experience_id = experience.id
								where vacancy.id = $1`, ID).Scan(
		&result.ID, &result.Emp_ID, &result.Employer_Name, &result.Vac_Name,
		&result.Price, &result.Email, &result.PhoneNumber, &result.Location,
		&result.Experience, &result.About, &result.Is_visible)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: %w", op, err)

		} else {
			return result, fmt.Errorf("%s: %w", op, err)
		}
	}

	return result, nil
}

type VacancyTake struct {
	ID               int    `json:"ID"`
	Emp_ID           int    `json:"emp_id"`
	Vac_Name         string `json:"vac_name"`
	NameOrganization string `json:"nameOrg"`
	Price            int    `json:"price"`
	Email            string `json:"email"`
	PhoneNumber      string `json:"phoneNumber"`
	Location         string `json:"location"`
	Experience       string `json:"exp"`
	About            string `json:"about"`
	Is_visible       bool   `json:"is_visible"`
}

func (s *Storage) VacancyByLimit(limit, last_id int) ([]VacancyTake, error) {
	const op = "storage.sqlite.Get.VacancyByIDs"
	var result []VacancyTake
	rows, err := s.db.Query(`SELECT vacancy.id, emp_id, vacancy.name , employer.nameOrganization, price, employer.email, employer.phoneNumber,  location, experience.name, is_visible, aboutWork FROM vacancy 
							INNER JOIN employer on vacancy.emp_id = employer.id
							INNER JOIN experience on vacancy.experience_id = experience.id
							where vacancy.id > ? order by vacancy.id limit ?
							`, last_id, limit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, fmt.Errorf("%s: ошибка в бд (xdd). %w", op, err)

		} else {
			return result, fmt.Errorf("%s: какая-то ошибка в получении вакансии по её id. Если вы это видите, то напишите разрабу и скажите что он даун xdd. %w", op, err)
		}
	}
	defer rows.Close()

	for rows.Next() {
		vac := VacancyTake{}
		err := rows.Scan(&vac.ID, &vac.Emp_ID, &vac.Vac_Name,
			&vac.NameOrganization,
			&vac.Price, &vac.Email, &vac.PhoneNumber, &vac.Location,
			&vac.Experience, &vac.Is_visible, &vac.About)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, vac)
	}

	return result, nil
}

func (s *Storage) AddEmployee(nameOrganization string, phoneNumber string, email string, inn string, statusID int) (int64, error) {
	const op = "storage.sqlite.Add.Emp"
	stmt, err := s.db.Prepare("INSERT INTO employer(nameOrganization,phoneNumber,email,INN,status_id) VALUES (?,?,?,?,?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(nameOrganization, phoneNumber, email, inn, statusID)
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

func (s *Storage) AddVacancy(employee_id int, name string, price int, email string, phoneNumber string,
	location string, experience int, about string, visible bool) (int64, error) {

	const op = "storage.sqlite.Add.Vacancy"

	stmtVacancy, err := s.db.Prepare("INSERT INTO vacancy(emp_id,name ,price,email,phoneNumber,location, experience_id, aboutWork, is_visible)" +
		"VALUES (?,?,?,?,?,?,?,?,?)")

	if err != nil {
		return -1, fmt.Errorf("%s: %w\n\t error in try to prepare sql request", op, err)
	}

	resultd, err := stmtVacancy.Exec(employee_id, name, price, email, phoneNumber, location, experience, about, visible)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: error in try to get sqlite", op)
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	vac_id, err := resultd.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return vac_id, nil
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
	payload.Exp = time.Now().Add(time.Minute * 15).Unix()

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

func (s *Storage) CreateAccessToken(email string) (string, error) {
	const op = "sqlite.CreateAccessToken.User"
	token, err := CreateToken(email)
	if err != nil {
		return "Error", fmt.Errorf("%s: %w", op, err)
	}
	token = strings.ReplaceAll(token, "+", "-")
	token = strings.ReplaceAll(token, "/", "_")
	return token, nil
}

func (s *Storage) CheckVacancyExist(vacancyID int) error {
	row := s.db.QueryRow("SELECT id from vacancy where vacancy.id = $1", vacancyID)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("your vacancyID doesn't exist! Pls check ur request and try again")
	}
	return nil
}

func (s *Storage) CheckUserExist(UID int) error {
	row := s.db.QueryRow("SELECT id from user where id = $1", UID)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("your UID doesn't exist! Pls check ur request and try again")
	}
	return nil
}

func (s *Storage) MakeResponse(UID, vacancyID int) (int64, error) {
	row := s.db.QueryRow("Select id from response where user_id = $1 and vacancy_id = $2", UID, vacancyID)
	var userID int
	err := row.Scan(&userID)
	if err != sql.ErrNoRows {
		return -1, fmt.Errorf("you have already applied for this position. Please check ur request")
	}
	result, err := s.db.Exec("Insert into response (user_id, vacancy_id, status_id) values ($1, $2, $3)", UID, vacancyID, 6)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return id, nil
}

/*
name TEXT NOT NULL,
phoneNumber TEXT NOT NULL UNIQUE,
email TEXT NOT NULL UNIQUE,
password TEXT NOT NULL,
resume_id INTEGER,
status_id INTEGER,
*/
func (s *Storage) AddUser(email string, password string, name string, phoneNumber string) (int, error) {
	const op = "storage.sqlite.Add.User"
	stmtUser, err := s.db.Prepare("INSERT INTO user(name, phoneNumber,email ,password , status_id) VALUES (?,?,?,?,?)")
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer stmtUser.Close()
	indexd, err := stmtUser.Exec(name, phoneNumber, email, password, 8)
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

func (s *Storage) CheckPasswordNEmail(email, password string) (RequestNewToken, error) {
	const op = "sqlite.CheckPasswordNEmail.New.Token"

	row := s.db.QueryRow("select email, password from user where email = $1 and password = $2", email, password)
	var body RequestNewToken
	err := row.Scan(&body.Email, &body.Password)
	if err != nil {
		return body, fmt.Errorf("%s: %w", op, err)
	}
	return body, nil
}

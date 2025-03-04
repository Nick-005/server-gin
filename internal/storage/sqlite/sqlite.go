package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

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

// done
func CreateVacancyTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	stmtVacancy, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vacancy(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			emp_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			price REAL,
			mail TEXT,
			phoneNumber TEXT,
			location TEXT,
			experience_id INTEGER,
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
		id INTEGER PRIMARY KEY AUTOINCREMENT,
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
		nameOrganization TEXT NOT NULL UNIQUE,
		phoneNumber TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE ,
		INN TEXT NOT NULL,
		about TEXT,
		limitVac INTEGER);
		CREATE INDEX IF NOT EXISTS about ON employer(limitVac);
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
		resume_id INTEGER,
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
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS status(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);
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
func CreateResumeTable(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New.Resume"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS resume(
		id INTEGER PRIMARY KEY,
		experience_id INTEGER,
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
	res, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS experience(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
	);
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

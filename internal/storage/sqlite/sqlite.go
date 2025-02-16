package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type ResponseVac struct {
	ID         int    `json:"ID"`
	Emp_ID     int    `json:"emp_id"`
	Vac_Name   string `json:"vac_name"`
	Price      int    `json:"price"`
	Location   string `json:"location"`
	Experience string `json:"exp"`
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

package sqlite

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type RequestNewToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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

type GetStatus struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Crated_At time.Time `db:"created_at"`
}

func GetAllStatus(storage *sqlx.DB) ([]GetStatus, error) {
	var result []GetStatus
	const op = "storage.sqlite.Get.AllVacancy"
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

package structs

import "time"

type SuccessEmployer struct {
	ID               int       `db:"id"`
	NameOrganization string    `db:"name_organization"`
	PhoneNumber      string    `db:"phone_number"`
	Email            string    `db:"email"`
	INN              string    `db:"inn"`
	Status           string    `db:"status_id"`
	Created_at       time.Time `db:"created_at"`
	Updated_at       time.Time `db:"updated_at"`
}

type Vacancy_Body struct {
	Emp_ID      int    `json:"emp_id"`
	Vac_Name    string `json:"vac_name"`
	Price       int    `json:"price"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Location    string `json:"location"`
	Experience  int    `json:"exp"`
	About       string `json:"about"`
	Is_visible  bool   `json:"is_visible"`
}

type RequestVac struct {
	Limit   int `json:"limit"`
	Last_id int `json:"last_id"`
}

type RequestAdd struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// const secretKEY = "ISP-7-21-borodinna"

type Status struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Vacancies struct {
	ID           int       `db:"id"`
	Emp_id       int       `db:"emp_id"`
	Name         string    `db:"name"`
	Price        string    `db:"price"`
	Email        string    `db:"email"`
	PhoneNumber  string    `db:"phone_number"`
	Location     string    `db:"location"`
	ExperienceId int       `db:"experience_id"`
	AboutWork    string    `db:"about_work"`
	IsVisible    bool      `db:"is_visible"`
	Created_at   time.Time `db:"created_at"`
	Updated_at   time.Time `db:"updated_at"`
}
type RequestEmployee struct {
	NameOrganization string `json:"name_organization"`
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

type RequestNewToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

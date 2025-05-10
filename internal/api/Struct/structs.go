package structs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SuccessEmployer struct {
	ID               int    `db:"id"`
	NameOrganization string `db:"name_organization"`
	PhoneNumber      string `db:"phone_number"`
	Email            string `db:"email"`
	INN              string `db:"inn"`
	Password         string `db:"password"`
	Status           struct {
		ID        int       `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
	} `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Vacancy_Body struct {
	Emp_ID      int    `json:"emp_id"`
	Vac_Name    string `json:"vac_name"`
	Price       int    `json:"price"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Location    string `json:"location"`
	Experience  int    `json:"exp"`
	About       string `json:"about"`
	Is_visible  bool   `json:"is_visible"`
}

type RequestCandidate struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Status_id   int    `json:"status_id"`
}

type RequestVac struct {
	Limit   int `json:"limit"`
	Last_id int `json:"last_id"`
}

type RequestResume struct {
	Experience  int    `json:"experience_id"`
	Description string `json:"description"`
}

type RequestResumeUpdate struct {
	Experience  int    `json:"experience_id"`
	Description string `json:"description"`
	Resume_id   int    `json:"resume_id"`
}

type AllUserResponseOK struct {
	Status  string
	Otkliks string
}

type Ok struct {
	Status string
}

type SimpleError struct {
	Status string
	Error  string
}

type InfoError struct {
	SimpleError
	Info string
}

type Authorization struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestAdd struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type Status struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type ResponsePatch struct {
	Response_id int `json:"response_id"`
	Status_id   int `json:"status_id"`
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
	Password         string `json:"password"`
	Status_id        int    `json:"status_id"`
}

type DBResponse struct {
	ID            int       `db:"id"`
	Candidates_id int       `db:"candidates_id"`
	Vacancy_id    int       `db:"vacancy_id"`
	Created_at    time.Time `db:"created_at"`
	Status_id     int       `db:"status_id"`
}

type VacanciesToResponse struct {
	ID            int       `db:"id"`
	Employee_name string    `db:"employee_name"`
	Name          string    `db:"name"`
	Price         string    `db:"price"`
	Email         string    `db:"email"`
	PhoneNumber   string    `db:"phone_number"`
	Location      string    `db:"location"`
	Experience    GetStatus `db:"experience"`
	AboutWork     string    `db:"about_work"`
	IsVisible     bool      `db:"is_visible"`
	Created_at    time.Time `db:"created_at"`
	Updated_at    time.Time `db:"updated_at"`
}

type SuccessResponse struct {
	Vacancy   VacancyData `db:"vacancy"`
	Responses []struct {
		ID         int           `db:"id"`
		Candidate  InfoCandidate `db:"candidate"`
		Created_at time.Time     `db:"created_at"`
		Status     GetStatus     `db:"status"`
	} `db:"responses"`
}

type ResponseByVac struct {
	ID      int                 `db:"id"`
	Vacancy VacanciesToResponse `db:"vacancy"`
	Status  GetStatus           `db:"status"`
}

type RequestResponse struct {
	Vacancy_id int `json:"vac_id"`
	Status_id  int `json:"status_id"`
}

type SuccessVacancy struct {
	Employee  SuccessEmployer `db:"employer"`
	Vacancies []VacancyData   `db:"vacancies"`
}

type VacancyData_Limit struct {
	ID          int             `db:"id"`
	Employee    SuccessEmployer `db:"employer"`
	Name        string          `db:"name"`
	Price       string          `db:"price"`
	Email       string          `db:"email"`
	PhoneNumber string          `db:"phone_number"`
	Location    string          `db:"location"`
	Experience  GetStatus       `db:"experience"`
	AboutWork   string          `db:"about_work"`
	IsVisible   bool            `db:"is_visible"`
	Created_at  time.Time       `db:"created_at"`
	Updated_at  time.Time       `db:"updated_at"`
}

type VacancyData struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Price       string    `db:"price"`
	Email       string    `db:"email"`
	PhoneNumber string    `db:"phone_number"`
	Location    string    `db:"location"`
	Experience  GetStatus `db:"experience"`
	AboutWork   string    `db:"about_work"`
	IsVisible   bool      `db:"is_visible"`
	Created_at  time.Time `db:"created_at"`
	Updated_at  time.Time `db:"updated_at"`
}

type VacancyPut struct {
	ID            int    `json:"id"`
	Vac_Name      string `json:"vac_name"`
	Price         string `json:"price"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phoneNumber"`
	Location      string `json:"location"`
	Experience_Id int    `json:"exp_id"`
	About         string `json:"about"`
	Is_visible    bool   `json:"is_visible"`
}

type ResponseVac struct {
	Vac_Name      string `json:"vac_name"`
	Price         string `json:"price"`
	Email         string `json:"email"`
	PhoneNumber   string `json:"phoneNumber"`
	Location      string `json:"location"`
	Experience_Id int    `json:"exp_id"`
	About         string `json:"about"`
	Is_visible    bool   `json:"is_visible"`
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
type ResumeResult_slice struct {
	Id          int       `db:"id"`
	Experience  GetStatus `db:"experience"`
	Description string    `db:"description"`
	Created_at  time.Time `db:"created_at"`
	Updated_at  time.Time `db:"updated_at"`
}
type ResumeResult struct {
	Resumes   []ResumeResult_slice `db:"resume"`
	Candidate InfoCandidate        `db:"candidate"`
}
type InfoCandidate struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	PhoneNumber string `db:"phone_number"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Status      struct {
		ID        int       `db:"id"`
		Name      string    `db:"name"`
		Crated_At time.Time `db:"created_at"`
	} `db:"status"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}

type SuccessResume struct {
	Id            int       `db:"id"`
	Candidate_id  int       `db:"candidate_id"`
	Experience_id int       `db:"experience_id"`
	Description   string    `db:"description"`
	Created_at    time.Time `db:"created_at"`
	Updated_at    time.Time `db:"updated_at"`
}

type RequestNewToken struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	ID    int    `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

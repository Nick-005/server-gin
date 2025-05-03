package structs

import "time"

type SuccessEmployer struct {
	ID               int    `db:"id"`
	NameOrganization string `db:"name_organization"`
	PhoneNumber      string `db:"phone_number"`
	Email            string `db:"email"`
	INN              string `db:"inn"`
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
	UserStatus  string `json:"u_status"`
}

type RequestVac struct {
	Limit   int `json:"limit"`
	Last_id int `json:"last_id"`
}

type RequestResume struct {
	UserEmail   string `json:"user_email"`
	Experience  string `json:"exp_name"`
	Description string `json:"description"`
}

type RequestAdd struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
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

type DBResponse struct {
	ID            int       `db:"id"`
	Candidates_id int       `db:"candidates_id"`
	Vacancy_id    int       `db:"vacancy_id"`
	Created_at    time.Time `db:"created_at"`
	Status_id     int       `db:"status_id"`
}

type RequestResponse struct {
	Candidate_email string `json:"can_email"`
	Vacancy_id      int    `json:"vac_id"`
	Status_name     string `json:"status"`
}

type ResponseVac struct {
	Emp_Email   string `json:"emp_email"`
	Vac_Name    string `json:"vac_name"`
	Price       string `json:"price"`
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

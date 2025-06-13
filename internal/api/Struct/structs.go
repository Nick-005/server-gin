package structs

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type SuccessAllEmployers struct {
	Status        string            `json:"Status"`
	EmployersInfo []SuccessEmployer `json:"EmployersInfo"`
}

type SuccessEmployer struct {
	ID               int    `db:"id" json:"ID"`
	NameOrganization string `db:"name_organization" json:"NameOrganization"`
	PhoneNumber      string `db:"phone_number"  json:"PhoneNumber"`
	Email            string `db:"email" json:"Email"`
	INN              string `db:"inn" json:"INN"`
	Password         string `db:"password" json:"Password"`
	Status           struct {
		ID        int       `db:"id" json:"ID"`
		Name      string    `db:"name" json:"Name"`
		CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	} `db:"status" json:"Status"`
	CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt time.Time `db:"updated_at" json:"UpdatedAt"`
}

type ResponseOnVacancy struct {
	IsResponsed bool `json:"IsResponsed"`
	Status      struct {
		ID        int       `db:"id" json:"ID"`
		Name      string    `db:"name" json:"Name"`
		CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	} `db:"status" json:"Status"`
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
	IsVisible   bool   `json:"is_visible"`
}

type RequestCandidate struct {
	Name        string `json:"Name"`
	PhoneNumber string `json:"PhoneNumber"`
	Email       string `json:"Email"`
	Password    string `json:"Password"`
	Status_id   int    `json:"StatusId"`
}

type RequestVac struct {
	Limit   int `json:"Limit"`
	Last_id int `json:"LastID"`
}

type RequestResume struct {
	Experience  int    `json:"ExperienceID"`
	Description string `json:"Description"`
}

type RequestResumeUpdate struct {
	Experience  int    `json:"ExperienceID"`
	Description string `json:"description"`
	Resume_id   int    `json:"ResumeID"`
}

type AllUserResponseOK struct {
	Status  string `json:"Status"`
	Otkliks string ` json:"Responses"`
}

type Ok struct {
	Status string `json:"Status"`
}

type SimpleError struct {
	Status string `json:"Status"`
	Error  string `json:"Error"`
}

type InfoError struct {
	SimpleError
	Info string `json:"Info"`
}

type Authorization struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type RequestAdd struct {
	Name        string `json:"Name"`
	PhoneNumber string `json:"PhoneNumber"`
	Email       string `json:"Email"`
	Password    string `json:"Password"`
}

type Status struct {
	ID   int    `db:"id" json:"ID"`
	Name string `db:"name" json:"Name"`
}

type ResponsePatch struct {
	Response_id int `json:"ResponseID"`
	Status_id   int `json:"StatusID"`
}

type EmployerInfoForPatch struct {
	Employer_id int `json:"EmployerID"`
	Status_id   int `json:"StatusID"`
}

type Vacancies struct {
	ID           int       `db:"id" json:"ID"`
	Emp_ID       int       `db:"emp_id" json:"EmployerID"`
	Name         string    `db:"name" json:"Name"`
	Price        int       `db:"price" json:"Price"`
	Email        string    `db:"email" json:"Email"`
	PhoneNumber  string    `db:"phone_number" json:"PhoneNumber"`
	Location     string    `db:"location" json:"Location"`
	ExperienceId int       `db:"experience_id" json:"ExperienceID"`
	AboutWork    string    `db:"about_work" json:"AboutWork"`
	IsVisible    bool      `db:"is_visible" json:"IsVisible"`
	CreatedAt    time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"UpdatedAt"`
}
type RequestEmployee struct {
	NameOrganization string `json:"NameOrganization"`
	PhoneNumber      string `json:"PhoneNumber"`
	Email            string `json:"Email"`
	INN              string `json:"INN"`
	Password         string `json:"Password"`
	Status_id        int    `json:"StatusID"`
}

type RequestEmployer struct {
	NameOrganization string `json:"NameOrganization"`
	PhoneNumber      string `json:"PhoneNumber"`
	Email            string `json:"Email"`
	Password         string `json:"Password"`
	Status_id        int    `json:"StatusID"`
}

// неиспользуется
type DBResponse struct {
	ID            int       `db:"id"`
	Candidates_id int       `db:"candidates_id"`
	Vacancy_id    int       `db:"vacancy_id"`
	CreatedAt     time.Time `db:"created_at"`
	Status_id     int       `db:"status_id"`
}

type VacanciesToResponse struct {
	ID            int       `db:"id" json:"ID"`
	Employer_name string    `db:"employee_name" json:"EmployerName"`
	Name          string    `db:"name" json:"Name"`
	Price         int       `db:"price" json:"Price"`
	Email         string    `db:"email" json:"Email"`
	PhoneNumber   string    `db:"phone_number" json:"PhoneNumber"`
	Location      string    `db:"location" json:"Location"`
	Experience    GetStatus `db:"experience" json:"ExperienceInfo"`
	AboutWork     string    `db:"about_work" json:"AboutWork"`
	IsVisible     bool      `db:"is_visible" json:"IsVisible"`
	CreatedAt     time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"UpdatedAt"`
}

type SuccessResponse struct {
	Vacancy   VacancyData `db:"vacancy" json:"VacancyInfo"`
	Responses []struct {
		ID        int           `db:"id" json:"ID"`
		Candidate InfoCandidate `db:"candidate" json:"CandidateInfo"`
		CreatedAt time.Time     `db:"created_at" json:"CreatedAt"`
		Status    GetStatus     `db:"status" json:"Status"`
	} `db:"responses"`
}

type NumberOfVacancies struct {
	Status   string `json:"Status"`
	Quantity int    `json:"Quantity"`
}

type ResponseByVac struct {
	ID      int                 `db:"id" json:"ID"`
	Vacancy VacanciesToResponse `db:"vacancy" json:"VacancyInfo"`
	Status  GetStatus           `db:"status" json:"StatusInfo"`
}

type ResponsesByVac struct {
	Status    string          `json:"Status"`
	Responses []ResponseByVac `json:"Responses"`
}

// неиспользуется
type RequestResponse struct {
	Vacancy_id int `json:"vacancy_id"`
	Status_id  int `json:"status_id"`
}

type SuccessVacancy struct {
	Employer  SuccessEmployer `db:"employer" json:"EmployerInfo"`
	Vacancies []VacancyData   `db:"vacancies" json:"VacancyInfo"`
}

type VacanciesByLimitResponse struct {
	Status        string              `json:"Status"`
	VacanciesInfo []VacancyData_Limit `json:"VacancyInfo"`
}

type VacancyData_Limit struct {
	ID          int             `db:"id" json:"ID"`
	Employer    SuccessEmployer `db:"employer" json:"EmployerInfo"`
	Name        string          `db:"name" json:"Name"`
	Price       int             `db:"price" json:"Price"`
	Email       string          `db:"email" json:"Email"`
	PhoneNumber string          `db:"phone_number" json:"PhoneNumber"`
	Location    string          `db:"location" json:"Location"`
	Experience  GetStatus       `db:"experience" json:"ExperienceInfo"`
	AboutWork   string          `db:"about_work" json:"AboutWork"`
	IsVisible   bool            `db:"is_visible" json:"IsVisible"`
	CreatedAt   time.Time       `db:"created_at" json:"CreatedAt"`
	UpdatedAt   time.Time       `db:"updated_at" json:"UpdatedAt"`
}

type VacancyData struct {
	ID          int       `db:"id" json:"ID"`
	Name        string    `db:"name" json:"Name"`
	Price       int       `db:"price" json:"Price"`
	Email       string    `db:"email" json:"Email"`
	PhoneNumber string    `db:"phone_number" json:"PhoneNumber"`
	Location    string    `db:"location" json:"Location"`
	Experience  GetStatus `db:"experience" json:"Experience"`
	AboutWork   string    `db:"about_work" json:"AboutWork"`
	IsVisible   bool      `db:"is_visible" json:"IsVisible"`
	CreatedAt   time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"UpdatedAt"`
}

type ResponseAllResponsesOnVacancy struct {
	Status    string      `json:"Status"`
	Vacancy   VacancyData ` json:"VacancyInfo"`
	Responses []struct {
		ID        int           `db:"id" json:"ID"`
		Candidate InfoCandidate `db:"candidate" json:"CandidateInfo"`
		CreatedAt time.Time     `db:"created_at" json:"CreatedAt"`
		Status    GetStatus     `db:"status" json:"StatusInfo"`
	} `db:"responses" json:"ResponseInfo"`
}

type ResponseEmployerInfo struct {
	Status   string          ` json:"Status"`
	Employer SuccessEmployer `json:"EmployerInfo"`
}

type ResponseCreateNewResponse struct {
	Response_id     int         `json:"ID"`
	Vacancy         VacancyData `json:"VacancyInfo"`
	Response_status GetStatus   `json:"ResponseStatus"`
	Status          string      `json:"Status"`
}

type ResponseAllVacancyByEmployee struct {
	Status      Ok            `json:"Status`
	Vacancies   []VacancyData `json:"VacanciesInfo"`
	Employer_id int           `json:"EmployerID"`
}

type InfoAboutAllCandidates struct {
	CandidatesInfo []InfoCandidate `json:"CandidatesInfo"`
	Status         string          `json:"Status"`
}

type ResponseAuthorization struct {
	CandidateInfo ResponseCreateCandidate `json:"CandidateInfo"`
	EmployerInfo  ResponseCreateEmployer  `json:"EmployerInfo"`
}

type ResponseCreateCandidate struct {
	Status         string        `json:"Status"`
	Candidate_Info InfoCandidate `json:"CandidateInfo"`
	Token          string        `json:"Token"`
}

type ResponseCreateEmployer struct {
	Status       Ok              `json:"Status"`
	EmployerInfo SuccessEmployer `json:"EmployerInfo"`
	Token        string          `json:"Token"`
}

type ResponseCreateNewVacancy struct {
	Status   string          `json:"Status"`
	Vacancy  VacancyData     `json:"VacancyInfo"`
	Employer SuccessEmployer `json:"EmployerInfo"`
}

type ResponseInfoByVacancy struct {
	VacancyInfo VacancyData_Limit `json:"VacancyInfo"`
	Status      string            `json:"Status"`
}

type ResponseInfoByVacancyByTimes struct {
	VacanciesInfo []VacancyData_Limit `json:"VacanciesInfo"`
	Status        string              `json:"Status"`
}

type VacancyPut struct {
	ID           int    `json:"ID"`
	VacancyName  string `json:"VacancyName"`
	Price        int    `json:"Price"`
	Email        string `json:"Email"`
	PhoneNumber  string `json:"PhoneNumber"`
	Location     string `json:"Location"`
	ExperienceId int    `json:"ExperienceId"`
	About        string `json:"About"`
	IsVisible    bool   `json:"IsVisible"`
}

type ResponseVac struct {
	VacancyName  string `json:"VacancyName"`
	Price        int    `json:"Price"`
	Email        string `json:"Email"`
	PhoneNumber  string `json:"PhoneNumber"`
	Location     string `json:"Location"`
	ExperienceId int    `json:"ExperienceId"`
	About        string `json:"About"`
	IsVisible    bool   `json:"IsVisible"`
}

type ResponseSearchVac struct {
	ID          int    `json:"ID"`
	EmployerID  int    `json:"EmployerID"`
	VacancyName string `json:"VacancyName"`
	Price       int    `json:"Price"`
	Location    string `json:"Location"`
	Experience  string `json:"Experience"`
}

type StatusInfo struct {
	Status string `json:"Status"`
	Info   string `json:"Info"`
}

type GetAllStatuses struct {
	Status string      `json:"Status"`
	Data   []GetStatus `json:"Data"`
}

type GetStatus struct {
	ID        int       `db:"id" json:"ID"`
	Name      string    `db:"name" json:"Name"`
	CreatedAt time.Time `db:"created_at"  json:"CreatedAt"`
}

type ResumeResult_slice struct {
	Id          int       `db:"id" json:"ID"`
	Experience  GetStatus `db:"experience" json:"ExperienceInfo"`
	Description string    `db:"description" json:"Description"`
	CreatedAt   time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"UpdatedAt"`
}
type ResumeResult struct {
	Resumes   []ResumeResult_slice `db:"resume" json:"ResumesInfo"`
	Candidate InfoCandidate        `db:"candidate" json:"CandidateInfo"`
}
type InfoCandidate struct {
	ID          int    `db:"id" json:"ID"`
	Name        string `db:"name" json:"Name"`
	PhoneNumber string `db:"phone_number" json:"PhoneNumber"`
	Email       string `db:"email" json:"Email"`
	Password    string `db:"password" json:"Password"`
	Status      struct {
		ID        int       `db:"id" json:"ID"`
		Name      string    `db:"name" json:"Name"`
		CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	} `db:"status" json:"StatusInfo"`
	CreatedAt time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt time.Time `db:"updated_at" json:"UpdatedAt"`
}

type GetAllFromCandidates struct {
	Candidate_Info InfoCandidate `json:"CandidateInfo"`
	Status         string        `json:"Status"`
}

type SuccessResume struct {
	Id            int       `db:"id" json:"ID"`
	Candidate_id  int       `db:"candidate_id" json:"CandidateID"`
	Experience_id int       `db:"experience_id" json:"ExperienceID"`
	Description   string    `db:"description" json:"Description"`
	CreatedAt     time.Time `db:"created_at" json:"CreatedAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"UpdatedAt"`
}

type RequestNewToken struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Claims struct {
	ID    int    `json:"uid"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

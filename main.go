package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func main() {
	config := NewConfig()

	db, err := ConnectDatabase(config)

	if err != nil {
		panic(err)
	}

	personRepository := NewPersonRepository(db)

	personService := NewPersonService(config, personRepository)

	server := NewServer(config, personService)

	server.Run()
}

// Person  ..
type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Config ..
type Config struct {
	Enabled      bool
	DatabasePath string
	Port         string
}

// NewConfig ..
func NewConfig() *Config {
	return &Config{
		Enabled:      true,
		DatabasePath: "./test.db",
		Port:         "8001",
	}
}

// ConnectDatabase ..
func ConnectDatabase(config *Config) (*sql.DB, error) {
	return sql.Open("sqlite3", config.DatabasePath)
}

// PersonRepository ..
type PersonRepository struct {
	database *sql.DB
}

// FindAll ..
func (p *PersonRepository) FindAll() []*Person {
	rows, _ := p.database.Query(`SELECT id, name, age FROM people`)
	defer rows.Close()

	people := []*Person{}

	for rows.Next() {
		var (
			id   int
			name string
			age  int
		)

		rows.Scan(&id, &name, &age)

		people = append(people, &Person{
			ID:   id,
			Name: name,
			Age:  age,
		})

	}
	return people
}

// NewPersonRepository ..
func NewPersonRepository(database *sql.DB) *PersonRepository {
	return &PersonRepository{database: database}
}

// PersonService ..
type PersonService struct {
	config     *Config
	repository *PersonRepository
}

// FindAll ..
func (p *PersonService) FindAll() []*Person {
	if p.config.Enabled {
		return p.repository.FindAll()
	}

	return []*Person{}
}

// NewPersonService ..
func NewPersonService(config *Config, repository *PersonRepository) *PersonService {
	return &PersonService{
		config:     config,
		repository: repository,
	}
}

// Server ..
type Server struct {
	config        *Config
	personService *PersonService
}

// Handler ..
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/people", s.people)

	return mux
}

// Run ..
func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.Handler(),
	}

	httpServer.ListenAndServe()
}

func (s *Server) people(w http.ResponseWriter, r *http.Request) {
	people := s.personService.FindAll()
	bytes, _ := json.Marshal(people)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

// NewServer ..
func NewServer(config *Config, service *PersonService) *Server {
	return &Server{
		config:        config,
		personService: service,
	}
}

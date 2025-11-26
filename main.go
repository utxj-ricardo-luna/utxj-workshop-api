package main

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Workshop string `json:"workshop"`
}

type Statistics struct {
	Workshops map[string]int `json:"workshops"`
	Version   string         `json:"version"`
}

var (
	students     = make(map[int]Student)
	nextID       = 1
	studentsLock sync.RWMutex
)

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/api/students", listStudents)
	e.POST("/api/students", registerStudent)
	e.GET("/stats", getStats)
	e.GET("/health", healthCheck)

	initializeStudents()

	e.Logger.Fatal(e.Start(":8080"))
}

func getStats(c echo.Context) error {
	studentsLock.RLock()
	defer studentsLock.RUnlock()

	workshops := make(map[string]int)
	for _, alumno := range students {
		workshops[alumno.Workshop]++
	}

	stats := Statistics{
		Workshops: workshops,
		Version:   "2.0",
	}

	return c.JSON(http.StatusOK, stats)
}

func listStudents(c echo.Context) error {
	studentsLock.RLock()
	defer studentsLock.RUnlock()

	lista := make([]Student, 0, len(students))
	for _, Student := range students {
		lista = append(lista, Student)
	}

	return c.JSON(http.StatusOK, lista)
}

func registerStudent(c echo.Context) error {
	var student Student
	if err := c.Bind(&student); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid  data"})
	}

	if student.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name is required"})
	}

	if student.Workshop == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Workshop is required"})
	}

	studentsLock.Lock()
	defer studentsLock.Unlock()

	student.ID = nextID
	students[nextID] = student
	nextID++

	return c.JSON(http.StatusCreated, student)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "2.0",
	})
}

func initializeStudents() {
	studentsLock.Lock()
	defer studentsLock.Unlock()

	students[1] = Student{ID: 1, Name: "Ana García", Workshop: "GitOps"}
	students[2] = Student{ID: 2, Name: "Carlos López", Workshop: "GitOps"}
	nextID = 3
}

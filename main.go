package main

import (
	"encoding/json"
	"net/http"

	"github.com/Dryluigi/golang-todos/database"
	"github.com/labstack/echo"
)

type CreateRequest struct{
	Title 		string `json:"title"`
	Description string `json:"description"`
}

type TodoResponse struct {
	Id 			int `json:"id"`
	Title		string `json:"title"`
	Description string `json:"description"`
	Done 		bool `json:"done"`
}

func main() {
	db := database.InitDb()
	defer db.Close()

	err := db.Ping()
	if err != nil{
		panic(err)
	}

	e  := echo.New()
	e.DELETE("/todos/:id", func(ctx echo.Context) error {
		id := ctx.Param("id")

		_, err := db.Exec(
			"DELETE from todos where id = ? ",
			id,
		)

		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		return ctx.String(http.StatusOK, "OK")
	})


	e.GET("/todos", func(ctx echo.Context) error {
		rows, err := db.Query("SELECT * FROM todos")
		if err != nil{
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		var res []TodoResponse
		for rows.Next(){
			var id int
			var title string
			var description string
			var done bool

			err = rows.Scan(&id, &title, &description, &done)
			if err != nil {
				return ctx.String(http.StatusInternalServerError, err.Error())
			}

			var todo TodoResponse
			todo.Id = id
			todo.Title = title
			todo.Description = description
			todo.Done = done

			res = append(res, todo)
		}

		return ctx.JSON(http.StatusOK, res)
	})

	e.POST("/todos", func(ctx echo.Context) error {
		var request CreateRequest
		json.NewDecoder(ctx.Request().Body).Decode(&request)

		_, err := db.Exec(
			"INSERT INTO todos (title, description, done) VALUES (?, ?, 0)",
			request.Title,
			request.Description,
		)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, err.Error())
		}

		return ctx.String(http.StatusOK, "OK")
	})

	e.Start(":5001")
}
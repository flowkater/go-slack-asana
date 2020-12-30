package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const asanaGetTodoListEndpoint = "https://app.asana.com/api/1.0/user_task_lists/1197488625310378/tasks?completed_since=now&opt_fields=name,assignee_status"

type Todo struct {
	Gid            string `json:"gid"`
	AssigneeStatus string `json:"assignee_status"`
	Name           string `json:"name"`
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	godotenv.Load(".env")
	port, err := determineListenAddress()
	if err != nil {
		return
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Gzip())

	e.GET("/", func(c echo.Context) error {
		var client http.Client

		req, err := http.NewRequest("GET", asanaGetTodoListEndpoint, nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		req.Header.Add("Authorization", "Bearer "+os.Getenv("TONY_ACCESS_TOKEN"))

		res, err := client.Do(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		allTodos, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		var result map[string]interface{}

		if err := json.Unmarshal(allTodos, &result); err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		arr := result["data"].([]interface{})

		todayArr := []string{}

		for _, v := range arr {
			t := v.(map[string]interface{})
			if t["assignee_status"] == "today" {
				todayArr = append(todayArr, fmt.Sprintf("%v", t["name"]))
			}
		}

		return c.JSON(http.StatusOK, todayArr)
	})

	e.Logger.Fatal(e.Start(port))
}

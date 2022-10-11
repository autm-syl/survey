package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

var mySQLXContext *sqlx.DB

/**
 * Defines a response object
 */
type Response struct {
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func main() {
	fmt.Println("Hello, ccc.")

	mySqlHost := "45.76.156.52"
	mySqlUserName := "root1"
	mySqlPassword := "@Dinhmenh1"
	mySqlDatabase := "survey_db"
	mySqlMaxOpenConnections := 20
	mySqlMaxIdleConnections := 20

	//Open a MySql infrastructure
	mySQLXContext = ConnectMySQLSqlx(mySqlHost, mySqlUserName, mySqlPassword, mySqlDatabase, mySqlMaxOpenConnections, mySqlMaxIdleConnections, 0)
	if mySQLXContext == nil {
		os.Exit(1)
	}

	err := mySQLXContext.Ping()
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		err := mySQLXContext.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	timeout := 10 * time.Second

	e := echo.New()
	// Set timeout and disable keep alive
	e.Server.SetKeepAlivesEnabled(false)
	e.Server.ReadTimeout = timeout
	e.Server.WriteTimeout = timeout

	// api
	api := e.Group("api")

	api.GET("/ping", Ping)
	api.POST("/create_type_1", Create_type_1)
	api.POST("/create_type_2", Create_type_2)
	api.POST("/create_type_3", Create_type_3)

	e.Start(":5055")
}

func Ping(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

func Create_type_1(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type_1 := Quest_type_1{}
	if err := c.Bind(&type_1); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_1 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	// insert db
	_, err := mySQLXContext.NamedExec(
		`INSERT INTO 
		quest_type_1 (session, quest_num, quest_name, answer)
		 VALUES 
		 (:session, :quest_num, :quest_name, :answer)`,
		type_1,
	)
	if err != nil {
		res := Response{
			Message: "Failed",
			Data:    err,
		}

		fmt.Println("Create_type_1 buoi", err)

		return c.JSON(http.StatusBadRequest, res)
	}

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

func Create_type_2(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type_2 := Quest_input_type_2{}
	if err := c.Bind(&type_2); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_2 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	// insert db
	resIp, err := mySQLXContext.NamedExec(
		`INSERT INTO 
		quest_type_2 (session, quest_num, quest_name)
		 VALUES 
		 (:session, :quest_num, :quest_name)`,
		type_2,
	)
	if err != nil {
		res := Response{
			Message: "Failed",
			Data:    err,
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	idQ2, err := resIp.LastInsertId()
	if err != nil {
		res := Response{
			Message: "Failed",
			Data:    err,
		}
		return c.JSON(http.StatusBadRequest, res)
	}

	for _, inside := range type_2.Inside_quest {
		// insert db
		inside.Quest_type_2_id = idQ2
		_, err = mySQLXContext.NamedExec(
			`INSERT INTO 
			inside_quest_type_2 (quest_type_2_id, quest_name, quest_answer)
			VALUES 
			(:quest_type_2_id, :quest_name, :quest_answer)`,
			inside,
		)

		if err != nil {
			fmt.Println("Create_type_12 inside buoi", err)
		}

	}

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

func Create_type_3(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type_3 := Quest_type_3{}
	if err := c.Bind(&type_3); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_3 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	// insert db
	_, err := mySQLXContext.NamedExec(
		`INSERT INTO 
		quest_type_3 (session, quest_num, quest_name, vi_tri, thoi_gian, ngoi_tren_xe, di_bo, calo, chi_phi, rui_ro, tham_gia)
		 VALUES 
		 (:session, :quest_num, :quest_name, :vi_tri, :thoi_gian, :ngoi_tren_xe, :di_bo, :calo, :chi_phi, :rui_ro, :tham_gia)`,
		type_3,
	)
	if err != nil {
		res := Response{
			Message: "Failed",
			Data:    err,
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

// (limit int, offset int, nameQuery string, costLessthan int, categoryId int64)

func GetAllItemProductByQuery(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	// limit := 1
	// offset := 1

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

func connect_db() {

}

/**
 * Initializes MySql infrastructure
 */
func ConnectMySQLSqlx(host string, userName string, password string, database string, maxOpenConnections int, maxIdleConnections int, maxConnectionLifeTime int) (db *sqlx.DB) {
	strConnection := fmt.Sprintf("%s:%s@%s(%s)/%s?%s", userName, password, "tcp", host, database, "parseTime=true")
	db, err := sqlx.Connect("mysql", strConnection)

	if err != nil {
		fmt.Println("Failed to connect to MySql", err.Error())
		return nil
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	if maxConnectionLifeTime > 0 {
		db.SetConnMaxLifetime(time.Duration(maxConnectionLifeTime))
	}
	return db
}

type Quest_type_1 struct {
	Id         int64  `db:"id" json:"Id"`
	Session    string `db:"session" json:"Session"`
	Quest_num  int    `db:"quest_num" json:"Quest_num"`
	Quest_name string `db:"quest_name" json:"Quest_name"`
	Answer     string `db:"answer" json:"Answer"`
}

type Inside_quest_type_2 struct {
	Id              int64  `db:"id" json:"Id"`
	Quest_type_2_id int64  `db:"quest_type_2_id" json:"Quest_type_2_id"`
	Quest_name      string `db:"quest_name" json:"Quest_name"`
	Quest_answer    string `db:"quest_answer" json:"Quest_answer"`
}

type Quest_type_2 struct {
	Id         int64  `db:"id" json:"Id"`
	Session    string `db:"session" json:"Session"`
	Quest_num  int    `db:"quest_num" json:"Quest_num"`
	Quest_name string `db:"quest_name" json:"Quest_name"`
}

type Quest_type_3 struct {
	Id           int64  `db:"id" json:"Id"`
	Session      string `db:"session" json:"Session"`
	Quest_num    int    `db:"quest_num" json:"Quest_num"`
	Quest_name   string `db:"quest_name" json:"Quest_name"`
	Vi_tri       string `db:"vi_tri" json:"Vi_tri"`
	Thoi_gian    int    `db:"thoi_gian" json:"Thoi_gian"`
	Ngoi_tren_xe int    `db:"ngoi_tren_xe" json:"Ngoi_tren_xe"`
	Di_bo        int    `db:"di_bo" json:"Di_bo"`
	Calo         int    `db:"calo" json:"Calo"`
	Chi_phi      int    `db:"chi_phi" json:"Chi_phi"`
	Rui_ro       string `db:"rui_ro" json:"Rui_ro"`
	Tham_gia     string `db:"tham_gia" json:"Tham_gia"`
}

type Quest_input_type_2 struct {
	Id           int64                 `db:"id" json:"Id"`
	Session      string                `db:"session" json:"Session"`
	Quest_num    int                   `db:"quest_num" json:"Quest_num"`
	Quest_name   string                `db:"quest_name" json:"Quest_name"`
	Inside_quest []Inside_quest_type_2 `db:"inside_quest" json:"Inside_quest"`
}

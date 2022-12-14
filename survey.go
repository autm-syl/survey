package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xuri/excelize/v2"
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

	timeout := 100 * time.Second

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

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
	api.GET("/all_data", GetAllItemProductByQuery)
	api.GET("/export_pdf", ExportExcel)

	// log
	api.POST("/logdata", WriteLog)
	api.GET("/all_log", GetAllWriteLog)

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

	var type_1s []Quest_type_1
	if err := c.Bind(&type_1s); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_1 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	// insert db
	var errors []interface{}
	for _, type_1 := range type_1s {
		type_1.Created_at = time.Now().Unix()
		_, err := mySQLXContext.NamedExec(
			`INSERT INTO 
		quest_type_1 (session, quest_num, quest_name, answer, created_at)
		 VALUES 
		 (:session, :quest_num, :quest_name, :answer, :created_at)`,
			type_1,
		)
		if err != nil {
			err_res := map[string]interface{}{
				"error": type_1,
			}

			fmt.Println("Create_type_1 buoi", err)

			errors = append(errors, err_res)
		}
	}

	if len(errors) != 0 {
		res := Response{
			Message: "Failed",
			Data:    errors,
		}
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

	var type_2s []Quest_input_type_2
	if err := c.Bind(&type_2s); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_2 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	var errors []interface{}
	for _, type_2 := range type_2s {
		type_2.Created_at = time.Now().Unix()
		// insert db
		resIp, err := mySQLXContext.NamedExec(
			`INSERT INTO 
			quest_type_2 (session, quest_num, quest_name, created_at)
			VALUES 
			(:session, :quest_num, :quest_name, :created_at)`,
			type_2,
		)
		if err != nil {
			err_res := map[string]interface{}{
				"error": type_2,
			}

			errors = append(errors, err_res)
		}

		idQ2, err := resIp.LastInsertId()
		if err != nil {
			err_res := map[string]interface{}{
				"error": err,
			}

			errors = append(errors, err_res)
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
	}

	if len(errors) != 0 {
		res := Response{
			Message: "Failed",
			Data:    errors,
		}
		return c.JSON(http.StatusBadRequest, res)
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

	var type_3s []Quest_type_3
	if err := c.Bind(&type_3s); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data type_3 error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	var errors []interface{}

	// insert db
	for _, type_3 := range type_3s {
		type_3.Created_at = time.Now().Unix()
		_, err := mySQLXContext.NamedExec(
			`INSERT INTO 
			quest_type_3 (session, quest_num, quest_name, vi_tri, thoi_gian, ngoi_tren_xe, di_bo, calo, chi_phi, rui_ro, tham_gia, created_at)
			 VALUES 
			 (:session, :quest_num, :quest_name, :vi_tri, :thoi_gian, :ngoi_tren_xe, :di_bo, :calo, :chi_phi, :rui_ro, :tham_gia, :created_at)`,
			type_3,
		)
		if err != nil {
			err_res := map[string]interface{}{
				"error": type_3,
			}

			errors = append(errors, err_res)
		}
	}

	if len(errors) != 0 {
		res := Response{
			Message: "Failed",
			Data:    errors,
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
	var sessions []string
	err := mySQLXContext.Select(&sessions, `SELECT session FROM survey_db.quest_type_1 group by session`)
	if err != nil {
		res := Response{
			Message: "No record was found!",
			Data:    err,
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	var result_all_data []map[string]interface{}

	for _, session := range sessions {

		var quest_1s []Quest_type_1
		_ = mySQLXContext.Select(&quest_1s, `SELECT * FROM survey_db.quest_type_1 where session = ? order by quest_num ASC`, session)

		// time
		var curentTime = time.Now().Unix()
		if len(quest_1s) != 0 {
			curentTime = quest_1s[0].Created_at
		}
		var session_time = 0
		time_sessions := strings.Split(session, "_")
		if len(time_sessions) == 2 {
			intVar, err := strconv.Atoi(time_sessions[1])
			if err == nil {
				session_time = intVar
			}
		}
		var duration = curentTime - int64(session_time/1000)
		if session_time == 0 {
			duration = 0
		}
		fmt.Printf("\ndata curentTime: %v\n", curentTime)
		fmt.Printf("\ndata session_time: %v\n", session_time)

		var quest_type_2_return []Quest_type_2_return
		var quest_2s []Quest_type_2
		_ = mySQLXContext.Select(&quest_2s, `SELECT * FROM survey_db.quest_type_2 where session = ? order by quest_num ASC`, session)
		for _, quest2 := range quest_2s {
			var quest_insite_2s []Inside_quest_type_2
			_ = mySQLXContext.Select(&quest_insite_2s, `SELECT * FROM survey_db.inside_quest_type_2 where quest_type_2_id = ?`, quest2.Id)
			x := Quest_type_2_return{
				Id:           quest2.Id,
				Session:      quest2.Session,
				Quest_num:    quest2.Quest_num,
				Quest_name:   quest2.Quest_name,
				Inside_quest: quest_insite_2s,
			}
			quest_type_2_return = append(quest_type_2_return, x)
		}

		var quest_3s []Quest_type_3
		_ = mySQLXContext.Select(&quest_3s, `SELECT * FROM survey_db.quest_type_3 where session = ? order by quest_num ASC`, session)

		var data []interface{}

		for _, quest1 := range quest_1s {

			var y interface{}
			y = quest1
			data = append(data, y)
		}
		for _, quest2 := range quest_type_2_return {
			var y interface{}
			y = quest2
			data = append(data, y)
		}
		for _, quest3 := range quest_3s {

			var y interface{}
			y = quest3
			data = append(data, y)
		}

		result_all_data = append(result_all_data, map[string]interface{}{
			"session":  session,
			"duration": duration,
			"data":     data,
		})
	}
	// step2 loop all session get quest_1, quest_2, quest_3
	// sort by quest_num

	res := Response{
		Message: "Success",
		Data:    result_all_data,
	}

	return c.JSON(http.StatusOK, res)
}

func ExportExcel(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(20)*time.Second)
	defer cancel()

	var sessions []string
	err := mySQLXContext.Select(&sessions, `SELECT session FROM survey_db.quest_type_1 group by session`)
	if err != nil {
		res := Response{
			Message: "No record was found!",
			Data:    err,
		}

		return c.JSON(http.StatusBadRequest, res)
	}

	var result_value_excel [][]string

	for _, session := range sessions {
		var result_answer []string

		var quest_1s []Quest_type_1
		_ = mySQLXContext.Select(&quest_1s, `SELECT * FROM survey_db.quest_type_1 where session = ? group by quest_num order by quest_num ASC`, session)

		// time
		var curentTime = time.Now().Unix()
		if len(quest_1s) != 0 {
			curentTime = quest_1s[0].Created_at
		}
		var session_time = 0
		time_sessions := strings.Split(session, "_")
		if len(time_sessions) == 2 {
			intVar, err := strconv.Atoi(time_sessions[1])
			if err == nil {
				session_time = intVar
			}
		}
		var duration = curentTime - int64(session_time/1000)
		if session_time == 0 {
			duration = 0
		}
		fmt.Printf("\ndata curentTime: %v\n", curentTime)
		fmt.Printf("\ndata session_time: %v\n", session_time)

		var quest_type_2_return []Quest_type_2_return
		var quest_2s []Quest_type_2
		_ = mySQLXContext.Select(&quest_2s, `SELECT * FROM survey_db.quest_type_2 where session = ? group by quest_num order by quest_num ASC`, session)
		for _, quest2 := range quest_2s {
			var quest_insite_2s []Inside_quest_type_2
			_ = mySQLXContext.Select(&quest_insite_2s, `SELECT * FROM survey_db.inside_quest_type_2 where quest_type_2_id = ?`, quest2.Id)
			x := Quest_type_2_return{
				Id:           quest2.Id,
				Session:      quest2.Session,
				Quest_num:    quest2.Quest_num,
				Quest_name:   quest2.Quest_name,
				Inside_quest: quest_insite_2s,
			}
			quest_type_2_return = append(quest_type_2_return, x)
		}

		var quest_3s []Quest_type_3
		_ = mySQLXContext.Select(&quest_3s, `SELECT * FROM survey_db.quest_type_3 where session = ? group by quest_num order by quest_num ASC`, session)

		if len(quest_1s) == 22 {
			phandau := quest_1s[0:19]
			phancuoi := quest_1s[19:]
			for _, quest1 := range phandau {
				result_answer = append(result_answer, quest1.Answer)
			}

			for _, quest2 := range quest_type_2_return {
				num := 0
				total := 0
				for _, insiteQ2 := range quest2.Inside_quest {
					answer := insiteQ2.Quest_answer
					// string to int
					answerInt, err := strconv.Atoi(answer)
					if err != nil {
						// ... skip
					} else {
						num = num + 1
						total = total + answerInt
					}
				}
				if num == 0 {
					num = 1
				}
				value := total / num
				result_answer = append(result_answer, strconv.Itoa(value))
			}
			for _, quest3 := range quest_3s {
				answer := fmt.Sprintf(`V??? tr??: %s
					Ng???i ch??? tr??n xe: %d
					??i b???: %d
					T???ng th???i gian: %d
					Calo: %d
					Chi ph?? ????a ????n: %d
					M???c an to??n: %s
					S??? tham gia c???a ph??? huynh: %s`,
					quest3.Vi_tri,
					quest3.Ngoi_tren_xe,
					quest3.Di_bo,
					quest3.Thoi_gian,
					quest3.Calo,
					quest3.Chi_phi,
					quest3.Rui_ro,
					quest3.Tham_gia)

				result_answer = append(result_answer, answer)
			}

			for _, quest1 := range phancuoi {
				result_answer = append(result_answer, quest1.Answer)
			}

			durations := fmt.Sprintf("%d", duration)
			result_answer = append(result_answer, durations)
			result_value_excel = append(result_value_excel, result_answer)
		} else {
			for _, quest1 := range quest_1s {
				result_answer = append(result_answer, quest1.Answer)
			}

			for _, quest2 := range quest_type_2_return {
				num := 0
				total := 0
				for _, insiteQ2 := range quest2.Inside_quest {
					answer := insiteQ2.Quest_answer
					// string to int
					answerInt, err := strconv.Atoi(answer)
					if err != nil {
						// ... skip
					} else {
						num = num + 1
						total = total + answerInt
					}
				}
				if num == 0 {
					num = 1
				}
				value := total / num
				result_answer = append(result_answer, strconv.Itoa(value))
			}
			for _, quest3 := range quest_3s {
				answer := fmt.Sprintf(`V??? tr??: %s
					Ng???i ch??? tr??n xe: %d
					??i b???: %d
					T???ng th???i gian: %d
					Calo: %d
					Chi ph?? ????a ????n: %d
					M???c an to??n: %s
					S??? tham gia c???a ph??? huynh: %s`,
					quest3.Vi_tri,
					quest3.Ngoi_tren_xe,
					quest3.Di_bo,
					quest3.Thoi_gian,
					quest3.Calo,
					quest3.Chi_phi,
					quest3.Rui_ro,
					quest3.Tham_gia)

				result_answer = append(result_answer, answer)
			}

			durations := fmt.Sprintf("%d", duration)
			result_answer = append(result_answer, durations)
			result_value_excel = append(result_value_excel, result_answer)
		}

	}

	// step2 loop all session get quest_1, quest_2, quest_3
	// sort by quest_num

	f := excelize.NewFile()

	nameHead := []string{"H??? t??n", "S??T", "Email", "Qu???n", "Ph?????ng", "Tr?????ng", "C??u 1", "C??u 2", "C??u 3", "C??u 4", "C??u 5", "C??u 6", "C??u 7", "C??u 8", "C??u 9", "C??u 10", "C??u 11", "C??u 12", "C??u 13", "C??u 14", "C??u 15", "C??u 16", "C??u 17", "C??u 18(TH1)", "C??u 18(TH2)", "C??u 18(TH3)", "C??u 19", "C??u 20", "C??u 21", "Th???i gian tr??? l???i"}

	cotName := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH", "AI", "AJ", "AK", "AL"}

	for i, v := range nameHead {
		position := fmt.Sprintf(`%s%d`, cotName[i], 1)
		f.SetCellValue("Sheet1", position, v)
	}

	for i, v := range result_value_excel {
		if len(v) > 30 {
			fmt.Printf("\ndata loi: %v", v)
			continue
		} else {
			for k, v2 := range v {
				position := fmt.Sprintf(`%s%d`, cotName[k], i+2)
				f.SetCellValue("Sheet1", position, v2)

			}
		}

	}
	// // Set value of a cell.
	// f.SetCellValue("Sheet1", "A2", "Hello world.")
	// Save spreadsheet by the given path.

	println("dau buoi re rech\n")
	fmt.Printf("%v", nameHead)

	if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)

		res := Response{
			Message: "Success",
			Data:    "dbrr 2222",
		}

		return c.JSON(http.StatusOK, res)
	}
	// res := Response{
	// 	Message: "Success",
	// 	Data:    "dbrr",
	// }

	// return c.JSON(http.StatusOK, res)
	return c.File("Book1.xlsx")
}

func connect_db() {
	// step1 get all session

}

func WriteLog(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var log Log_data
	if err := c.Bind(&log); err != nil {
		res := Response{
			Message: "Failed",
			Data:    "Bind data log error",
		}

		return c.JSON(http.StatusBadRequest, res)
	}
	fmt.Println("Bind data log_data", log.Log)

	// insert db
	_, err := mySQLXContext.NamedExec(
		`INSERT INTO 
		log_data (log)
		VALUES 
		(:log)`,
		log,
	)

	if err != nil {
		fmt.Println("Failed to log_data to log_data", err.Error())
	}

	res := Response{
		Message: "Success",
	}

	return c.JSON(http.StatusOK, res)
}

func GetAllWriteLog(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var logs []string
	_ = mySQLXContext.Select(&logs, `SELECT log FROM survey_db.log_data order by id DESC`)

	resS := ""
	for _, inte := range logs {
		resS = "<p>" + resS + inte + "</p>"
	}

	_ = Response{
		Message: "Success",
		Data:    resS,
	}

	return c.HTML(http.StatusOK, resS)
	// return c.JSON(http.StatusOK, res)
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

type Log_data struct {
	Id  int64  `db:"id" json:"Id"`
	Log string `db:"log" json:"Log"`
}

type Quest_type_1 struct {
	Id         int64  `db:"id" json:"Id"`
	Session    string `db:"session" json:"Session"`
	Quest_num  int    `db:"quest_num" json:"Quest_num"`
	Quest_name string `db:"quest_name" json:"Quest_name"`
	Answer     string `db:"answer" json:"Answer"`
	Created_at int64  `db:"created_at" json:"Created_at"`
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
	Created_at int64  `db:"created_at" json:"Created_at"`
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
	Created_at   int64  `db:"created_at" json:"Created_at"`
}

type Quest_input_type_2 struct {
	Id           int64                 `db:"id" json:"Id"`
	Session      string                `db:"session" json:"Session"`
	Quest_num    int                   `db:"quest_num" json:"Quest_num"`
	Quest_name   string                `db:"quest_name" json:"Quest_name"`
	Inside_quest []Inside_quest_type_2 `db:"inside_quest" json:"Inside_quest"`
	Created_at   int64                 `db:"created_at" json:"Created_at"`
}

type Quest_type_2_return struct {
	Id           int64                 `db:"id" json:"Id"`
	Session      string                `db:"session" json:"Session"`
	Quest_num    int                   `db:"quest_num" json:"Quest_num"`
	Quest_name   string                `db:"quest_name" json:"Quest_name"`
	Inside_quest []Inside_quest_type_2 `db:"inside_quest" json:"Inside_quest"`
}

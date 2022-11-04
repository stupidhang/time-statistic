package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DurationVO struct {
	Date     string `json:"date"`
	Duration int    `json:"duration"`
}

func main() {
	// 1.get mysql conn
	db, err := getDBConn()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	// 2.add log
	err = setLog()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	// 3.start server
	http.HandleFunc("/add_duration", addDuration(db))
	http.ListenAndServe(":8001", nil)
}

func getDBConn() (*sql.DB, error) {
	return sql.Open("mysql",
		"username:password@tcp(49.232.70.87:3306)/time-statistic")
}

func setLog() error {
	file, err := os.OpenFile("log/log.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	logrus.SetOutput(file)
	return nil
}

func addDuration(db *sql.DB) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			resp.WriteHeader(http.StatusForbidden)
			resp.Write([]byte("操作异常"))
			logrus.Info("操作异常")
			return
		}
		var vo DurationVO
		err := parse(req, &vo)
		if err != nil {
			resp.WriteHeader(http.StatusBadRequest)
			resp.Write([]byte("参数不合法"))
			logrus.Error(err)
			return
		}
		_, err = db.Exec("INSERT INTO duration (`date`, `duration`, `create_time`) "+
			"VALUES (?,?,?)",
			vo.Date, vo.Duration, time.Now())
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte("服务器内部异常"))
			logrus.Error(err)
			return
		}
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte("操作成功"))
		logrus.Info("操作成功")
	}
}

func parse(req *http.Request, dest interface{}) error {
	if reflect.TypeOf(dest).Kind() != reflect.Pointer {
		return errors.New("dest is not pointer")
	}
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(bytes, dest)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

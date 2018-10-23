package models

import (
	"database/sql"
	_ "database/sql/driver"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

func CString(userId int, charstr string) {
	db := DBOpen()
	PrepareAndExec(db, "INSERT INTO test_table (charstr, owner_id) VALUES (?, ?)", charstr, userId)
	db.Close()
	return
}

func ResearchString(userId int) (map[int]string, string) {
	var stmterr string
	values := map[int]string{}

	db := DBOpen()
	rows, err := db.Query("SELECT string_id, charstr FROM test_table WHERE owner_id = ?", userId)
	if err != nil {
		panic(err.Error())
	}
	db.Close()

	for rows.Next() {
		var stringId int
		var charstr string

		if err := rows.Scan(&stringId, &charstr); err != nil {
			panic(err.Error())
		}
		values[stringId] = charstr
	}

	valueCount := len(values)
	if valueCount == 0 {
		stmterr = "There is no data registered."
	} else {
		stmterr = "Count of your data: " + strconv.Itoa(valueCount)
	}
	return values, stmterr
}

func UString(updId int, charstr string) {
	db := DBOpen()
	PrepareAndExec(db, "UPDATE test_table SET charstr = ? WHERE string_id = ?", charstr, updId)
	db.Close()
	return
}

func DString(delId int) {
	db := DBOpen()
	PrepareAndExec(db, "DELETE FROM test_table WHERE string_id = ?", delId)
	db.Close()
	return
}

func SignUp(isTest bool, doSignup bool, username string, email string, password string) (bool, map[string]string) {
	var stmterr string
	var tooLongUser string
	var tooLongEmail string
	var tooLongPass string
	result := false

	if len(username) > 100 {
		tooLongUser = "Username is too long."
	}
	if len(email) > 200 {
		tooLongEmail = "E-mail address is too long."
	}
	if len(password) > 100 {
		tooLongPass = "Password is too long."
	}

	db := DBOpen()

	if doSignup == true {
		if username != "" && password != "" {
			if tooLongUser == "" && tooLongEmail == "" && tooLongPass == "" {
				var userId int
				db.QueryRow("SELECT user_id FROM users WHERE user_name = ?", username).Scan(&userId)

				if userId == 0 {
					hashPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					PrepareAndExec(db, "INSERT INTO users (user_name, email, password) VALUES (?, ?, ?)", username, email, string(hashPass))
					stmterr = "You signed up successfully!!"
					result = true
				} else {
					stmterr = "Sorry, this Username is already in use..."
				}
			} else {
				stmterr = "You failed to sign up."
			}
		} else {
			stmterr = "You must fill in both Username and Password."
		}
	}
	statements := map[string]string{"stmterr": stmterr, "tooLongUser": tooLongUser, "tooLongEmail": tooLongEmail, "tooLongPass": tooLongPass}
	if isTest == true {
		PrepareAndExec(db, "DELETE FROM users WHERE user_name = ?", "01")
		PrepareAndExec(db, "DELETE FROM test_table WHERE owner_id = ?", 1)
	}
	db.Close()
	return result, statements
}

func Login(username string, password string) (bool, int) {
	result := false
	var dbPass string
	var userId int

	db := DBOpen()
	db.QueryRow("SELECT password, user_id FROM users WHERE user_name = ?", username).Scan(&dbPass, &userId)
	db.Close()

	if bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(password)) == nil {
		result = true
	}
	return result, userId
}

func AddLoginFailCount(username string, password string) {
	db := DBOpen()
	PrepareAndExec(db, "INSERT INTO fail_login_list(fail_usr, fail_pwd) VALUES (?, ?)", username, password)
	PrepareAndExec(db, "DELETE FROM fail_login_list WHERE created < (Now() - INTERVAL 1 HOUR)")

	var fui int
	db.QueryRow("SELECT fail_usr_id FROM fail_usr_cnt WHERE fail_usr = ?", username).Scan(&fui)
	if fui == 0 {
		PrepareAndExec(db, "INSERT INTO fail_usr_cnt (fail_usr, usr_cnt_temp, usr_cnt_perm) VALUES (?, ?, ?)", username, 1, 1)
	} else {
		PrepareAndExec(db, "UPDATE fail_usr_cnt SET usr_cnt_temp = usr_cnt_temp + ?, usr_cnt_perm = usr_cnt_perm + ? WHERE fail_usr = ?", 1, 1, username)
	}
	var fpi int
	db.QueryRow("SELECT fail_pwd_id FROM fail_pwd_cnt WHERE fail_pwd = ?", password).Scan(&fpi)
	if fpi == 0 {
		PrepareAndExec(db, "INSERT INTO fail_pwd_cnt (fail_pwd, pwd_cnt_temp, pwd_cnt_perm) VALUE (?, ?, ?)", password, 1, 1)
	} else {
		PrepareAndExec(db, "UPDATE fail_pwd_cnt SET pwd_cnt_temp = pwd_cnt_temp + ?, pwd_cnt_perm = pwd_cnt_perm + ? WHERE fail_pwd = ?", 1, 1, password)
	}
	db.Close()
	return
}

func GetLoginFailCount(username string, password string) (int, string) {
	var FailCnt int
	var LastFailedTime string

	db := DBOpen()

	var FailCntUsr int
	rowsUsr, err := db.Query("SELECT COUNT(*) FROM fail_usr_cnt WHERE fail_usr = ?", username)
	if err != nil {
		panic(err.Error())
	}
	for rowsUsr.Next() {
		if err := rowsUsr.Scan(&FailCntUsr); err != nil {
			panic(err.Error())
		}
	}

	var FailCntPwd int
	rowsPwd, err := db.Query("SELECT COUNT(*) FROM fail_pwd_cnt WHERE fail_pwd = ?", password)
	if err != nil {
		panic(err.Error())
	}
	for rowsPwd.Next() {
		if err := rowsPwd.Scan(&FailCntPwd); err != nil {
			panic(err.Error())
		}
	}

	if FailCntUsr >= FailCntPwd {
		FailCnt = FailCntUsr
	} else {
		FailCnt = FailCntPwd
	}

	if FailCnt != 0 {
		db.QueryRow("SELECT max(created) FROM (SELECT * FROM fail_login_list ORDER BY created) AS tbl WHERE fail_usr = ? OR fail_pwd = ?", username, password).Scan(&LastFailedTime)
	}
	db.Close()
	return FailCnt, LastFailedTime
}

func RefreshFailCntTemp(username string, password string) {
	db := DBOpen()
	PrepareAndExec(db, "UPDATE fail_usr_cnt SET usr_cnt_temp = ? WHERE fail_usr = ?", 0, username)
	PrepareAndExec(db, "UPDATE fail_pwd_cnt SET pwd_cnt_temp = ? WHERE fail_pwd = ?", 0, password)
	db.Close()
	return
}

func DBOpen() *sql.DB {
	db, err := sql.Open("mysql", "test_user:test_pass@/test_db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func PrepareAndExec(db *sql.DB, prepare string, args ...interface{}) {
	stmt, err := db.Prepare(prepare)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	stmt.Exec(args...)
	return
}

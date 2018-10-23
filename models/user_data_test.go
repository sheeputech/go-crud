package models

import (
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"testing"
)

func TestResearchString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		tn      int
		userId  int
		stmterr string
	}{
		{tn: 0, userId: 2, stmterr: "There is no data registered."},
		{tn: 1, userId: 1, stmterr: "Count of your data: 1"},
	}
	for _, c := range cases {
		_, actual := ResearchString(c.userId)
		expected := c.stmterr
		if actual != expected {
			t.Error(strconv.Itoa(c.tn) + "(stmterr)")
			t.Errorf("actual %v\nexpected %v", actual, expected)
		}
	}
}

func TestSignUp(t *testing.T) {
	t.Parallel()

	var TooLong string
	for i := 0; i < 101; i++ {
		TooLong += "a"
	}
	var TooLongAddress string
	for i := 0; i < 257; i++ {
		TooLongAddress += "a"
	}
	cases := []struct {
		tn           int
		doSignup     bool
		username     string
		email        string
		password     string
		stmterr      string
		tooLongUser  string
		tooLongEmail string
		tooLongPass  string
	}{
		{tn: 0, doSignup: true, username: "01", password: "01", stmterr: "You signed up successfully!!"},
		{tn: 1, doSignup: true, stmterr: "You must fill in both Username and Password."},
		{tn: 2, doSignup: true, username: "03", stmterr: "You must fill in both Username and Password."},
		{tn: 3, doSignup: true, password: "04", stmterr: "You must fill in both Username and Password."},
		{tn: 4, doSignup: true, username: "default_user", password: "default_pass", stmterr: "Sorry, this Username is already in use..."},
		{tn: 5, doSignup: true, username: TooLong, email: TooLongAddress, password: TooLong, stmterr: "You failed to sign up.", tooLongUser: "Username is too long.", tooLongEmail: "E-mail address is too long.", tooLongPass: "Password is too long."},
		{tn: 6, doSignup: false},
	}
	actual := map[string]string{}
	for _, c := range cases {
		_, actual = SignUp(true, c.doSignup, c.username, c.email, c.password)
		expected := map[string]string{"stmterr": c.stmterr, "tooLongUser": c.tooLongUser, "tooLongEmail": c.tooLongEmail, "tooLongPass": c.tooLongPass}
		if actual["stmterr"] != expected["stmterr"] {
			t.Error(strconv.Itoa(c.tn) + "(stmterr)")
			t.Errorf("actual: %v /expected: %v end.", actual["stmterr"], expected["stmterr"])
		}
		if actual["tooLongUser"] != expected["tooLongUser"] {
			t.Error(strconv.Itoa(c.tn) + "(tooLongUser)")
			t.Errorf("actual: %v /expected: %v end.", actual["tooLongUser"], expected["tooLongUser"])
		}
		if actual["tooLongEmail"] != expected["tooLongEmail"] {
			t.Error(strconv.Itoa(c.tn) + "(tooLongEmail)")
			t.Errorf("actual: %v /expected: %v end.", actual["tooLongEmail"], expected["tooLongEmail"])
		}
		if actual["tooLongPass"] != expected["tooLongPass"] {
			t.Error(strconv.Itoa(c.tn) + "(tooLongPass)")
			t.Errorf("actual: %v /expected: %v end.", actual["tooLongPass"], expected["tooLongPass"])
		}
	}
}

func TestLogin(t *testing.T) {
	t.Parallel()

	cases := []struct {
		username string
		password string
		result   bool
		userId   int
	}{
		{username: "default_user", password: "default_pass", result: true, userId: 229},
		{username: "default_user", password: "", result: false},
		{username: "", password: "default_pass", result: false},
		{username: "", password: "", result: false},
	}
	actual := map[string]interface{}{}
	for i, c := range cases {
		actual["result"], actual["userId"] = Login(c.username, c.password)
		expected := map[string]interface{}{"result": c.result}
		if actual["result"] != expected["result"] {
			t.Errorf(strconv.Itoa(i) + "(result)")
			t.Errorf("actual: %t /expected: %t", actual["result"], expected["result"])
		}
	}
}

func BenchmarkCString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		CString(1, "fortestingfortestingfortestingfortestingfortesting")
	}
}

func BenchmarkResearchString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ResearchString(1)
	}
}

func BenchmarkUString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		UString(1, "ForTESTINGForTESTINGForTESTINGForTESTINGForTESTINGForTESTINGForTESTINGForTESTINGForTESTING")
	}
}

func BenchmarkAddLoginFailCount(b *testing.B) {
	for n := 0; n < b.N; n++ {
		AddLoginFailCount("testing", "testing")
	}
}

func BenchmarkGetLoginFailCount(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetLoginFailCount("testing", "testing")
	}
}

func BenchmarkRefreshFailCntTemp(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RefreshFailCntTemp("testing", "testing")
	}
}

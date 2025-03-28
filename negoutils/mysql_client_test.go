package negoutils

import (
	"fmt"
	"os"
	"testing"
)

func TestDBase_Conn(t *testing.T) {
	testStart()

	myconf := MySQLConf{
		Host:    "127.0.0.1",
		User:    "phpmyadmin",
		Passwd:  "123456",
		DbName:  "db_test",
		Charset: "utf8",
		Timeout: 5,
		Port:    3306,
	}
	db := NewDBase(myconf)
	db.SetDebug(true)
	_, err := db.Conn()
	defer db.Close()
	fmt.Println(err)
	sql := "select * from videoinfo limit 10"
	ret, err := db.FetchRows(sql)
	fmt.Println(ret, err)

	testEnd()
}

func TestDBase_InsertBatchData(t *testing.T) {
	testStart()

	myconf := MySQLConf{
		Host:       "127.0.0.1",
		User:       "phpmyadmin",
		Passwd:     "123456",
		DbName:     "db_test",
		Charset:    "utf8",
		Timeout:    5,
		Port:       3306,
		AutoCommit: true,
	}
	db := NewDBase(myconf)
	db.SetDebug(true)
	_, err := db.Conn()
	defer db.Close()
	fmt.Println(err)
	var data [][]interface{}
	fields := []string{"id1", "id2"}
	for i := 0; i < 100; i++ {
		var tmp []interface{}
		tmp = append(tmp, i)
		tmp = append(tmp, i+1)
		data = append(data, tmp)
	}
	ret, b, e := db.InsertBatchData("test", fields, data, true)
	fmt.Println(ret, b, e)

	testEnd()
}

func TestGetMySQLTableStruct(t *testing.T) {
	testStart()

	myconf := MySQLConf{
		Host:       "127.0.0.1",
		User:       "phpmyadmin",
		Passwd:     "123456",
		DbName:     "db_test",
		Charset:    "utf8",
		Timeout:    5,
		Port:       3306,
		AutoCommit: true,
	}
	db := NewDBase(myconf)
	_, err := db.Conn()
	if err != nil {
		fmt.Println(err)
		return
	}
	ret, err := GetMySQLTableStruct(db, "admin_menu")
	fmt.Println(err, ret)

	testEnd()
}

func TestGetAllMySQLTables(t *testing.T) {
	testStart()

	myconf := MySQLConf{
		Host:       "127.0.0.1",
		User:       "phpmyadmin",
		Passwd:     "123456",
		DbName:     "db_test",
		Charset:    "utf8",
		Timeout:    5,
		Port:       3306,
		AutoCommit: true,
	}
	db := NewDBase(myconf)
	_, err := db.Conn()
	if err != nil {
		fmt.Println(err)
		return
	}
	ret, err := GetAllMySQLTables(db)
	fmt.Println(err, ret)
	testEnd()
}

func TestGetMySQLAllTablesStruct(t *testing.T) {
	testStart()

	myconf := MySQLConf{
		Host:       "127.0.0.1",
		User:       "phpmyadmin",
		Passwd:     "123456",
		DbName:     "db_test",
		Charset:    "utf8",
		Timeout:    5,
		Port:       3306,
		AutoCommit: true,
	}
	db := NewDBase(myconf)
	_, err := db.Conn()
	fmt.Println(err)
	str, _ := GetMySQLAllTablesStruct(db)
	fmt.Println(str)
	testEnd()
}

func TestFormatFieldNameToGolangType(t *testing.T) {
	testStart()

	fields := []string{
		"api",
		"1user",
		"user_name",
		"user1",
		"menuName",
		"user_Name",
		"_http_status_",
	}
	for _, f := range fields {
		fmt.Fprintf(os.Stdout, "fieldName:%s\tformatFieldName:%s\n", f, FormatFieldNameToGolangType(f))
	}
	testEnd()
}

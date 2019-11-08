package cDbops

import (
	"XianZhi/models"
	"fmt"
	"testing"
)

func TestParseRow(t *testing.T) {
	Init()
	var res = new(models.GoodsCategory)
	stmt, err := dbConn.Prepare("select name, id from categories order by id")
	if err != nil {
		fmt.Println(err)
		return
	}
	row := stmt.QueryRow()
	err = ParseRow(row, res)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("id:%v,name:%v\n", res.ID, res.Name)
	}
}

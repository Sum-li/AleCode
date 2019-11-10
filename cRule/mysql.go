package cRule

import (
	"AleCode/cDbops"
	"database/sql"
	"fmt"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"strconv"
	"strings"
)

type Rule struct {
	PorG   string `json:"porg"`
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
	Field4 string `json:"field4"`
	Field5 string `json:"field5"`
	Field6 string `json:"field6"`
}

type Adapter struct {
	db    *sql.DB
	table string
}

func NewMysqlAdapter(dbconn *sql.DB, table string) *Adapter {
	return &Adapter{db: dbconn, table: table}
}

//从数据库中导出策略
func (a *Adapter) LoadPolicy(model model.Model) (err error) {
	var lines []*Rule
	if err = a.selectFields(&lines); err != nil {
		return
	}
	for _, line := range lines {
		a.loadPolicyLine(line, model)
	}
	return
}

//将策略保存到数据库中
func (a *Adapter) SavePolicy(model model.Model) (err error) {
	_, err = a.db.Exec("truncate " + a.table)
	if err != nil {
		return
	}
	for p, ast := range model["p"] {
		for _, rule := range ast.Policy {
			err = a.savePolicyLine(p, rule)
			if err != nil {
				return
			}
		}
	}
	for p, ast := range model["g"] {
		for _, rule := range ast.Policy {
			err = a.savePolicyLine(p, rule)
			if err != nil {
				return
			}
		}
	}
	return
}

//向数据库中添加一个策略
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) (err error) {
	err = a.savePolicyLine(ptype, rule)
	return
}

//删除数据库中的一个策略
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) (err error) {
	var value, sqlstr string
	value = "porg=" + `"` + ptype + `"`
	for i, v := range rule {
		value += " and " + "field" + strconv.Itoa(i+1) + "=" + `"` + v + `"`
	}
	sqlstr = fmt.Sprintf("delete from %s where %s", a.table, value)
	_, err = a.db.Exec(sqlstr)
	return
}

//删除满足要求的策略
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) (err error) {
	if len(fieldValues) == 0 {
		return
	}
	var value, sqlstr string
	value = "porg=" + `"` + ptype + `"`
	for i, v := range fieldValues {
		if i < fieldIndex {
			continue
		}
		value += " and " + "field" + strconv.Itoa(i) + "=" + `"` + v + `"`
	}
	sqlstr = fmt.Sprintf("delete from %s where %s", a.table, value)
	_, err = a.db.Exec(sqlstr)
	return
}

func (a *Adapter) loadPolicyLine(line *Rule, model model.Model) {
	var prefixLine = ", "
	var sb strings.Builder
	sb.WriteString(line.PorG)
	if len(line.Field1) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field1)
	}
	if len(line.Field2) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field2)
	}
	if len(line.Field3) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field3)
	}
	if len(line.Field4) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field4)
	}
	if len(line.Field5) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field5)
	}
	if len(line.Field6) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.Field6)
	}

	persist.LoadPolicyLine(sb.String(), model)
}

func (a *Adapter) selectFields(result *[]*Rule) error {
	sqlstr := "select * from " + a.table
	stmt, err := a.db.Prepare(sqlstr)
	if err != nil {
		fmt.Printf("select from db failed,err:%v", err)
		return err
	}
	rows, err := stmt.Query()
	if err != nil {
		fmt.Printf("select from db failed,err:%v", err)
		return err
	}
	err = cDbops.ParseRows(rows, result)
	if err != nil {
		fmt.Printf("select from db failed,err:%v", err)
		return err
	}
	return nil
}

func (a *Adapter) savePolicyLine(p string, rule []string) error {
	var key, value, sqlstr string
	key = "porg"
	value = `"` + p + `"`
	for i, v := range rule {
		key += ", " + "field" + strconv.Itoa(i+1)
		value += ", " + `"` + v + `"`
	}
	sqlstr = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", a.table, key, value)
	_, err := a.db.Exec(sqlstr)
	return err
}

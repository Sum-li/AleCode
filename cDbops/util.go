package cDbops

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

//获取结果列表的struct映射。传入struct的列表地址
func ParseRows(rows *sql.Rows, result interface{}) error {
	columns, err := rows.Columns()
	if err != nil {
		fmt.Printf("get columns failed,err:%v", err)
		return err
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		record := make(map[string]interface{})
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Printf("scan failed,err:%v", err)
			return err
		}
		for i, col := range values {
			if col != nil {
				if _, ok := col.([]uint8); ok {
					record[columns[i]] = string(col.([]uint8))
				} else {
					record[columns[i]] = col
				}
			}
		}
		records = append(records, record)
	}
	data, err := json.Marshal(&records)
	if err != nil {
		fmt.Printf("marshal failed,err:%v", err)
		return err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		fmt.Printf("unmarshal failed,err:%v", err)
	}
	return err
}

//将单个查询结果映射到struct上，传入struct的地址
func ParseRow(row *sql.Row, result interface{}) error {
	columns, err := row.Columns()
	if err != nil {
		fmt.Printf("get columns failed,err:%v", err)
		return err
	}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	record := make(map[string]interface{})
	err = row.Scan(scanArgs...)
	if err != nil {
		fmt.Printf("scan failed,err:%v", err)
		return err
	}
	for i, col := range values {
		if col != nil {
			if _, ok := col.([]uint8); ok {
				record[columns[i]] = string(col.([]uint8))
			} else {
				record[columns[i]] = col
			}
		}
	}
	data, err := json.Marshal(&record)
	if err != nil {
		fmt.Printf("marshal failed,err:%v", err)
		return err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		fmt.Printf("unmarshal failed,err:%v", err)
	}
	return err
}

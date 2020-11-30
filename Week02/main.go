package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err := updateBiz()
	if err != nil {
		fmt.Printf("err:%T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:%+v\n", err)
	}

	m, err := selectBiz()
	if err != nil {
		fmt.Printf("err:%T %v\n", errors.Cause(err), errors.Cause(err))
		fmt.Printf("stack trace:%+v\n", err)
	} else {
		fmt.Printf("get data:%+v\n", m)
	}
}

type model struct{}

// 更新业务，空数据返回错误
func updateBiz() error {
	data, err := dao()
	if err != nil {
		return errors.Wrap(err, "updateBiz error")
	}
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	// do update
	return nil
}

// 查询业务，空数据直接返回
func selectBiz() ([]model, error) {
	data, err := dao()
	if err != nil {
		return nil, errors.Wrap(err, "selectBiz error")
	}
	return data, nil
}

func dao() ([]model, error) {
	ml := make([]model, 0)
	err := sql.ErrNoRows
	if err == sql.ErrNoRows {
		return ml, nil
	}
	return nil, err
}

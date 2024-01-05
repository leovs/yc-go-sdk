// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

type JsonTime struct {
	time.Time
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	tm := strings.Replace(string(data), "\"", "", -1)
	if tm == "" {
		return
	}
	now, err := time.ParseInLocation(TimeFormat, tm, time.Local)
	*t = JsonTime{
		now,
	}
	return
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(TimeFormat))
	return []byte(formatted), nil
}

func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func Now() string {
	return time.Now().Format(TimeFormat)
}

type JsonDate struct {
	time.Time
}

func (t *JsonDate) UnmarshalJSON(data []byte) (err error) {
	tm := strings.Replace(string(data), "\"", "", -1)
	if tm == "" {
		return
	}
	now, err := time.ParseInLocation(DateFormat, strings.Split(tm, "T")[0], time.Local)
	*t = JsonDate{
		now,
	}
	return
}

func (t JsonDate) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.Format(DateFormat))
	return []byte(formatted), nil
}

func (t JsonDate) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *JsonDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

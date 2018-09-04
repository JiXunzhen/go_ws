package utils

import (
	"bytes"
	"io/ioutil"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// ConvertFunc 转换函数
type ConvertFunc func([]byte) ([]byte, error)

// GbkToUtf8 ...
func GbkToUtf8(s []byte) ([]byte, error) {
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder()))
}

// Utf8ToGbk ...
func Utf8ToGbk(s []byte) ([]byte, error) {
	return ioutil.ReadAll(transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder()))
}

// StrictConvert ...
func StrictConvert(src string, fc ConvertFunc) (string, error) {
	res, err := fc([]byte(src))
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// Convert ...
func Convert(src string, fc ConvertFunc) string {
	res, _ := StrictConvert(src, fc)
	return res
}

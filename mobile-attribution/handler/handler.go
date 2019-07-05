package handler

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	datas "mobile-attribution/dbs"

	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
)

func EnvGet(env, default_value string) string {
	value := os.Getenv(env)
	if value == "" {
		if default_value != "" {
			return default_value
		}
	}
	return value
}

func MD5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	return md5
}

type Mobile struct {
	Mobile   string `json:"mobile"`
	Operator string `json:"operator"`
	Province string `json:"province"`
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	AreaCode string `json:"areaCode"`
}

func validateMobile(mobile string) bool {
	var re = regexp.MustCompile(`(?m)((\+86)|(86))?(1)[3|4|5|6|7|8|9|]\d{9}$`)
	results := re.FindAllString(mobile, -1)
	if len(results) > 0 {
		return true
	} else {
		return false
	}
}

var headers = map[string]string{"Content-Type": "application/json"}

// OperatorType 运营商对应关系
var OperatorType = map[int]string{1: "移动", 2: "联通", 3: "电信"}

func getMobileByDB(number int) (*Mobile, error) {
	phoneModel := Phone{}
	phone := phoneModel.Get(number)
	if phone == nil {
		return nil, errors.New("mobile not found")
	}
	result := &Mobile{
		Operator: OperatorType[phone.Type],
		Province: phone.Province,
		City:     phone.City,
		ZipCode:  phone.ZipCode,
		AreaCode: phone.AreaCode,
	}
	return result, nil
}

func findNumber(number int) string {
	left, right := 0, len(datas.Phones)-1
	for left <= right {
		// fmt.Println(datas.Phones[left], datas.Phones[right])
		// fmt.Println(left, right)
		mid := (left + right) / 2
		line := datas.Phones[mid]
		num, _ := strconv.Atoi(strings.Split(line, ":")[0])
		if number == num {
			return line
		} else if num < number {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ""
}

func getMobileByFile(number int) (*Mobile, error) {
	// 使用二分法查找
	line := findNumber(number)
	// fmt.Println("line: ", line)
	if line == "" {
		return nil, errors.New("not Found")
	}
	numbers := strings.Split(line, ":")
	operator, _ := strconv.Atoi(numbers[1])
	regionID, _ := strconv.Atoi(numbers[2])
	region := datas.Regions[regionID-1]
	regions := strings.Split(region, ":")
	result := &Mobile{
		Operator: OperatorType[operator],
		Province: regions[0],
		City:     regions[1],
		ZipCode:  regions[2],
		AreaCode: regions[3],
	}
	return result, nil
}

func GetMobileAttribution(mobile string) (*Mobile, error) {
	isMobile := validateMobile(mobile)
	if !isMobile {
		return nil, errors.New("invalid mobile")
	}
	number := 0
	if strings.HasPrefix(mobile, "+86") {
		number, _ = strconv.Atoi(mobile[3:10])
	} else {
		number, _ = strconv.Atoi(mobile[0:7])
	}
	var result *Mobile
	var err error
	if DB_SOURCE == "FILE" {
		result, err = getMobileByFile(number)
	} else {
		result, err = getMobileByDB(number)
	}
	if err == nil {
		result.Mobile = mobile
	}
	// fmt.Println(result)
	return result, err
}

type ResponseError struct {
	ErrorCode string `json:"error_code"`
}

func Error(msg string) string {
	err := ResponseError{ErrorCode: msg}
	body, _ := json.Marshal(&err)
	return string(body)
}

func Response(statusCode int, body string) scf.APIGatewayProxyResponse {
	response := scf.APIGatewayProxyResponse{
		StatusCode: statusCode, Headers: headers,
		Body: body, IsBase64Encoded: false,
	}
	return response
}

func GetMobileAttributionHandler(req scf.APIGatewayProxyRequest) (scf.APIGatewayProxyResponse, error) {
	var body string
	statusCode := 200
	mobile := req.QueryString["mobile"]
	result, err := GetMobileAttribution(mobile)
	if err != nil {
		if err.Error() == "invalid mobile" {
			statusCode = 400
		} else if err.Error() == "not found" {
			statusCode = 404
		} else {
			statusCode = 500
		}
		body = err.Error()
		return Response(statusCode, Error(body)), nil
	}
	bodyByte, err := json.Marshal(result)
	if err != nil {
		body = "server error"
		statusCode = 500
	} else {
		body = string(bodyByte)
	}
	return Response(statusCode, body), nil
}

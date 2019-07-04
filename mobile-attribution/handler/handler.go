package handler

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	phoneModel := Phone{}
	phone := phoneModel.Get(number)
	result := &Mobile{
		Mobile:   mobile,
		Operator: OperatorType[phone.Type],
		Province: phone.Province,
		City:     phone.City,
		ZipCode:  phone.ZipCode,
		AreaCode: phone.AreaCode,
	}
	fmt.Println(result)
	return result, nil
}

func GetMobileAttributionHandler(req scf.APIGatewayProxyRequest) (scf.APIGatewayProxyResponse, error) {
	body := fmt.Sprintf("Hello world")
	response := scf.APIGatewayProxyResponse{StatusCode: 200, Headers: headers, Body: body, IsBase64Encoded: false}
	return response, nil
}

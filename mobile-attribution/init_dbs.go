package main

import (
	"fmt"
	"log"
	handler "mobile-attribution/handler"
	"os"
)

// var Phones = []string{"{number}:{type}:{region}"}

// var Regions = []string{"{province}:{city}:{zip_code}:{area_code}"}

var Phones = []string{
	"1300000:1:123",
	"1300001:1:123",
}

const DBFILE = "./dbs/phones.go"
const TEMP = `package dbs

// var Phones = []string{"{number}:{type}:{region}"}

// var Regions = []string{"{province}:{city}:{zip_code}:{area_code}"}

`

// WriteHead2File write head to file
func WriteHead2File(f *os.File) {
	_, err := f.WriteString(TEMP)
	if err != nil {
		log.Fatal("write temp error: %v", err)
	}
}

// WriteNumbers2File write numbers to file
func WriteNumbers2File(f *os.File) {
	fmt.Fprintln(f, "var Phones = []string{")
	phones := handler.Phones{}
	numbers, err := phones.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	for numbers.Next() {
		numbers.Scan(&phones.ID, &phones.Number, &phones.Type, &phones.RegionID)
		line := fmt.Sprintf(`"%d:%d:%d",`, phones.Number, phones.Type, phones.RegionID)
		fmt.Fprintln(f, line)
	}
	fmt.Fprintln(f, "}")
}

// WriteRegions2File write regions to file
func WriteRegions2File(f *os.File) {
	fmt.Fprintln(f, "var Regions = []string{")
	region := handler.Region{}
	regions, err := region.Fetch()
	if err != nil {
		log.Fatal(err)
	}
	for regions.Next() {
		regions.Scan(&region.ID, &region.Province, &region.City, &region.ZipCode, &region.AreaCode)
		line := fmt.Sprintf(`"%s:%s:%s:%s",`, region.Province, region.City, region.ZipCode, region.AreaCode)
		fmt.Fprintln(f, line)
	}
	fmt.Fprintln(f, "}")
}

func InitDB() {
	// open file
	f, err := os.Create(DBFILE)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	WriteHead2File(f)
	WriteNumbers2File(f)
	WriteRegions2File(f)
}

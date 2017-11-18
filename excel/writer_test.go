package excel

import (
	"log"
	"os"
	"testing"
)

func TestExecel(t *testing.T) {
	excel := New()
	sheet, err := excel.AddSheet("测试sheet页")
	if err != nil {
		log.Fatal(err)
	}
	sheet.AddHeader([]string{"姓名", "年龄", "邮箱"})
	sheet.AddRowData([]interface{}{"邓正超", 28, "cilendeng@gmail.com"})

	sheet2, err := excel.AddSheet("测试sheet页2")
	if err != nil {
		log.Fatal(err)
	}
	sheet2.AddHeaderColumns("姓名", "年龄", "邮箱")
	sheet2.AddRowData([]interface{}{"邓正超", 28, "cilendeng@gmail.com"})

	target, err := os.Create("hello.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	excel.Render(target)
}

func TestExecelStruct(t *testing.T) {
	excel := New()
	sheet, err := excel.AddSheet("测试sheet页")
	if err != nil {
		log.Fatal(err)
	}
	sheet.AddHeader([]string{"姓名", "年龄", "邮箱"})
	data := []struct {
		Name  string `excel:"姓名"`
		Age   int    `excel:"年龄"`
		Email string `excel:"邮箱"`
	}{
		{"dzc", 32, "124@132.com"},
		{Name: "aabbcc", Email: "fan@qq.com", Age: 18},
	}
	err = sheet.AddRows(data)
	if err != nil {
		log.Fatalf("AddRows error:%+v", err)
	}
	target, err := os.Create("hello2.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	excel.Render(target)
}
func TestExecelStructPointer(t *testing.T) {
	excel := New()
	sheet, err := excel.AddSheet("测试sheet页")
	if err != nil {
		log.Fatal(err)
	}
	sheet.AddHeader([]string{"姓名", "年龄", "邮箱"})
	data := []*struct {
		Name  string `excel:"姓名"`
		Age   int    `excel:"年龄"`
		Email string `excel:"邮箱"`
	}{
		{"dzc", 32, "124@132.com"},
		{Name: "aabbcc", Email: "fan@qq.com", Age: 18},
	}
	err = sheet.AddRows(data)
	if err != nil {
		log.Fatalf("AddRows error:%+v", err)
	}
	target, err := os.Create("hello2.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	excel.Render(target)
}

func TestAddAutomaticRows(t *testing.T) {
	excel := New()
	sheet, err := excel.AddSheet("测试sheet页")
	if err != nil {
		log.Fatal(err)
	}

	data := []*struct {
		Name  string `excel:"姓名123"`
		Age   int    `excel:"年龄"`
		Email string `excel:"邮箱"`
	}{
		{"dzc", 32, "124@132.com"},
		{Name: "aabbcc", Email: "fan@qq.com", Age: 18},
	}
	err = sheet.AddAutomaticRows(data)
	if err != nil {
		log.Fatalf("AddRows error:%+v", err)
	}
	target, err := os.Create("hello3.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	excel.Render(target)
}

func TestAddExcelBySteps(t *testing.T) {
	excel := New()
	sheet, err := excel.AddSheet("测试sheet页")
	if err != nil {
		log.Fatal(err)
	}

	row := sheet.AddHeader([]string{"姓名", "年龄", "邮箱"}).AddRow()
	row.Cell("zhengchao.deng").Cell(20).Cell("123@gmail.com")
	row2 := sheet.AddRow()
	data := struct {
		Name  string
		Age   int
		Email string
	}{
		Name:  "carbin-gun@gmail.com",
		Age:   19,
		Email: "carbin-gun@gmail.com",
	}
	row2.Data(data)

	row3 := sheet.AddRow()
	row3.Cells("carbin-gun-2333", 21, "999@111.com")
	target, err := os.Create("hello4.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	excel.Render(target)
}

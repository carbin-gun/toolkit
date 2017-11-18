package excel

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/tealeg/xlsx"
)

/*
	The following are samples of format samples.

	* "0.00e+00"
	* "0", "#,##0"
	* "0.00", "#,##0.00", "@"
	* "#,##0 ;(#,##0)", "#,##0 ;[red](#,##0)"
	* "#,##0.00;(#,##0.00)", "#,##0.00;[red](#,##0.00)"
	* "0%", "0.00%"
	* "0.00e+00", "##0.0e+0"
*/

const (
	Float_A = "0.00e+00"
	Float_B = "#,##0"
	Float_C = "#,##0"
	Float_D = "0.00e+00"
	Float_E = "0%"
	Float_F = "#,##0.00"
)

//Type for Error type
type Type int

const (
	//DuplicateSheet , duplicate sheet name
	DuplicateSheet Type = iota

	//OverflowSheetName ,sheetname size shoule be less than 31 characters
	OverflowSheetName
	//InvalidData , means the data input is invalid for the invocation
	InvalidData
)

//Error for error represent
type Error struct {
	Type
	Message string
}

func (e Error) Error() string {
	return fmt.Sprintf("Error, code = %d , message = %s", e.Type, e.Message)
}

//Excel represents excel file
type Excel struct {
	*xlsx.File
	FloatFormat string
}

//Sheet represents a sheet of a excel file
type Sheet struct {
	*xlsx.Sheet
}

//Row for excel row
type Row struct {
	*xlsx.Row
}

//New to construct a new Excel file
func New() *Excel {
	return &Excel{xlsx.NewFile(), Float_F}
}

//AddSheetIgnore , AddSheet() and ignore the error
func (e *Excel) AddSheetIgnore(sheetName string) *Sheet {
	s, _ := e.AddSheet(sheetName)
	return s
}

// AddSheet ,Add a new Sheet with the provided name, to a File
func (e *Excel) AddSheet(sheetName string) (*Sheet, error) {
	if _, exists := e.Sheet[sheetName]; exists {
		return nil, Error{DuplicateSheet, fmt.Sprintf("duplicate sheet name '%s'.", sheetName)}
	}
	if len(sheetName) >= 31 {
		return nil, Error{OverflowSheetName,
			fmt.Sprintf("sheet name must be less than 31 characters long.  It is currently '%d' characters long for name '%s'", len(sheetName), sheetName)}
	}
	sheet := &xlsx.Sheet{
		Name:     sheetName,
		File:     e.File,
		Selected: len(e.Sheets) == 0,
	}
	e.Sheet[sheetName] = sheet
	e.Sheets = append(e.Sheets, sheet)
	return &Sheet{sheet}, nil
}

//Render same with Write
func (e *Excel) Render(w io.Writer) {
	e.Write(w)
}

//Write , write to io.Writer to output
func (e *Excel) Write(w io.Writer) {
	e.File.Write(w)
}

//HTTPDownload , excel-attachment as an attachment
func (e *Excel) HTTPDownload(name string, w http.ResponseWriter) {
	header := w.Header()
	header.Add("Content-disposition", fmt.Sprintf("attachment;filename=%s", name))
	header.Add("Content-Type", "application/vnd.ms-excel")
	e.Write(w)
}

//AddHeader , Add first row
func (s *Sheet) AddHeader(headers []string) *Sheet {
	header := s.AddRow()
	for _, item := range headers {
		header.AddCell().SetString(item)
	}
	return s
}

//AddHeaderColumns , to add a row of header
func (s *Sheet) AddHeaderColumns(headers ...string) *Sheet {
	return s.AddHeader(headers)
}

//AddRowData ,add one row,each cell value is the correspoinding element of the input data slice
func (s *Sheet) AddRowData(data []interface{}) {
	r := s.AddRow()
	for _, item := range data {
		ref := reflect.ValueOf(item)
		r.addCellByReflect(ref)
	}
}

//AddStringRow , add a row ,all the types of the  cells are string
func (s *Sheet) AddStringRow(data []string) {
	r := s.AddRow()
	r.StringData(data)
}

// AddAutomaticRows , do the following things:
// 1. add header according to tag
// 2. add the data
func (s *Sheet) AddAutomaticRows(data interface{}) error {
	ref := reflect.ValueOf(data)
	if ref.Kind() != reflect.Slice {
		return Error{InvalidData, "Sheet.AddRows the input parameter should be a slice of struct"}
	}
	for i := 0; i < ref.Len(); i++ {
		if i == 0 {
			s.addAutomaticHeader(ref.Index(i))
		}
		s.AddRows(data)
	}
	return nil
}

//addAutomaticHeader , the val is the input element , the val should be a struct pointer or a struct
func (s *Sheet) addAutomaticHeader(val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return Error{InvalidData, "the row data input type should be struct or struct pointer"}
	}
	header := make([]string, 0, val.NumField())
	for i := 0; i < val.NumField(); i++ {
		header = append(header, val.Type().Field(i).Tag.Get("excel"))
	}
	s.AddStringRow(header)
	return nil
}

//AddRows to add rows of a slice of struct
func (s *Sheet) AddRows(data interface{}) error {
	ref := reflect.ValueOf(data)
	if ref.Kind() != reflect.Slice {
		return Error{InvalidData, "Sheet.AddRows the input parameter should be a slice of struct"}
	}
	for i := 0; i < ref.Len(); i++ {
		err := s.addReflectStructRow(ref.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

//addReflectStructRow , add the reflected value into a sheet row
func (s *Sheet) addReflectStructRow(ref reflect.Value) error {
	r := s.AddRow()
	//if the ref is a struct pointer , use ref.Elm() to return the element it points to
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	if ref.Kind() != reflect.Struct {
		return Error{InvalidData, "the row data input type should be struct or struct pointer"}
	}
	for i := 0; i < ref.NumField(); i++ {
		r.addCellByReflect(ref.Field(i))
	}
	return nil
}

//AddStructRow , a struct means a row
func (s *Sheet) addStructRow(data interface{}) error {
	ref := reflect.ValueOf(data)
	return s.addReflectStructRow(ref)
}

//AddRow ...
func (s *Sheet) AddRow() *Row {
	return &Row{s.Sheet.AddRow()}
}

//StringData , add row for all string
func (r *Row) StringData(items []string) {
	for _, item := range items {
		r.AddCell().SetString(item)
	}
}

//Data ,write data into a row according to reflect thing
func (r *Row) Data(item interface{}) {
	data := reflect.ValueOf(item)
	r.reflectData(data)
}

//Cell add a new cell
func (r *Row) Cell(data interface{}) *Row {
	r.addCellByReflect(reflect.ValueOf(data))
	return r
}

//Cells to add a couple of cells ,and set the value
func (r *Row) Cells(data ...interface{}) *Row {
	for _, item := range data {
		r.addCellByReflect(reflect.ValueOf(item))
	}
	return r
}

//reflectData , add a row of data ,the data is represented by the ref
func (r *Row) reflectData(ref reflect.Value) error {
	//if the ref is a struct pointer , use ref.Elm() to return the element it points to
	if ref.Kind() == reflect.Ptr {
		ref = ref.Elem()
	}
	if ref.Kind() != reflect.Struct {
		return Error{InvalidData, "the row data input type should be struct or struct pointer"}
	}
	for i := 0; i < ref.NumField(); i++ {
		r.addCellByReflect(ref.Field(i))
	}
	return nil
}

func (r *Row) addCellByReflect(val reflect.Value) {
	switch val.Kind() {
	case reflect.Int, reflect.Int64, reflect.Int32:
		r.AddCell().SetInt64(val.Int())
	case reflect.String:
		r.AddCell().SetString(val.String())
	case reflect.Float32, reflect.Float64:
		r.AddCell().SetFloatWithFormat(val.Float(), Float_F)
	}
}

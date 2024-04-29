package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)
type Postgress struct {
	db *sql.DB
}

type Storage interface {
	CreateStudentStore(acc *requestStudentId) (error)
	AttendanceStore(acc *requestEsp)(error)
	AllAttendance()([]requestAll,error )
	NewAttendanceEntires(tag time.Time)([]requestAll,error)


}

func InitEnv(){
	os.Setenv("DB_HOST","localhost");
	os.Setenv("DB_PORT","5432");
	os.Setenv("DB_NAME","postgres");
	os.Setenv("DB_PASSWORD","root")
	fmt.Println("ALERT THE DB_HOST IS SET TO localhost")
}

func Newpostgress() (*Postgress,error){
	connStr := fmt.Sprintf("host=%s port=%s user=postgres dbname=%s password=%s sslmode=disable",os.Getenv("DB_HOST"),os.Getenv("DB_PORT"),os.Getenv("DB_NAME"),os.Getenv("DB_PASSWORD"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil ,err
	}

	if err:=db.Ping();err!=nil{
		return nil,err
	}
	return &Postgress{db:db},nil



}


func (e *Postgress) CreateStudentStore(acc *requestStudentId) (error) {
	query:=`INSERT INTO studentdetails(rollno,fname,lname,branch,year,created_at) values($1,$2,$3,$4,$5,$6)`
	
	_,err := e.db.Query(query,acc.RollNo,acc.Fname,acc.Lname,acc.Branch,acc.Year,time.Now().UTC())
	
	if err!=nil{
		return err
	}

	return nil
}


func (s *Postgress) init() error{

	query:=`CREATE TABLE IF NOT EXISTS studentdetails (
		rollno char(9),
		fname varchar(50),
		lname varchar(50),
		branch varchar(15),
		year int,
		created_at timestamp
		)`
	
	_,err:=s.db.Exec(query)
	
	if err!=nil{
		return err
	}

	
	query1:=`CREATE TABLE IF NOT EXISTS attendance (
		rollno char(9),
		subject varchar(50),
		created_at timestamp
		)`

	_,err1:=s.db.Exec(query1)
	
	if err1!=nil{
		return err1
	}


	return nil
	

}

func (s * Postgress) AttendanceStore(acc *requestEsp)(error){

	query:=`INSERT INTO attendance(rollno,subject,created_at) values($1,$2,$3)`
	
	_,err := s.db.Query(query,acc.RollNo,acc.Subject,time.Now().UTC())
	
	if err!=nil{
		return err
	}

	return nil
}

func (s* Postgress) AllAttendance()( []requestAll,error){
	var acc []requestAll;
	rows,err := s.db.Query(`SELECT * FROM attendance`)
	if err!=nil{

		return nil,fmt.Errorf("bad-url")
	}
	defer rows.Close()

	requestAllStruct:= &requestAll{}
	
	for rows.Next(){
		err:=ScanIntoStructAttendance(rows,requestAllStruct)
		if err!=nil{
			return nil,err
		}
		row1:=s.db.QueryRow(`SELECT fname,lname,branch,year FROM studentdetails where rollno = $1`,requestAllStruct.RollNo)
	
		err2:=ScanIntoStructdetails(row1,requestAllStruct)

		if err2!=nil{
			return nil,err2
		}

		acc = append(acc,*requestAllStruct)


	}

	return acc,err
}

func (s* Postgress) NewAttendanceEntires(tag time.Time)([]requestAll,error){
	
	var acc []requestAll;
	rows,err := s.db.Query(`SELECT * FROM attendance where created_at > $1`,tag)
	if err!=nil{

		return nil,fmt.Errorf("bad-url")
	}
	defer rows.Close()

	requestAllStruct:= &requestAll{}
	
	for rows.Next(){
		err:=ScanIntoStructAttendance(rows,requestAllStruct)
		if err!=nil{
			return nil,err
		}
		row1:=s.db.QueryRow(`SELECT fname,lname,branch,year FROM studentdetails where rollno = $1`,requestAllStruct.RollNo)
	
		err2:=ScanIntoStructdetails(row1,requestAllStruct)

		if err2!=nil{
			return nil,err2
		}

		acc = append(acc,*requestAllStruct)


	}

	return acc,err
	
}


func ScanIntoStructAttendance(rows *sql.Rows,AccountStruct *requestAll) error{
	
	err:=rows.Scan(
		&AccountStruct.RollNo,
		&AccountStruct.Subject,
		&AccountStruct.CreatedAt,
	)
	if err!=nil{
	 return	err
	}
	return nil
}
func ScanIntoStructdetails(rows *sql.Row,AccountStruct *requestAll) error{
	
	err:=rows.Scan(
		&AccountStruct.Fname,
		&AccountStruct.Lname,
		&AccountStruct.Branch,
		&AccountStruct.Year,
	)
	if err!=nil{
	 return	err
	}

	return nil
}
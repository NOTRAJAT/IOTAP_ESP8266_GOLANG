package main

import (
	"encoding/json"
	"fmt"
	"iotap_mini_proc/templ"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ApiServer struct {
	ListenAddr string
	// loadRequest time.Time
	Store Storage

}

type requestAll struct{
	RollNo string `json:"Roll_no"`
	Subject string `json:"Subject"`
	Fname string  `json:"Fname"`
	Lname string  `json:"Lname"`
	Branch string  `json:"Branch"`
	Year uint8  `json:"Year"`
	CreatedAt time.Time `json:"Createdat"`
	
}

type requestStudentId struct{
	RollNo string `json:"Roll_no"`
	
	Fname string  `json:"Fname"`
	Lname string  `json:"Lname"`
	Branch string  `json:"Branch"`
	Year uint8  `json:"Year"`
	
}

type requestEsp struct{
	RollNo string `json:"Roll_no"`
	Subject string `json:"Subject"`
}

type ApiError struct{
	Error string
}

type ApiFunc func (http.ResponseWriter,*http.Request) error

func setLoadReqCookie(w http.ResponseWriter){
	cook := http.Cookie{Name: "loadRequest",Value: time.Now().UTC().String(),Expires: time.Now().AddDate(0, 0, 1), SameSite: http.SameSiteLaxMode, Domain: "localhost" ,Secure: false}
	http.SetCookie(w,&cook)
}
func getLoadReqCookie(r *http.Request) (time.Time, error){
	tag,err:=r.Cookie("loadRequest")
	if err!=nil{
		return time.Now() ,err
	}
	t,_:=time.Parse("2006-01-02 15:04:05.999999999 -0700 MST",tag.Value)

	return t , nil
}

func Runserver(s *ApiServer) {
	fmt.Println("Sever on ", s.ListenAddr)
	fmt.Println("network ip", GetDhcpIp())
	printGoApi()
	router := mux.NewRouter()
	// fs:= http.FileServer(http.Dir("/templ/css"))
	// http.Handle("/static",fs)

	router.HandleFunc("/espurl",makeHttpHandlefunc(s.EspAttendenceRequest))
	router.HandleFunc("/createstudent",makeHttpHandlefunc(s.handleCreateStudent))

	router.HandleFunc("/test",makeHttpHandlefunc(s.loadTestPage))
	router.HandleFunc("/test/load",makeHttpHandlefunc(s.loadData))
	router.HandleFunc("/test/loadNewEntries",makeHttpHandlefunc(s.loadDataNewEntries))
	router.HandleFunc("/static/{id}",func(w http.ResponseWriter, r* http.Request){
		val:=mux.Vars(r)["id"]
		http.ServeFile(w,r,fmt.Sprintf("./templ/css/%v",val))
	})



	http.ListenAndServe(s.ListenAddr, router)
}

func (s * ApiServer) handleCreateStudent(w http.ResponseWriter, r* http.Request) error{
	
	if(r.Method=="POST"){
	reqStudentsStruct:= &requestStudentId{}

	if err:=ReadJson(r,reqStudentsStruct);err!=nil{
		return fmt.Errorf("invalid json format")
	}

	// StudentIdStruct:= &StudentId{}
	// createStudentIdStruct(reqStudentsStruct,StudentIdStruct)
	if err:=s.Store.CreateStudentStore(reqStudentsStruct);err!=nil{
		return err;
	}
	WriteJson(w,200,reqStudentsStruct)
}else {
	return fmt.Errorf("invalid method request")
}
	return nil
}

func (s *ApiServer)loadTestPage(w http.ResponseWriter, r* http.Request) error{
	
	// WriteJson(w,200,res)
	setLoadReqCookie(w)
	templ.Hello().Render(r.Context(),w)
	return nil
}
func (s *ApiServer)loadData(w http.ResponseWriter, r* http.Request) error{

		var all string
		res,err:=s.Store.AllAttendance()
		if err!=nil{
			return err
		}
		for _,first:=range res{
			all+=fmt.Sprintf(` <tr >
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
			<td class="p-2 border border-gray-800 text-sm">%s</td>
		</tr>`,first.RollNo,first.Fname,first.Lname,first.Branch,strconv.Itoa(int(first.Year)),first.Subject,first.CreatedAt.Local().Format(time.RFC850))
		}
		w.Write([]byte(all+`<tr id="append" class"hidden"> </tr>`))
	return nil
}
func (s *ApiServer)loadDataNewEntries(w http.ResponseWriter, r* http.Request) error{
val,err:=getLoadReqCookie(r)
if err!=nil{
	return err
}
var all string
	res,err:=s.Store.NewAttendanceEntires(val)
	if err!=nil{
		return err
	}
	setLoadReqCookie(w)
	// WriteJson(w,200,res)
	for _,first:=range res{
		all+=fmt.Sprintf(` <tr >
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>
        <td class="p-2 border border-gray-800 text-sm">%s</td>

    </tr>`,first.RollNo,first.Fname,first.Lname,first.Branch,strconv.Itoa(int(first.Year)),first.Subject,first.CreatedAt.Local().Format(time.RFC850))
	}
	w.Write([]byte(all+`<tr id="append" class"hidden"> </tr>`))
	return nil
}
func (s *ApiServer) EspAttendenceRequest(w http.ResponseWriter, r *http.Request) error{
	if(r.Method=="POST"){
	    studentid := &requestEsp{}
		if err:=ReadJson(r,studentid);err !=nil {
			return fmt.Errorf("invalid id")
		}
		if err:=s.Store.AttendanceStore(studentid);err!=nil{
			return err
		}
		WriteJson(w,200,studentid)
		// fmt.Print(studentid)
		
	}else {
		return fmt.Errorf("invalid method request")
	}
	return nil
}





func printGoApi(){
	fmt.Println(`   
 ________  ________          ________  ________  ___     
|\   ____\|\   __  \        |\   __  \|\   __  \|\  \    
\ \  \___|\ \  \|\  \       \ \  \|\  \ \  \|\  \ \  \   
 \ \  \  __\ \  \\\  \       \ \   __  \ \   ____\ \  \  
  \ \  \|\  \ \  \\\  \       \ \  \ \  \ \  \___|\ \  \ 
   \ \_______\ \_______\       \ \__\ \__\ \__\    \ \__\
    \|_______|\|_______|        \|__|\|__|\|__|     \|__|
																`)	
                                              
}
func GetDhcpIp() string{
	netInterfaceAddresses, _ := net.InterfaceAddrs()
	for _,val:= range netInterfaceAddresses{
		ip := val.String()
		
		if(strings.HasPrefix(ip,"192") || strings.HasPrefix(ip,"172")|| strings.HasPrefix(ip,"10")){
			return ip[:strings.Index(ip,"/")]
		}
	}
	return ""

}

func NewApiServerAddr(ListenAddr string,store Storage) *ApiServer { //constuctor functions
	instance := &ApiServer{
		ListenAddr: ListenAddr,
		Store: store,
		}
	return instance
}

func logline(r* http.Request){
	host:=r.RemoteAddr
	log.Println(r.URL,host[:strings.Index(host,":")],r.UserAgent())
}

func WriteJson(w http.ResponseWriter,status int, args any) error{
	w.Header().Add("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(args)

}
func makeHttpHandlefunc(e ApiFunc) http.HandlerFunc {
	
	return func (w http.ResponseWriter,r *http.Request){
		logline(r)
		if err:=e(w,r);err!=nil{
								//
			WriteJson(w,http.StatusBadRequest,ApiError{Error: err.Error()})					
		}

	}

}

func  ReadJson[T requestStudentId| requestEsp](r* http.Request,e *T) (error) {
		if err := json.NewDecoder(r.Body).Decode(e); err!=nil{
		return err
	}
	defer r.Body.Close()
	return nil
}
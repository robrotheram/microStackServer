package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"test/app/dbc"
)

type Audit struct {
	Id       int    `json:"id"`
	TITLE    string `json:"title"`
	TEXT     string `json:"text"`
	USERNAME string `json:"username"`
	VMUUID   string `json:"vmuuid"`
	DATE     int64  `json:"date"`
}

type VMINFO struct {
	UUID      string `json:"uuid"`
	TITLE     string `json:"Title"`
	OS        string `json:"OS"`
	CORES     int    `json:"Cores"`
	MEMORRY   int    `json:"Memory"`
	DISK      string `json:"Disk"`
	IPADDRESS int64  `json:"Ipaddress"`
}

func setAudit(a Audit) {
	log.Println(a.TITLE, a.TEXT, a.DATE, a.VMUUID)
	db := dbc.Connect()

	if db != nil {
		defer db.Close()
		stmt, err := db.Prepare("INSERT INTO Audit (AuditTitle, AuditText, AuditDate, AuditVmid) values (?,?,?,?) ")
		if err != nil {
			log.Print(err.Error()) // proper error handling instead of panic in your app
			log.Print("Failed to insert post into database")
		} else {
			res, err := stmt.Exec(a.TITLE, a.TEXT, a.DATE, a.VMUUID)
			if err != nil {
				log.Print(err.Error()) // proper error handling instead of panic in your app
				log.Print("Failed to insert post into database")
			}
			id, err := res.LastInsertId()
			if err != nil {
				log.Print("Failed to get last insert id")
			}
			log.Print("Successfully Inserted Post %i Into the Database", id)
		}
	}
}

func creatVM(vm VMINFO) {
	db := dbc.Connect()
	defer db.Close()
	stmt, _ := db.Prepare("INSERT INTO VMS (idVMS, Title, OS, Cores, Memory, Disk, Ipaddress) values (?,?,?,?,?,?,?) ")
	res, err := stmt.Exec(vm.UUID, vm.TITLE, vm.OS, vm.CORES, vm.MEMORRY, vm.DISK, vm.IPADDRESS)
	if err != nil {
		log.Print(err.Error()) // proper error handling instead of panic in your app
		log.Print("Failed to insert post into database")
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print("Failed to get last insert id")
	}
	log.Print("Successfully Inserted Post %i Into the Database", id)
}

func GetAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	decoder := json.NewDecoder(r.Body)
	var vms VMINFO
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.UUID != "" {
			json.NewEncoder(w).Encode(vmAUDIT(vms.UUID))
		}
	}
}

/*
func aetAudit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)

	bodyBytes, err2 := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

	decoder := json.NewDecoder(r.Body)
	var vm VMINFO
	err := decoder.Decode(&vm)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		fmt.Println(vm)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vmAUDIT(vm.idVMS))
	}
}

*/
func vmAUDIT(id string) []Audit {
	db := dbc.Connect()
	defer db.Close()
	var audits []Audit

	fmt.Println("ID=:" + id)

	rows, err := db.Query("SELECT AuditTitle, AuditText, AuditDate, AuditVmid  FROM Audit inner join VMS on Audit.AuditVmid = VMS.idVMS  where Audit.AuditVmid = ?", id)
	//.Scan(audit.TITLE, audit.TEXT, audit.DATE, audit.VMUUID)
	defer rows.Close()
	for rows.Next() {
		var audit Audit
		err = rows.Scan(&audit.TITLE, &audit.TEXT, &audit.DATE, &audit.VMUUID)
		if err != nil {
			panic(err)
		}
		audits = append(audits, audit)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return audits
}

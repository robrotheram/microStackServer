package Config

import (
	"encoding/json"
	"fmt"
	"os"

	libvirt "github.com/libvirt/libvirt-go"		
)

type Configuration struct {
	URL     string `json:"URL"`
	DB_HOST string `json:"DB_HOST"`
	DB_PORT string `json:"DB_PORT"`
	DB_USER string `json:"DB_USERNAME"`
	DB_PASS string `json:"DB_PASSWORD"`
	DB_NAME string `json:"DB_NAME"`
}

var Conf = Configuration{}
var Conn *libvirt.Connect

func ReadConfig(filePath string) *Configuration {
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Conf)
	if err != nil {
		fmt.Println("error:", err)
		return &Conf
	} else {
		return &Conf
	}
}


func LibConnect(){
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		fmt.Println(err.Error())
	}
	Conn = conn;
}
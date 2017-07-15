package handlers

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
	//	"log"

	"test/app/Background"
	"test/app/Config"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/ssh"
	yaml "gopkg.in/yaml.v2"
)

type diskInfo struct {
	Format     string `json:"format"`
	Path       string `json:"filename"`
	Size       uint64 `json:"virtual-size"`
	ActualSize uint64 `json:"actual-size"`
}

type VMS struct {
	UUID      string     `json:"uuid"`
	HostName  string     `json:"hostname"`
	Name      string     `json:"name"`
	OS        string     `json:"OS"`
	Cores     uint       `json:"cores"`
	Memory    uint64     `json:"memory"`
	Disks     []diskInfo `json:"disks"`
	DiskSize  string     `json:"disk"`
	Ipaddress string     `json:"ipaddress"`
	State     string     `json:"state"`
	FreeMem   uint64     `json:freemem`
	Load      float64    `json:load`
}

type userinit struct {
	Name              string   `yaml:"name"`
	Groups            string   `yaml:"group"`
	Shell             string   `yaml:"shell"`
	Sudo              []string `yaml:"sudo"`
	Password          string   `yaml:"password"`
	Chpassword        string   `yaml:"chpassword"`
	Ssh_pwauth        bool     `yaml:"ssh_pwauth"`
	Sshauthorizedkeys []string `yaml:"ssh-authorized-keys"`
}
type cloudConfig struct {
	Hostname         string     `json:"hostname"`
	Manage_etc_hosts bool       `json:"mange_etc_hosts"`
	Users            []userinit `json:"users"`
}

func CommandRunner(cmdName string, cmdArgs []string) (err error, result bool) {
	fmt.Println("running commandL: %s", cmdName)
	fmt.Println(cmdArgs)
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("OutPUT out | %s\n", scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		return err, false
	}
	for scanner.Scan() {
		fmt.Printf("CMD out | %s\n", scanner.Text())
	}

	err = cmd.Wait()
	if err != nil {
		return err, false
	}
	return nil, true

}

func Info(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info())
}

func Start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var vms VMS
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.UUID != "" {
			start(vms.UUID)
			setAudit(Audit{TITLE: "VMS", TEXT: "VM start", DATE: time.Now().Unix(), VMUUID: vms.UUID})
			json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "VM Started"})
		}
	}
}

func Stop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var vms VMS
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.UUID != "" {
			shutdown(vms.UUID)
			json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "VM stoped"})
		}
	}
}

func Pause(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var vms VMS
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.UUID != "" {
			pause(vms.UUID)
			json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "VM Paused"})
		}
	}
}

func Resume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var vms VMS
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.UUID != "" {
			resume(vms.UUID)
			json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "VM Resumed"})
		}
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var vms VMS
	err := decoder.Decode(&vms)
	if err != nil {
		json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
	} else {
		if vms.HostName != "" &&
			vms.Cores != 0 &&
			vms.Memory != 0 &&
			vms.DiskSize != "" {
			result, err := create(vms.HostName, vms.DiskSize, int(vms.Memory), int(vms.Cores))
			if err != nil {
				json.NewEncoder(w).Encode(HttpResp{Status: 500, Description: "Failed to registered", Error: err.Error()})
			} else {
				json.NewEncoder(w).Encode(HttpResp{Status: 200, Description: "VM Created", UUID: result})
			}
		}
	}

	//start(u1.String())
	//--bridge=$NETWORK -m $macadd*/
}

func create(name string, diskSize string, ram int, cpu int) (result string, err error) {
	u1 := uuid.NewV4()
	dir := "data/" + u1.String()
	os.Mkdir(dir, 0777)

	pub := MakeSSHKeyPair((dir + "/pub.key"), (dir + "/private.key"))
	fmt.Println(pub)

	cd := cloudConfig{
		Hostname:         name,
		Manage_etc_hosts: true,
		Users: []userinit{
			userinit{
				Name:              "ubuntu",
				Groups:            "sudo",
				Shell:             "/bin/bash",
				Password:          "mypassword",
				Chpassword:        "expire:True",
				Sudo:              []string{"ALL=(ALL) NOPASSWD:ALL"},
				Ssh_pwauth:        true,
				Sshauthorizedkeys: []string{pub},
			},
		},
	}

	//fmt.Println(cd)

	d, err := yaml.Marshal(cd)
	if err != nil {
		return "", err
	}
	fmt.Printf(string(d))

	file, err := os.OpenFile((dir + "/cloudconfig"), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()

	wr := bufio.NewWriter(file)
	fmt.Fprintf(wr, "#cloud-config\n")
	fmt.Fprintf(wr, "%s\n", string(d))
	wr.Flush()

	CommandRunner("cloud-localds", []string{(dir + "/seed.img"), (dir + "/cloudconfig")})
	CommandRunner("qemu-img", []string{"convert", "-O", "qcow2", "/home/robert/VMCreator/images/trusty-server-cloudimg-amd64-disk1.img", (dir + "/" + u1.String() + ".qcow2")})
	CommandRunner("qemu-img", []string{"resize", (dir + "/" + u1.String() + ".qcow2"), diskSize})

	disk := (dir + "/" + u1.String() + ".qcow2")
	cdPath := (dir + "/seed.img")

	CommandRunner("virt-install", []string{
		"--vnc",
		"--noautoconsole",
		"--noreboot",
		"--uuid=" + u1.String(),
		"--name=" + name,
		"--ram=" + strconv.Itoa(ram),
		"--vcpu=" + strconv.Itoa(cpu),
		"--disk", "path=" + disk + ",format=qcow2",
		"--cdrom=" + cdPath,
		"--boot=hd",
		"--livecd",
	})
	setAudit(Audit{TITLE: "VMS", TEXT: "VM Created", DATE: time.Now().Unix(), VMUUID: u1.String()})
	creatVM(VMINFO{UUID: u1.String(), TITLE: name, CORES: cpu, MEMORRY: ram, DISK: diskSize})
	return u1.String(), nil
}
func delete() {}

func pause(uuid string) {
	conn := Background.LibConnect()
	defer conn.Close()
	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err == nil {
		err := dom.Suspend()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	setAudit(Audit{TITLE: "VMS", TEXT: ("VM PAUSED"), DATE: time.Now().Unix(), VMUUID: uuid})
}
func resume(uuid string) {
	conn := Background.LibConnect()
	defer conn.Close()
	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err == nil {
		err := dom.Resume()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	setAudit(Audit{TITLE: "VMS", TEXT: ("VM RESUMED"), DATE: time.Now().Unix(), VMUUID: uuid})
}
func shutdown(uuid string) {
	conn := Background.LibConnect()
	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err == nil {
		err := dom.Shutdown()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	conn.Close()
	setAudit(Audit{TITLE: "VMS", TEXT: ("VM shutdown"), DATE: time.Now().Unix(), VMUUID: uuid})
}
func start(uuid string) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		fmt.Println(err.Error())
	}
	dom, err := conn.LookupDomainByUUIDString(uuid)
	if err == nil {
		err := dom.Create()
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	conn.Close()
}

func domainStateToString(ds libvirt.DomainState) string {
	switch ds {
	case libvirt.DOMAIN_NOSTATE:
		return "NOSTATE"
	case libvirt.DOMAIN_RUNNING:
		return "RUNNING"
	case libvirt.DOMAIN_BLOCKED:
		return "BLOCKED"
	case libvirt.DOMAIN_PAUSED:
		return "PAUSED"
	case libvirt.DOMAIN_SHUTDOWN:
		return "SHUTDOWN"
	case libvirt.DOMAIN_CRASHED:
		return "CRASHED"
	case libvirt.DOMAIN_PMSUSPENDED:
		return "PMSUSPENDED"
	case libvirt.DOMAIN_SHUTOFF:
		return "SHUTOFF"
	default:
		return "ERROR"
	}
}

func GetVNC(luuid string) (string, int) {
	var conn = Config.Conn
	if conn != nil {
		doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		if err == nil {
			for _, dom := range doms {
				uuid, _ := dom.GetUUIDString()
				if uuid == luuid {
					xmldoc, _ := dom.GetXMLDesc(0)
					domcfg := &libvirtxml.Domain{}
					_ = domcfg.Unmarshal(xmldoc)
					return domcfg.Devices.Graphics[0].Listen, domcfg.Devices.Graphics[0].Port
				}
			}
		}
	}
	return "", 0
}

func info() []VMS {
	var conn = Config.Conn
	if conn != nil {
		domArr := []VMS{}
		doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
		domsA, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
		domss := append(doms, domsA...)
		if err != nil {
			//fmt.Println(err.Error())
		}
		var dARR []*libvirt.Domain
		for _, dom := range domss {
			dARR = append(dARR, &dom)
			uuid, err := dom.GetUUIDString()
			name, err := dom.GetName()
			info, err := dom.GetInfo()
			ram := info.MaxMem
			cpu := info.NrVirtCpu

			state, _, err := dom.GetState()

			xmldoc, err := dom.GetXMLDesc(0)
			domcfg := &libvirtxml.Domain{}
			err = domcfg.Unmarshal(xmldoc)
			if err == nil {
				vm := VMS{UUID: uuid, HostName: name, Cores: cpu, Memory: ram, Ipaddress: "", State: domainStateToString(state)}
				if state == libvirt.DOMAIN_RUNNING {
					ms, err := dom.MemoryStats(1, 0)
					if err == nil {
						for _, mem := range ms {
							if mem.Tag == 6 {
								//vm.FreeMem = mem.Val
							}
						}
					}
					cpustats, err := dom.GetCPUStats(0, 0, 0)
					fmt.Println(cpustats)
				}
				diskSource := []diskInfo{}
				for _, disk := range domcfg.Devices.Disks {
					if disk.Source != nil {
						var info diskInfo
						cmd := exec.Command("qemu-img", "info", "--output=json", disk.Source.File)
						out, err := cmd.CombinedOutput()
						if err != nil {
							//fmt.Println("Error" + string(out))
						}
						err = json.Unmarshal(out, &info)
						if err == nil {
							diskSource = append(diskSource, info)
						}
					}
				}
				vm.Disks = diskSource
				vm.Load = Background.VTS[name].Load
				vm.FreeMem = Background.VTS[name].Mem
				domArr = append(domArr, vm)
			}
		}
		return domArr
	}
	return []VMS{}
}

func Round(val float64, places int) (newVal float64) {
	var round float64
	roundOn := 0.5
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

/*

type VMS struct {
	UUID     string `json:"uuid"`
	Title     string `json:"title"`
	OS        string `json:"os"`
	Cores     int    `json:"cores"`
	Memory    int    `json:"memory"`
	Disk      int    `json:"disk"`
	Ipaddress string `json:"ipaddress"`
	State     string `json:"state"`
}*/
func healthcheck() {}

// MakeSSHKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
func MakeSSHKeyPair(pubKeyPath, privateKeyPath string) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Print(err.Error())
	}
	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		fmt.Print(err.Error())
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		fmt.Print(err.Error())
	}
	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0644)
	if err != nil {
		fmt.Print(err.Error())
	}
	return string(ssh.MarshalAuthorizedKey(pub))
}

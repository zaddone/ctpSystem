package config
import(
	//"net/http"
	"github.com/BurntSushi/toml"
	"flag"
	"os"
	"os/exec"
	"context"
	"fmt"
	"log"
	"time"
	"bufio"
	"net"
	//"io"
	//"path/filepath"
)
var(
	LogFileName   = flag.String("c", "conf.log", "config log")
	Conf *Config
	forenoon = [2][2]int{[2]int{8,55},[2]int{11,30}}
	Afternoon = [2][2]int{[2]int{12,55},[2]int{15,15}}
	night = [2][2]int{[2]int{20,55},[2]int{2,30}}

	RunTime = [][2]int{
		[2]int{8,55},[2]int{11,30},
		[2]int{12,55},[2]int{15,15},
		[2]int{20,55},[2]int{2,30}}
	//StopTime = [][2]int{[2]int{11,30},[2]int{15,15},[2]int{2,30}}
)
type runStart struct {
	b time.Time
	e time.Time
	o bool
}

func WaitRunTime() {

	n := time.Now()
	//h := n.Hour()
	var runS []*runStart
	I := 0
	for i:=0;i<len(RunTime);i++{
		r := &runStart{
			b:time.Date(n.Year(),n.Month(),n.Day(),RunTime[i][0],RunTime[i][1],0,0,n.Location()),
		}
		I = len(runS)
		if I==0{
			r.o = true
		}else{
			runS[I-1].e = r.b
			if I%2 == 0 {
				r.o = true
			}
		}
		runS = append(runS,r)
	}
	runS[I].e = runS[0].b.Add(time.Hour*24)
	//nu := n.Unix()
	//var r_ *runStart
	for _,r := range runS {
		if r.b.Before(n) && r.e.After(n){
		if !r.o{
			<-time.After(time.Duration(n.Unix() - r.e.Unix()))
		}
			break
		}
	}
}
func init(){
	//EntryList = make(chan *Entry,1000)
	//flag.Parse()
	Conf = NewConfig(*LogFileName)
}

type UserInfo struct {
	BrokerID string
	UserID string
	Password string
	PasswordBak string
	Taddr []string
	Maddr []string
	DefAdd int
	sendTr chan []byte
	sendMd chan []byte
}

func (self *UserInfo)RunTr(path string,local string,hand func([]byte)){
	for{
		self.sendTr = make(chan []byte,100)
		runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Taddr[self.DefAdd],
			local,
		},self.sendTr,hand)
		WaitRunTime()
	}
}
func (self *UserInfo)RunMd(path string,local string,hand func([]byte)){

	for{
		self.sendMd = make(chan []byte,100)
		runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Maddr[self.DefAdd],
			local,
		},self.sendMd,hand)
		WaitRunTime()
	}
}

func runComm(
	path string,
	word []string,
	send chan []byte,
	hand func([]byte)){
	//fmt.Println(path,word)
	ctx,cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx,path,word...)

	out,err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(out)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	//var isPrefix bool
	ch := make(chan bool,1)
	go func(addr string){
		var line []byte
		//line,isPrefix,err = buf.ReadLine()
		line,err = buf.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
			cancel()
			ch<-true
			return
		}
		le := len(line)
		if le==0{
			ch<-true
			return
		}
		line = line[:len(line)-1]
		fmt.Println("conn",string(line))
		//hand(line)
		ch <- true
		go func(){
			rAddr, err := net.ResolveUnixAddr("unixgram",addr)
			if err != nil {
				panic(err)
				log.Fatal(err)
				return
			}
			c,err := net.DialUnix("unixgram",nil,rAddr)
			if err != nil {
				panic(err)
				log.Fatal(err)
			}
			sum:=0
			var n int
			for db:= range send {
				//fmt.Println(addr,sum,"<---------------",string(db))
				n,err = c.Write(db)
				if err != nil {
					//panic(err)
					log.Println(err)
					//break
				}
				sum+=n
				//fmt.Println(addr,sum)
			}
			fmt.Println("close--------------->",addr)
			c.Close()
		}()
		for{
			line,err = buf.ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
				cancel()
				break
			}
			line = line[:len(line)-1]
			//log.Println(addr,"------>",string(line))
			hand(line)
		}
	}(word[5])
	//defer log.Println(path,"end")

	select{
	case <-ch:
		//break
		//if t == 1{
		//	return
		//}
	case <-time.After(time.Second*30):
		log.Println("time out")
		cancel()
		//return
	}
	cmd.Wait()
	close(send)
}
func (self *UserInfo) SendMd(db []byte) {
	//fmt.Println("md",len(self.sendMd),string(db))
	if self.sendMd != nil {
		self.sendMd <- db
	}
	//select{
	//case self.sendMd <- db:
	//default:
	//	fmt.Println("md lose",string(db))
	//}
}
func (self *UserInfo) SendTr(db []byte) {
	//log.Println(string(db))
	if self.sendTr != nil {
		self.sendTr <- db
	}
	//select{
	//case self.sendTr <- db:
	//default:
	//	fmt.Println("tr lose",string(db))
	//}
}

type Config struct {
	User map[string]*UserInfo
	DefaultUser string
	Port string
	Static string
	Templates string
	DbPath string
	TPath string
	//RunAll bool
	MdServer string
	TrServer string
	IsTrader bool
	Weight int
	Example []string
}

func (self *Config) GetTPath() string {
	_,err := os.Stat(self.TPath)
	if err != nil {
		os.MkdirAll(self.TPath,0777)
	}
	return self.TPath
}
func (self *Config) GetDbPath() string {
	_,err := os.Stat(self.DbPath)
	if err != nil {
		os.MkdirAll(self.DbPath,0777)
	}
	return self.DbPath
}

func (self *Config)DefUser() (u *UserInfo) {
	u = self.User[self.DefaultUser]
	if u == nil || !self.IsTrader {
		return &UserInfo{}
	}
	return
}

func (self *Config)Save(fileName string){
	fi,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,0777)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	e := toml.NewEncoder(fi)
	err = e.Encode(self)
	if err != nil {
		panic(err)
	}
}
func NewConfig(fileName string)  *Config {
	var c Config
	_,err := os.Stat(fileName)
	if err != nil {
		c.Port=":8080"
		c.Static = "static"
		c.Templates = "templates"
		c.DbPath = "ctpdb"
		c.TPath = "ctpTrader"

		//c.RunAll = false
		c.DefaultUser = "150797"
		c.User = map[string]*UserInfo{
			"150797":&UserInfo{
				BrokerID : "9999",
				UserID : "150797",
				Password : "Dimon2019",
				PasswordBak : "abc2019",
				DefAdd : 2,
				Taddr : []string{
					"tcp://180.168.146.187:10100",
					"tcp://180.168.146.187:10101",
					"tcp://218.202.237.33:10102",
					"tcp://180.168.146.187:10130",
				},
				Maddr : []string{
					"tcp://180.168.146.187:10110",
					"tcp://180.168.146.187:10111",
					"tcp://218.202.237.33:10112",
					"tcp://180.168.146.187:10131",
				},
			},
			"153985":&UserInfo{
				BrokerID : "9999",
				UserID : "153985",
				Password : "Dimon2019",
				PasswordBak : "abc2019",
				DefAdd : 2,
				Taddr : []string{
					"tcp://180.168.146.187:10100",
					"tcp://180.168.146.187:10101",
					"tcp://218.202.237.33:10102",
					"tcp://180.168.146.187:10130",
				},
				Maddr : []string{
					"tcp://180.168.146.187:10110",
					"tcp://180.168.146.187:10111",
					"tcp://218.202.237.33:10112",
					"tcp://180.168.146.187:10131",
				},
			},
		}
		c.MdServer = "mdServer/mdServer"
		c.TrServer = "traderServer/traderServer"
		c.Example = []string{
			"AP001",
			"AP003",
			"AP005",
			"AP007",
			"AP010",
			"AP911",
			"AP912",
		}

		c.Weight = 4
		c.IsTrader = true


		c.Save(fileName)
	}else{
		if _,err := toml.DecodeFile(fileName,&c);err != nil {
			panic(err)
		}
	}
	return &c
}


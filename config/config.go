package config
import(
	//"net/http"
	"github.com/BurntSushi/toml"
	"flag"
	"os"
	"os/exec"
	"context"
	"bufio"
	"fmt"
	"log"
	"time"
	//"path/filepath"
)
var(
	LogFileName   = flag.String("c", "conf.log", "config log")
	Conf *Config
)
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

}
func (self *UserInfo)RunTr(path string,local string,hand func(string,[]byte)){
	for{
	runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Taddr[self.DefAdd],
			local,
		},hand)
	}
}
func (self *UserInfo)RunMd(path string,local string,hand func(string,[]byte)){
	for{
	runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Maddr[self.DefAdd],
			local,
		},hand)
	}
}
func runComm(path string,word []string,hand func(string,[]byte)){

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
	var line,db []byte
	var isPrefix bool
	ch := make(chan int,1)
	go func(){
		line,isPrefix,err = buf.ReadLine()

		if err != nil {
			fmt.Println(err)
			cancel()
			ch<-1
			return
		}
		hand(word[5],line)
		ch <- 2
	}()
	select{
	case t:=<-ch:
		if t == 1{
			return
		}
	case <-time.After(time.Second*30):
		cancel()
		return
	}
	//fmt.Println("run")
	for{
		line,isPrefix,err = buf.ReadLine()
		if err != nil {
			fmt.Println(err)
			cancel()
			break
		}
		if isPrefix {
			db = append(db,line...)
		}else{
			db = line
		}
		hand(word[5],db)
	}
	cmd.Wait()
	log.Println(path,"end")
}
type Config struct {
	User map[string]*UserInfo
	DefaultUser string
	Port string
	Static string
	Templates string
	//RunAll bool

	MdServer string
	TrServer string


	Weight int

	//BrokerID string
	//UserID string
	//Password string
	//PasswordBak string
	//Taddr []string
	//Maddr []string
	//DefAdd int

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

		//c.RunAll = false
		c.DefaultUser = "simnow"
		c.User = map[string]*UserInfo{
			"simnow":&UserInfo{
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
		}
		c.MdServer = "mdServer"
		c.TrServer = "TraderServer"

		c.Weight = 4


		c.Save(fileName)
	}else{
		if _,err := toml.DecodeFile(fileName,&c);err != nil {
			panic(err)
		}
	}
	return &c
}

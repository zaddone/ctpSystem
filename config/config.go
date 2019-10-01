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
	//"bufio"
	"net"
	//"io"
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

//func runCommbak(path string,word []string,hand func(string,[]byte)){
//
//	//fmt.Println(path,word)
//	ctx,cancel := context.WithCancel(context.Background())
//	cmd := exec.CommandContext(ctx,path,word...)
//	//cmd.Stdout = io.SeekStart
//	out,err := cmd.StdoutPipe()
//	if err != nil {
//		panic(err)
//	}
//	buf := bufio.NewReader(out)
//	err = cmd.Start()
//	if err != nil {
//		panic(err)
//	}
//	var line [8196]byte
//	//var isPrefix bool
//	ch := make(chan int,1)
//	addr := word[5]
//	go func(){
//		//line,isPrefix,err = buf.ReadLine()
//		line,err = buf.ReadBytes('\n')
//		line = line[:len(line)-1]
//		if err != nil {
//			fmt.Println(err)
//			cancel()
//			ch<-1
//			return
//		}
//		//log.Println(string(line))
//		//hand(addr,line)
//		ch <- 2
//	}()
//	select{
//	case t:=<-ch:
//		if t == 1{
//			return
//		}
//	case <-time.After(time.Second*30):
//		log.Println("time out")
//		cancel()
//		return
//	}
//
//	//fmt.Println("run")
//	for{
//		//line,isPrefix,err = buf.ReadLine()
//		line,err := buf.ReadBytes('\n')
//		if err != nil {
//			fmt.Println(err)
//			cancel()
//			break
//		}
//		line = line[:len(line)-1]
//		//fmt.Println(string(line))
//		//if isPrefix {
//		//	db = append(db,line...)
//		//}else{
//		//	db = line
//		//}
//		//n,err := out.Read(line[:])
//		//if err != nil {
//		//	fmt.Println(err)
//		//	cancel()
//		//	break
//		//}
//		log.Println("log",string(line))
//		go hand(addr,line)
//	}
//	cmd.Wait()
//	log.Println(path,"end")
//}

func runComm(path string,word []string,hand func(string,[]byte)){
	ctx,cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx,path,word...)
	out,err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	//buf := bufio.NewReader(out)
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	//var isPrefix bool
	ch := make(chan int,1)
	addr := word[5]
	go func(){
		var line [4096]byte
		var err error
		var n int
		n,err = out.Read(line[:])
		//line,isPrefix,err = buf.ReadLine()
		if err != nil {
			log.Println(err)
			cancel()
			ch<-1
			return
		}
		log.Println(addr,string(line[:n]))
		//hand(addr,line)
		ch <- 2
		for{
			n,err = out.Read(line[:])
			//line,err = buf.ReadBytes('\n')
			//line = line[:len(line)-1]
			if err != nil {
				log.Println(addr,err)
				cancel()
				return
			}
			log.Println(addr,string(line[:n]))
		}
	}()
	select{
	case t:=<-ch:
		if t == 1{
			return
		}
	case <-time.After(time.Second*30):
		log.Println("time out")
		cancel()
		return
	}
	//ctx.Done()
	local := addr + "_"
	err = os.Remove(local)
	if err != nil {
		fmt.Println(err)
	}
	lAddr, err := net.ResolveUnixAddr("unixgram",local)
	if err != nil {
		fmt.Println(err)
		return
	}
	ln, err := net.ListenUnixgram("unixgram", lAddr )
	if err!= nil {
		fmt.Println(err)
		return
	}
	go func (){
		for {
			var buf [4096]byte
			n,err := ln.Read(buf[:])
			if err != nil {
				log.Println(addr,err)
				return
				//panic(err)
			}
			hand(addr,buf[:n])
		}
	}()
	cmd.Wait()
	ln.Close()

}
//func UnixServer(local,addr string,hand func(string,[]byte)){
//	err := os.Remove(local)
//	if err != nil {
//		fmt.Println(err)
//	}
//	lAddr, err := net.ResolveUnixAddr("unixgram",local)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	ln, err := net.ListenUnixgram("unixgram", lAddr )
//	if err!= nil {
//		fmt.Println(err)
//		return
//	}
//	for{
//		var buf [1024]byte
//		n,_,err := ln.ReadFromUnix(buf[:])
//		if err != nil {
//			panic(err)
//		}
//		go hand(addr,buf[:n])
//	}
//	ln.Close()
//}
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

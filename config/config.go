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
	sendTr chan []byte
	sendMd chan []byte
}
func (self *UserInfo)RunTr(path string,local string,hand func([]byte)){
	self.sendTr = make(chan []byte,100)
	for{
	runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Taddr[self.DefAdd],
			local,
		},self.sendTr,hand)
	}
}
func (self *UserInfo)RunMd(path string,local string,hand func([]byte)){
	self.sendMd = make(chan []byte,100)
	for{
	runComm(path,[]string{
			self.BrokerID,
			self.UserID,
			self.Password,
			self.PasswordBak,
			self.Maddr[self.DefAdd],
			local,
		},self.sendMd,hand)
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
	//cmd.Stdout = io.SeekStart
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
		line = line[:len(line)-1]
		if err != nil {
			fmt.Println(err)
			cancel()
			ch<-true
			return
		}
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
			//log.Println(addr,"<------",string(line))
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
		return
	}
	cmd.Wait()
	close(send)
}
func (self *UserInfo) SendMd(db []byte) {
	//fmt.Println("md",len(self.sendMd),string(db))
	self.sendMd <- db
	//select{
	//case self.sendMd <- db:
	//default:
	//	fmt.Println("md lose",string(db))
	//}
}
func (self *UserInfo) SendTr(db []byte) {
	//fmt.Println("tr",len(self.sendTr),string(db))
	self.sendTr <- db
	//select{
	//case self.sendTr <- db:
	//default:
	//	fmt.Println("tr lose",string(db))
	//}
}

//func UnixSend(raddr string,db chan string) error {
//	//fmt.Println("send",raddr)
//	rAddr, err := net.ResolveUnixAddr("unixgram",raddr)
//	if err != nil {
//		return err
//	}
//	c,err := net.DialUnix("unixgram",nil,rAddr)
//	if err != nil {
//		return err
//	}
//	for{
//		_,err = c.Write([]byte(<-db))
//		if err != nil {
//			log.Println(err)
//			break
//		}
//	}
//	c.Close()
//	return err
//}
//func runComm_(path string,word []string,hand func(string,[]byte)){
//	ctx,cancel := context.WithCancel(context.Background())
//	cmd := exec.CommandContext(ctx,path,word...)
//	out,err := cmd.StdoutPipe()
//	if err != nil {
//		panic(err)
//	}
//	//buf := bufio.NewReader(out)
//	err = cmd.Start()
//	if err != nil {
//		panic(err)
//	}
//	//var isPrefix bool
//	ch := make(chan int,1)
//	addr := word[5]
//	go func(){
//		var line [4096]byte
//		var err error
//		var n int
//		n,err = out.Read(line[:])
//		//line,isPrefix,err = buf.ReadLine()
//		if err != nil {
//			log.Println(err)
//			cancel()
//			ch<-1
//			return
//		}
//		log.Println(addr,string(line[:n]))
//		//hand(addr,line)
//		ch <- 2
//		for{
//			n,err = out.Read(line[:])
//			//line,err = buf.ReadBytes('\n')
//			//line = line[:len(line)-1]
//			if err != nil {
//				log.Println(addr,err)
//				cancel()
//				return
//			}
//			log.Println(addr,string(line[:n]))
//		}
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
//	//ctx.Done()
//	local := addr + "_"
//	err = os.Remove(local)
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
//	go func (){
//		for {
//			var buf [4096]byte
//			n,err := ln.Read(buf[:])
//			if err != nil {
//				log.Println(addr,err)
//				return
//				//panic(err)
//			}
//			go hand(addr,buf[:n])
//		}
//	}()
//	cmd.Wait()
//	ln.Close()
//
//}
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

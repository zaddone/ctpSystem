package main
import(
	"fmt"
	"os"
	"log"
	"net"
	"flag"
	"time"
	"strings"
	"github.com/boltdb/bolt"
	"github.com/zaddone/ctpSystem/cache"
	"github.com/zaddone/ctpSystem/config"
)
var (
	traderAddr = flag.String("traddr", "/tmp/trader", "trader addr")
	mdAddr = flag.String("mdaddr", "/tmp/market", "trader addr")
	dbName = flag.String("db","Ins.db","db name")
	DB *bolt.DB
	//Farmat = "20060102T15:04:05"
	Cache = cache.NewCache()
)
func init(){
	flag.Parse()
	var err error
	DB,err = bolt.Open(*dbName,0600,nil)
	if err != nil {
		panic(err)
	}
}

func main(){
	fmt.Println("start")

	taddr := getAddr(config.Conf.Taddr)
	maddr := getAddr(config.Conf.Maddr)
	fmt.Println(taddr,maddr)
	go UnixServer(*traderAddr,RouterTrader)
	go UnixServer(*mdAddr,RouterMarket)
	select{}
}


func RouterTrader(db []byte){
	dbs := strings.Split(string(db)," ")
	switch(dbs[0]){
	case "ins":
		err := UnixSend(*mdAddr,db)
		if err != nil {
			panic(err)
		}


	case "addr":
		addr := getAddr(config.Conf.Taddr)
		if len(addr) ==0 {
			return
		}
		err := UnixSend(*traderAddr,[]byte(addr))
		if err != nil {
			panic(err)
		}
	case "config":
		addr := getAddr(config.Conf.Taddr)
		if len(addr) ==0 {
			return
		}
		str := fmt.Sprintf(
			"config %s %s %s %s %s",
			config.Conf.BrokerID,
			config.Conf.UserID,
			config.Conf.Password,
			config.Conf.PasswordBak,
			addr,
		)
		fmt.Println(str)
		err := UnixSend(*traderAddr,[]byte(str))
		if err != nil {
			panic(err)
		}

	default:
		fmt.Println(string(db))
	}
}
func RouterMarket(db []byte){
	dbs := strings.Split(string(db)," ")
	//fmt.Println(dbs)
	switch(dbs[0]){
	case "market":
		c := &cache.Candle{}
		err := c.Load(dbs[1])
		if err != nil {
			log.Println(err)
			return
		}
		Cache.Add(c)
	case "addr":
		addr := getAddr(config.Conf.Maddr)
		if len(addr) ==0 {
			return
		}
		err := UnixSend(*mdAddr,[]byte(addr))
		if err != nil {
			panic(err)
		}
	case "config":
		addr := getAddr(config.Conf.Maddr)
		if len(addr) ==0 {
			return
		}
		str := fmt.Sprintf(
			"config %s %s %s %s %s",
			config.Conf.BrokerID,
			config.Conf.UserID,
			config.Conf.Password,
			config.Conf.PasswordBak,
			addr,
		)
		err := UnixSend(*mdAddr,[]byte(str))
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println(string(db))
	}
}
func UnixSend(raddr string,db []byte) error {
	rAddr, err := net.ResolveUnixAddr("unixgram",raddr +"_")
	if err != nil {
		return err
	}
	c,err := net.DialUnix("unixgram",nil,rAddr)
	if err != nil {
		return err
	}
	_,err = c.Write(db)
	return err
}
func UnixServer(local string,hand func([]byte)){
	err := os.Remove(local)
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
	for{
		var buf [1024]byte
		n,_,err := ln.ReadFromUnix(buf[:])
		if err != nil {
			panic(err)
		}
		go hand(buf[:n])
		//fmt.Println(string(buf[:n]),raddr)

	}
	ln.Close()
}

func checkTcp(addr string) (error,int64) {

	b := time.Now().UnixNano()
	ad := strings.Split(addr,"://")
	fmt.Println(ad)
	hawkServer,err := net.ResolveTCPAddr("tcp", ad[1])
	if err != nil {
		return err,0
	}
	c,err := net.DialTCP("tcp",nil,hawkServer)
	if err != nil {
		return err,0
	}
	defer c.Close()
	return nil,time.Now().UnixNano() - b

}
func getAddr(addrs []string) (addr string) {

	var min int64
	for _,a_ := range addrs {
		err,d := checkTcp(a_)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(a_,d)
		if min==0 {
			addr=a_
			min = d
		}
	}
	return

}

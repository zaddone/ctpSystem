package main
import(
	"fmt"
	"os"
	"net"
	"flag"
	"bytes"
	"github.com/boltdb/bolt"

)
//type handle func([]byte)error
var (
	traderAddr = flag.String("traddr", "/tmp/trader", "trader addr")
	mdAddr = flag.String("mdaddr", "/tmp/market", "trader addr")
	dbName = flag.String("db","Ins.db","db name")
	db *bolt.DB
	//MarketRouter = map[string]handle
)
func init(){
	flag.Parse()
	var err error
	db,err = bolt.Open(*dbName,0600,nil)
	if err != nil {
		panic(err)
	}
}

func main(){
	fmt.Println("start")
	go UnixServer(*traderAddr,RouterTrader)
	go UnixServer(*mdAddr,RouterMarket)
	select{}
}


func RouterTrader(db []byte){
	dbs := bytes.Split(db,[]byte{':'})
	switch(string(dbs[0])){
	case "ins":
		err := UnixSend(*mdAddr,db)
		if err != nil {
			panic(err)
		}
	default:
		fmt.Println(string(db))
	}
}
func RouterMarket(db []byte){
	dbs := bytes.Split(db,[]byte{':'})
	switch(string(dbs[0])){
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
	var buf [1024]byte
	for{
		n,_,err := ln.ReadFromUnix(buf[:])
		if err != nil {
			panic(err)
		}
		hand(buf[:n])
		//fmt.Println(string(buf[:n]),raddr)

	}
	ln.Close()
}


package main
import(
	"fmt"
	"os"
	"time"
	//"io"
	"log"
	"net"
	"flag"
	"bytes"
	"strings"
	"github.com/boltdb/bolt"
	"encoding/gob"
	"encoding/binary"
	"strconv"

)
//type handle func([]byte)error
var (
	traderAddr = flag.String("traddr", "/tmp/trader", "trader addr")
	mdAddr = flag.String("mdaddr", "/tmp/market", "trader addr")
	dbName = flag.String("db","Ins.db","db name")
	DB *bolt.DB
	Farmat = "20060102T15:04:05"
	//MarketRouter = map[string]handle
)
func init(){
	flag.Parse()
	var err error
	DB,err = bolt.Open(*dbName,0600,nil)
	if err != nil {
		panic(err)
	}
}

type Candle struct{
	ins string
	date int64
	Ask float64
	Bid float64
}
func (self *Candle) encode() ([]byte,error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(self)
	return buf.Bytes(),err
}

func (self *Candle)decode(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(self)
}
func (self *Candle) load(db string)(err error){
	db_ := strings.Split(db,",")
	//fmt.Println(db_)
	self.ins = db_[0]
	d,err := time.Parse(Farmat,db_[1])
	if err != nil {
		return err
	}
	self.date = d.Unix()
	if len(db_[2])>30 || len(db_[3])>30{
		fmt.Println(db)
		return fmt.Errorf("too long")
	}
	self.Ask,err = strconv.ParseFloat(db_[2],64)
	if err != nil {
		return err
	}
	self.Bid,err = strconv.ParseFloat(db_[2],64)
	if err != nil {
		return err
	}
	return nil
}



func (self *Candle) toSave(db *bolt.DB)error{
	return db.Batch(func(t *bolt.Tx)error{
		b,err := t.CreateBucketIfNotExists([]byte(self.ins))
		if err != nil {
			return err
		}
		k := make([]byte,8)
		binary.BigEndian.PutUint64(k,uint64(self.date))
		v,err := self.encode()
		if err != nil {
			return err
		}
		//fmt.Println(k)
		return b.Put(k,v)
	})

}


func main(){
	fmt.Println("start")
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
	default:
		fmt.Println(string(db))
	}
}
func RouterMarket(db []byte){
	dbs := strings.Split(string(db)," ")
	//fmt.Println(dbs)
	switch(dbs[0]){
	case "market":
		c := &Candle{}
		err := c.load(dbs[1])
		if err != nil {
			log.Println(err)
			return
		}
		err = c.toSave(DB)
		if err != nil {
			log.Println(err)
			return
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


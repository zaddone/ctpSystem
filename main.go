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
	"github.com/gin-gonic/gin"
	"github.com/axgle/mahonia"
	"net/http"
	"strconv"
	"os/exec"
	"bytes"

	//_ "github.com/zaddone/ctpSystem/route"
)

var (
	traderAddr = flag.String("traddr", "/tmp/trader", "trader addr")
	mdAddr = flag.String("mdaddr", "/tmp/market", "trader addr")
	dbName = flag.String("db","Ins.db","db name")
	DB *bolt.DB
	//Farmat = "20060102T15:04:05"
	Cache = cache.NewCache()
	DefaultAddr = config.Conf.DefAdd
)

func init(){
	flag.Parse()
	var err error
	DB,err = bolt.Open(*dbName,0600,nil)
	if err != nil {
		panic(err)
	}
	Router := gin.Default()
	Router.Static("/"+config.Conf.Static,"./"+config.Conf.Static)
	Router.LoadHTMLGlob(config.Conf.Templates+"/*")
	Router.GET("/",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",nil)
	})

	Router.GET("/defaultaddr/:id",func(c *gin.Context){
		id ,err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err})
			return
		}
		DefaultAddr = id
		c.JSON(http.StatusOK,gin.H{"msg":"Success","t":config.Conf.Taddr[DefaultAddr],"m":config.Conf.Maddr[DefaultAddr]})
		//c.HTML(http.StatusOK,"index.tmpl",nil)
	})
	Router.GET("/wrun",func(c *gin.Context){
		words := c.DefaultQuery("word","")
		err = UnixSend(*mdAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}
	})
	Router.GET("/mlink",func(c *gin.Context){
		words := "addr "+config.Conf.Maddr[DefaultAddr]
		err = UnixSend(*traderAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}

	})
	Router.GET("/tlink",func(c *gin.Context){
		words := "addr "+config.Conf.Taddr[DefaultAddr]
		err = UnixSend(*traderAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}

	})
	Router.GET("/trun",func(c *gin.Context){
		words := c.DefaultQuery("word","")
		err = UnixSend(*traderAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}
	})
	Router.GET("/close",func(c *gin.Context){
		words := fmt.Sprintf("close %s",c.DefaultQuery("word","AP001"))
		err = UnixSend(*traderAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}
	})
	Router.GET("/open",func(c *gin.Context){
		words := fmt.Sprintf("open %s %s",c.DefaultQuery("word","AP001"),c.DefaultQuery("dis","buy"))
		err = UnixSend(*traderAddr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}
	})
	go Router.Run(config.Conf.Port)
	fmt.Println("start")

	//taddr := getAddr(config.Conf.Taddr)
	//maddr := getAddr(config.Conf.Maddr)
	//fmt.Println(taddr,maddr)
	go UnixServer(*traderAddr,RouterTrader)
	go UnixServer(*mdAddr,RouterMarket)
	go cache.Send(func(k *cache.MsgKey){
		//return
		if k.T{
			UnixSend(*traderAddr,k.DB)
		}else{
			UnixSend(*mdAddr,k.DB)
		}
	})
}
func Runserver(){
	cmd := exec.Command("./ctpServer")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
	    log.Fatal(err)
	}
	log.Printf("ctpServer: %q\n", out.String())
}

func main(){

	if config.Conf.RunServer{
		go Runserver()
	}
	select{}
	//exec.Command()

}

func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}





func RouterTrader(db []byte){
	//fmt.Println(string(db))
	dbs := strings.Split(string(db)," ")
	fmt.Println("trader",dbs)
	switch(dbs[0]){
	case "ins":
		err := DB.Batch(func(t *bolt.Tx)error{
			_,err := t.CreateBucketIfNotExists([]byte(dbs[1]))
			return err
		})
		if err != nil {
			panic(err)
		}
		err = UnixSend(*mdAddr,db)
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
		fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
	}
}
func RouterMarket(db []byte){
	dbs := strings.Split(string(db)," ")
	//fmt.Println("market",dbs)
	switch(dbs[0]){
	case "ins":
		err := DB.View(func(t *bolt.Tx)error{
			return t.ForEach(
				func(name []byte,b *bolt.Bucket)error{
				return UnixSend(*mdAddr,append(append(db,' '),name...))
			})
		})
		if err != nil {
			panic(err)
		}
	case "market":
		c := &cache.Candle{}
		err := c.Load(dbs[1])
		if err != nil {
			//log.Println(err)
			return
		}
		c.ToSave(DB)
		//Cache.Add(c)
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
		fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
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
	c.Close()
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
		//fmt.Println(local,ConvertToString(string(buf[:n]),"gbk","utf-8"))

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

	return addrs[DefaultAddr]
	//var min int64
	//for _,a_ := range addrs {
	//	err,d := checkTcp(a_)
	//	if err != nil {
	//		fmt.Println(err)
	//		continue
	//	}
	//	fmt.Println(a_,d)
	//	if min==0 || min > d {
	//		addr=a_
	//		min = d
	//	}
	//}
	//return

}



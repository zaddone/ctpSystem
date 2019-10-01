package main
import(
	"fmt"
	"os"
	"log"
	"net"
	"flag"
	//"time"
	"strings"
	"github.com/boltdb/bolt"
	"github.com/zaddone/ctpSystem/cache"
	"github.com/zaddone/ctpSystem/config"
	"github.com/gin-gonic/gin"
	"github.com/axgle/mahonia"
	"net/http"
	//"strconv"
	//"os/exec"
	//"bytes"
	"path/filepath"
	//_ "github.com/zaddone/ctpSystem/route"
)

var (

	traderAddr string
	mdAddr string
	dbName = flag.String("db","Ins.db","db name")
	DB *bolt.DB
	//Farmat = "20060102T15:04:05"
	Cache = cache.NewCache()
	//DefaultAddr = config.Conf.DefAdd
)
func Runserver(){
	_,md := filepath.Split(config.Conf.MdServer)
	mdAddr = "/tmp/"+ md
	_,tr := filepath.Split(config.Conf.TrServer)
	traderAddr = "/tmp/"+ tr

	for k,v := range config.Conf.User{
		go v.RunMd(
			config.Conf.MdServer,
			mdAddr+"_"+k,
			RouterMarket)
		go v.RunTr(
			config.Conf.TrServer,
			traderAddr+"_"+k,
			RouterTrader)
	}
}


func init(){
	flag.Parse()
	//_,tr := filepath.Split(config.Conf.TrServer)
	//traderAddr = "/tmp/"+ tr
	//_,md := filepath.Split(config.Conf.MdServer)
	//mdAddr = "/tmp/"+ md
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

	Router.GET("/wrun",func(c *gin.Context){
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		err = UnixSend(mdAddr+"_"+c.DefaultQuery("user",config.Conf.DefaultUser),[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
		}
	})
	Router.GET("/trun",func(c *gin.Context){
		//words := c.DefaultQuery("word","")
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		addr := traderAddr+"_"+c.DefaultQuery("user",config.Conf.DefaultUser)
		err = UnixSend(addr,[]byte(words))
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
		}else{
			c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words,"addr":addr})
		}
	})

	go Router.Run(config.Conf.Port)
	fmt.Println("start")

	//taddr := getAddr(config.Conf.Taddr)
	//maddr := getAddr(config.Conf.Maddr)
	//fmt.Println(taddr,maddr)
	//go UnixServer(*traderAddr,RouterTrader)
	//go UnixServer(*mdAddr,RouterMarket)
	go cache.Send(func(k *cache.MsgKey){
		//return
		if k.T{
			UnixSend(traderAddr+config.Conf.DefaultUser,k.DB)
		}else{
			UnixSend(mdAddr+config.Conf.DefaultUser,k.DB)
		}
	})
}

func main(){

	//if config.Conf.RunServer{
	Runserver()
	//}
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

func RouterTrader(path string,db []byte){
	//fmt.Println(string(db))
	dbs := strings.Split(string(db)," ")
	//fmt.Println("trader",dbs)
	switch(dbs[0]){
	case "ins":
		err := DB.Batch(func(t *bolt.Tx)error{
			_,err := t.CreateBucketIfNotExists([]byte(dbs[1]))
			return err
		})
		if err != nil {
			panic(err)
		}
		err = UnixSend(mdAddr+"_"+strings.Split(path,"_")[1],db)
		if err != nil {
			log.Println(dbs[1],err)
		}
	default:
		fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
	}
}
func RouterMarket(path string,db []byte){
	dbs := strings.Split(string(db)," ")
	//fmt.Println("market",dbs)
	switch(dbs[0]){
	case "ins":
		err := DB.View(func(t *bolt.Tx)error{
			return t.ForEach(
				func(name []byte,b *bolt.Bucket)error{
				return UnixSend(path,append(append(db,' '),name...))
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
		Cache.Add(c)
	default:
		fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
	}
}
func UnixSend(raddr string,db []byte) error {
	//fmt.Println("send",raddr)
	rAddr, err := net.ResolveUnixAddr("unixgram",raddr)
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

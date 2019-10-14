package main
import(
	"fmt"
	//"os"
	//"log"
	//"net"
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
	"bytes"
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
	TraderChan  = make(chan []byte,100)
	MarketChan  = make(chan []byte,100)
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
	Router.GET("/wsend",func(c *gin.Context){
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		if len(words) == 0 {
			c.JSON(http.StatusOK,gin.H{"msg":"Word is nil","word":words})
			return
		}
		RouterMarket([]byte(words))
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
	})
	Router.GET("/tsend",func(c *gin.Context){
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		if len(words) == 0 {
			c.JSON(http.StatusOK,gin.H{"msg":"Word is nil","word":words})
			return
		}
		RouterTrader([]byte(words))
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
	})
	Router.GET("/wrun",func(c *gin.Context){
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		if len(words) == 0 {
			c.JSON(http.StatusOK,gin.H{"msg":"Word is nil","word":words})
			return
		}
		u := config.Conf.User[c.DefaultQuery("user",config.Conf.DefaultUser)]
		if u== nil {
			c.JSON(http.StatusOK,gin.H{"msg":"Fount not","word":words})
			return
		}
		u.SendMd([]byte(words))
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words})
	})
	Router.GET("/trun",func(c *gin.Context){
		words := strings.Join(strings.Split(c.DefaultQuery("word",""),"_")," ")
		if len(words) == 0 {
			c.JSON(http.StatusOK,gin.H{"msg":"Word is nil","word":words})
			return
		}
		u := config.Conf.User[c.DefaultQuery("user",config.Conf.DefaultUser)]
		if u== nil {
			c.JSON(http.StatusOK,gin.H{"msg":"Fount not","word":words})
			return
		}
		u.SendTr([]byte(words))
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words,"user":u})
	})
	go Router.Run(config.Conf.Port)
	fmt.Println("start")
	//taddr := getAddr(config.Conf.Taddr)
	//maddr := getAddr(config.Conf.Maddr)
	//fmt.Println(taddr,maddr)
	//go UnixServer(*traderAddr,RouterTrader)
	//go UnixServer(*mdAddr,RouterMarket)
	//Def := config.Conf.User[config.Conf.DefaultUser]
	//go cache.Send(func(k *cache.MsgKey){
	//	//return
	//	if k.T{
	//		Def.SendTr(k.DB)
	//		//UnixSend(traderAddr+config.Conf.DefaultUser,k.DB)
	//	}else{
	//		Def.SendMd(k.DB)
	//		//UnixSend(mdAddr+config.Conf.DefaultUser,k.DB)
	//	}
	//})
	go runRouterMarket()
	go runRouterTrader()
}

func main(){

	Runserver()
	select{}

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
	TraderChan<-db
}
func runRouterTrader(){
	for{
		db:=<-TraderChan

		dbs := bytes.SplitN(db,[]byte{' '},2)
		switch(string(dbs[0])){
		case "ins":
			err := DB.Batch(func(t *bolt.Tx)error{
				_,err := t.CreateBucketIfNotExists([]byte(dbs[1]))
				return err
			})
			if err != nil {
				panic(err)
			}
			//fmt.Println(dbs)
			insMap := make(map[string]string)
			for _,mb := range  bytes.Split(dbs[1],[]byte{','}){
				vs := strings.Split(string(mb),":")
				insMap[vs[0]] = ConvertToString(vs[1],"gbk","utf-8")
			}
			ins := insMap["InstrumentID"]
			//fmt.Println(insMap)
			cache.StoreInsInfo(ins,insMap)
			for _,v := range config.Conf.User{
				v.SendMd([]byte("ins "+ins))
			}
		default:
			sdb := string(db)
			if strings.HasPrefix(sdb,"msg:"){
				fmt.Println(string(db))
			}else{
				fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
			}
		}

	}
}
func RouterMarket(db []byte){
	//fmt.Println("------>",len(MarketChan),string(db))
	MarketChan<-db
}
func runRouterMarket(){
	for{
	db:=<-MarketChan
	dbs := strings.Split(string(db)," ")
	//fmt.Println("market",dbs)
	switch(dbs[0]){
	case "ins":
		err := DB.View(func(t *bolt.Tx)error{
			return t.ForEach(
				func(name []byte,b *bolt.Bucket)error{
				db_ := append(append(db,' '),name...)
				for _,v := range config.Conf.User{
					v.SendMd(db_)
				}
				return nil
				//return UnixSend(path,append(append(db,' '),name...))
			})
		})
		if err != nil {
			panic(err)
		}
	case "market":
		c := &cache.Candle{}
		err := c.Load(dbs[1])
		if err == nil {
			Cache.Add(c)
			c.ToSave(DB)
			//log.Println(err)
		}
	default:
		fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
	}
	}
}
//func UnixSend(raddr string,db []byte) error {
//	//fmt.Println("send",raddr)
//	rAddr, err := net.ResolveUnixAddr("unixgram",raddr)
//	if err != nil {
//		return err
//	}
//	c,err := net.DialUnix("unixgram",nil,rAddr)
//	if err != nil {
//		return err
//	}
//	_,err = c.Write(db)
//	c.Close()
//	return err
//}
//func UnixServer(local string,hand func([]byte)){
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
//
//		go hand(buf[:n])
//		//fmt.Println(local,ConvertToString(string(buf[:n]),"gbk","utf-8"))
//
//	}
//	ln.Close()
//}

package main
import(
	"fmt"
	"os"
	//"log"
	"sync"
	"flag"
	//"time"
	"strings"
	//"github.com/boltdb/bolt"
	"github.com/zaddone/ctpSystem/cache"
	"github.com/zaddone/ctpSystem/config"
	"github.com/gin-gonic/gin"
	"github.com/axgle/mahonia"
	"net/http"
	"strconv"
	//"os/exec"
	"bytes"
	"path/filepath"
	//_ "github.com/zaddone/ctpSystem/route"
)
var (
	InsTraderMap sync.Map
	OrderCount [2]float64
	traderAddr string
	mdAddr string
	dbName = flag.String("db","Ins.db","db name")
	//DB *bolt.DB
	OrderRef int = 1
	//Farmat = "20060102T15:04:05"
	//Cache = cache.NewCache()
	//DefaultAddr = config.Conf.DefAdd
	TraderChan  = make(chan []byte,100)
	MarketChan  = make(chan []byte,100)
	TraderRouterMap = map[string]func([]byte){
		"ins":traderIns,
		"trade":tradeBack,
		"orderCancel":orderCancelBack,
		"orderWait":orderWaitBack,
		"InvestorPositionDetail":traderPosition,
	}
	MarketRouterMap =map[string]func([]byte){
		"ins":marketIns,
		"market":marketInfo,
	}
)
func Runserver(){
	_,md := filepath.Split(config.Conf.MdServer)
	mdAddr = "/tmp/"+ md
	_,tr := filepath.Split(config.Conf.TrServer)
	traderAddr = "/tmp/"+ tr
	u := config.Conf.DefUser()
	go u.RunMd(
		config.Conf.MdServer,
		mdAddr+"_"+config.Conf.DefaultUser,
		RouterMarket)
	go u.RunTr(
		config.Conf.TrServer,
		traderAddr+"_"+config.Conf.DefaultUser,
		RouterTrader)

	//for k,v := range config.Conf.User{
	//	go v.RunMd(
	//		config.Conf.MdServer,
	//		mdAddr+"_"+k,
	//		RouterMarket)
	//	go v.RunTr(
	//		config.Conf.TrServer,
	//		traderAddr+"_"+k,
	//		RouterTrader)
	//}
}
func initHttpRouter(){
	Router := gin.Default()
	Router.Static("/"+config.Conf.Static,"./"+config.Conf.Static)
	Router.LoadHTMLGlob(config.Conf.Templates+"/*")
	Router.GET("/",func(c *gin.Context){
		c.HTML(http.StatusOK,"index.tmpl",nil)
	})
	Router.GET("/start",func(c *gin.Context){

		//u := config.Conf.DefUser()
		u := config.Conf.User[c.DefaultQuery("user",config.Conf.DefaultUser)]
		if u== nil {
			c.JSON(http.StatusOK,gin.H{"msg":"Fount not"})
			return
		}
		Start(u)
		c.JSON(http.StatusOK,gin.H{
			"msg":"start",
			"user":u,
		})
	})
	Router.GET("/show",func(c *gin.Context){
		words := c.DefaultQuery("word","")
		if words=="" {
			c.JSON(http.StatusOK,gin.H{"msg":"Success","list":cache.ShowAll()})
			return
		}
		ca := cache.Show(words)
		if ca==nil {
			c.JSON(http.StatusOK,gin.H{"msg":"Word is nil","word":words})
			return
		}
		list,err := ca.GetOrderList()
		if err != nil {
			c.JSON(http.StatusOK,gin.H{"msg":err,"word":words})
			return
		}
		c.JSON(http.StatusOK,gin.H{"msg":"Success","word":words,"list":list})

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
		words := c.DefaultQuery("word","")
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
		words := c.DefaultQuery("word","")
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
	fmt.Println("start")
	go Router.Run(config.Conf.Port)
}

func init(){
	flag.Parse()
	//var err error
	//DB,err = bolt.Open(*dbName,0600,nil)
	//if err != nil {
	//	panic(err)
	//}
	initHttpRouter()

	go runRouter(MarketChan,MarketRouterMap)
	go runRouter(TraderChan,TraderRouterMap)
}
func main(){
	Runserver()
	select{}
}
func Start(s *config.UserInfo){
	s.SendTr([]byte("ReqQrySettlementInfo"))
	s.SendTr([]byte("ReqSettlementInfoConfirm"))
	s.SendTr([]byte("Instrument"))
	s.SendTr([]byte("InvestorPositionDetail"))
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
func RouterMarket(db []byte){
	MarketChan<-db
}
func traderPosition(db []byte){
	//fmt.Println(string(db))
	dbs := bytes.Split(db,[]byte{','})
	if len(dbs)<6 {
		return
	}
	u := config.Conf.User[string(dbs[0])]
	if u== nil {
		return
	}

	can_ := cache.Show(string(dbs[3]))
	if can_ == nil {
		return
	}
	c_ := can_.GetLast()
	if c_== nil {
		return
	}
	can := c_.(*cache.Candle)
	t := "4"
	if bytes.Equal(dbs[1],dbs[2]){
		t = "3"
	}
	f := can.GetUpperLimitPrice()
	d := "0"
	if bytes.Equal(dbs[5],[]byte{'0'}) {
		d = "1"
		f = can.GetLowerLimitPrice()
	}
	send := fmt.Sprintf("OrderInsert,%s,%s,%d,%s,%s,%.2f",dbs[3],dbs[4],OrderRef,t,d,f)
	OrderRef++
	fmt.Println(send)
	u.SendTr([]byte(send))

}
func traderIns(db []byte){
	insMap := make(map[string]string)
	for _,mb := range  bytes.Split(db,[]byte{','}){
		vs := strings.Split(string(mb),":")
		//insMap[vs[0]] = ConvertToString(vs[1],"gbk","utf-8")
		insMap[vs[0]] = vs[1]
	}
	cache.StoreCache(insMap)
	config.Conf.DefUser().SendMd([]byte("ins,"+insMap["InstrumentID"]))
	//for _,v := range config.Conf.User{
	//	v.SendMd([]byte("ins "+insMap["InstrumentID"]))
	//}
}

func marketIns(db []byte){

	u := config.Conf.DefUser()
	err := filepath.Walk(
		config.Conf.GetDbPath(),
		func(path string,
		f os.FileInfo,
		err error)error{
			if f.IsDir() {
				return nil
			}
			db_ := []byte(string(db)+","+f.Name())
			u.SendMd(db_)
			//for _,v := range config.Conf.User{
			//	v.SendMd(db_)
			//}
			return nil
	})
	if err != nil {
		panic(err)
	}

}
func marketInfo(db []byte){
	c := &cache.Candle{}
	err := c.Load(string(db))
	if err == nil {
		cache.AddCandle(c)
		//go c.ToSave(DB)
	}
}

func runRouter(c chan []byte,rm map[string]func([]byte)){
	for{
		db:=<-c
		dbs := bytes.SplitN(db,[]byte{' '},2)
		h := rm[string(dbs[0])]
		if h == nil {
			fmt.Println(string(db))

			//sdb := string(db)
			//if strings.HasPrefix(sdb,"msg:"){
			//	fmt.Println(string(db))
			//}else{
			//	fmt.Println(ConvertToString(string(db),"gbk","utf-8"))
			//}
		}else{
			h(dbs[1])
		}
	}
}
func orderWaitBack(db []byte){
	dbs := strings.Split(string(db),"-")
	ca := cache.Show(dbs[0])
	if ca==nil {
		return
	}
	o := ca.GetOrder(dbs[1])
	if o != nil {
		o.Wait(dbs[2])
	}
	//ca.Order.Update(2,dbs[1],dbs[2])
}
func orderCancelBack(db []byte){

	dbs := strings.Split(string(db),"-")
	ca := cache.Show(dbs[0])
	if ca==nil {
		return
	}
	//ca.DelOrder(dbs[1])
	ca.Order = nil

	//ca.Order.Update(2,dbs[1],false)
}
func tradeBack(db []byte){
	dbs := strings.Split(string(db),"-")
	ca := cache.Show(dbs[0])
	if ca==nil {
		return
	}
	c,err := strconv.ParseFloat(dbs[1],64)
	if err != nil {
		//return
		panic(err)
	}
	if dbs[2]=="0"{
		o := ca.GetOrder(dbs[3])
		if o == nil {
			return
		}
		o.SetOpenPrice(c)
		//ca.Order.Update(2,dbs[3],c)
	}else{
		//ca.Order.Update(4,dbs[3],c)
		o := ca.GetOrder(dbs[3])
		if o == nil {
			return
		}
		o.EndOrder(ca,c)
		//ca.DelOrder(dbs[3])
	}

}
//func _tradeBack(db []byte){
//	dbs := strings.Split(string(db)," ")
//	fmt.Println(dbs)
//	or_,ok :=  cache.InsOrderMap.Load(dbs[0])
//	if !ok {
//		//fmt.Println(string(db))
//		return
//	}
//	or:=or_.(*cache.InsOrder)
//	if dbs[2]=="0"{
//
//		c,err := strconv.ParseFloat(dbs[1],64)
//		if err != nil {
//			//return
//			panic(err)
//		}
//		or.OpenPrice = c
//		fmt.Println(
//			"_open",
//			or.Open.Ask,
//			or.Open.Bid,
//			or.Open.Val(),
//			or.Dis)
//			//var r map[string]*cache.InsOrder
//			//r_,ok := InsTraderMap.Load(dbs[0])
//			//if !ok{
//			//	r = map[string]*cache.InsOrder{or}
//			//}else{
//			//	r = append(r_.([]*cache.InsOrder),or)
//			//}
//			//InsTraderMap.Store(dbs[0],r)
//	}else{
//		fmt.Println(
//			"__close",
//			or.OpenPrice,
//			dbs[1],
//			or.Dis)
//		c,err := strconv.ParseFloat(dbs[1],64)
//		if err != nil {
//			return
//			//panic(err)
//		}
//		o,err := strconv.ParseFloat(or.OpenPrice,64)
//		if err != nil {
//			return
//			//panic(err)
//		}
//		if (o<c) == or.Dis{
//			OrderCount[1]++
//		}else{
//			OrderCount[0]++
//		}
//		fmt.Println("------------------------------------")
//		fmt.Println(OrderCount,OrderCount[0]/OrderCount[1])
//		if or.Close != nil {
//			fmt.Println("-------",or.Open.Val(),or.Close.Val())
//		}
//	}
//
//}

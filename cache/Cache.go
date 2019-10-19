package cache
import(
	"fmt"
	"time"
	"sync"
	"github.com/zaddone/ctpSystem/config"
)
var (
	Count [5][3]float64
	//KeyChan = make(chan *MsgKey,100)
	InsInfoMap sync.Map
	InsOrderMap sync.Map
)
type InsOrder struct {
	Dis bool
	Open *Candle
	OpenPrice string
	Stop *Candle
	Close *Candle
	Orderdef string
	Closedef string
	Status bool
}

func (self *InsOrder)OpenOrder(){
	ins := self.Open.Name()
	ex,ok := InsInfoMap.Load(ins)
	if !ok {
		panic(fmt.Errorf("%s is nil",ins))
	}
	var dis string
	var price float64
	if self.Dis {
		dis = "0"
		//stop = self.Stop.Bid
		price = self.Open.Ask
	}else{
		dis = "1"
		//stop = self.Stop.Ask
		price = self.Open.Bid
	}
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("open %s %s %s %.5f 0",
		ins,
		ex.(map[string]string)["ExchangeID"],
		dis,
		price,
		//stop,
	)))
}

func (self *InsOrder)CloseOrder(c interface{}){
	self.Close = c.(*Candle)
	ins := self.Open.Name()
	ex,ok := InsInfoMap.Load(ins)
	if !ok {
		panic(fmt.Errorf("%s is nil",ins))
	}
	var dis string
	var f float64
	if self.Dis {
		dis = "1"
		f = self.Open.GetLowerLimitPrice()
	}else{
		dis = "0"
		f = self.Open.GetUpperLimitPrice()
	}
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("OrderInsert %s %s 3 %s %.5f",
		ins,
		ex.(map[string]string)["ExchangeID"],
		dis,
		f,
	)))
}
func NewInsOrder(_c,_c_ interface{},dis bool) (or *InsOrder) {
	c:= _c.(*Candle)
	c_:= _c_.(*Candle)
	//dis := c.Val()<c_.Val()
	//fmt.Println(c.Ask,c_.Ask,c.Ask-c_.Ask,c.Time()-c_.Time())
	//or_,ok := InsOrderMap.Load(c.Name())
	//if ok {
	//	or = or_.(*InsOrder)
	//	//fmt.Println(or.Open.Val(),c.Val())
	//	//if or_.(*InsOrder).Status {
	//		return nil
	//	//}
	//}
	or = &InsOrder{
		Open:c,
		Stop:c_,
		Dis:dis,
		Orderdef:fmt.Sprintf("%d",time.Now().Unix()),
	}
	//InsOrderMap.Load()
	InsOrderMap.Store(c.Name(),or)
	//fmt.Println(io.Open,io.Stop)
	or.OpenOrder()
	return
}
//type MsgKey struct{
//	T bool
//	DB []byte
//	Ins string
//}
//func Send(h func(*MsgKey)){
//	for{
//		//mk := <-KeyChan
//		//_m := MkMap.Load(mk.Ins)
//		h(<-KeyChan)
//		//MkMap.Store(mk.Ins,mk)
//	}
//}
func StoreInsInfo(k,v interface{}){
	InsInfoMap.Store(k,v)
}
type Cache struct {
	LayerMap sync.Map
	//BaseLayer map[string]*Layer
	//candleChan chan *Candle
}
func NewCache() *Cache {
	return &Cache{}
}
func (self *Cache) Add(c *Candle){
	self.add(c)
}
func (self *Cache) add(c *Candle){

	L,ok := self.LayerMap.Load(c.ins)
	if !ok{
		L = NewLayer()
		self.LayerMap.Store(c.ins,L)
	}
	L.(*Layer).Add(c)

}
func (self *Cache) Show(ins string) interface{} {

	L,ok := self.LayerMap.Load(ins)
	if !ok {
		return nil
	}
	return L.(*Layer).lastEl
}

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
	Close *Candle
	Orderdef string
	Status bool
}
func (self *InsOrder)OpenOrder(){
	ins := self.Open.Name()
	ex,ok := InsInfoMap.Load(ins)
	if !ok {
		panic(fmt.Errorf("%s is nil",ins))
	}
	var dis string
	var stop,price float64
	func() {
		if self.Dis {
			dis = "0"
			stop = self.Close.Bid
			price = self.Open.Ask
		}else{
			dis = "1"
			stop = self.Close.Ask
			price = self.Open.Bid
		}
	}()
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("open %s %s %s %s %.2f %.2f",
		ins,
		ex.(map[string]string)["ExchangeID"],
		self.Orderdef,
		dis,
		stop,
		price,
	)))
}
func UpdateOrder(_c interface{}){
	c:= _c.(*Candle)
	or_,ok := InsOrderMap.Load(c.Name())
	if !ok {
		return
	}
	or :=or_.(*InsOrder)
	if  or.Status{
		//InsOrderMap.Delete(c.Name())
		return
	}
	or.Open = c
	or.OpenOrder()
}
func CloseOrder(c Element){
	//or_,ok := InsOrderMap.Load(c.Name())
}
func NewInsOrder(_c,_c_ interface{}) (io *InsOrder) {
	c:= _c.(*Candle)
	c_:= _c_.(*Candle)
	dis := c.Val()>c_.Val()
	or_,ok := InsOrderMap.Load(c.Name())
	if ok && or_.(*InsOrder).Status {
		return nil
	}
	io = &InsOrder{
		Open:c,
		Close:c_,
		Dis:dis,
		Orderdef:fmt.Sprintf("%d",time.Now().Unix()),
	}

	InsOrderMap.Store(c.Name(),io)
	fmt.Println(io.Open,io.Close)
	io.OpenOrder()
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

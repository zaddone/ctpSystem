package cache
import(
	"fmt"
	//"time"
	"sync"
	"github.com/zaddone/ctpSystem/config"
	"github.com/boltdb/bolt"
	"path/filepath"
	//"os"
)
var (
	Count [5][3]float64
	OrderCount [4]float64
	//KeyChan = make(chan *MsgKey,100)
	//InsInfoMap sync.Map
	//InsOrderMap sync.Map
	CacheMap sync.Map
)
func AddCandle(c *Candle) {
	c_,ok := CacheMap.Load(c.Name())
	if !ok {
		panic(c.Name())
	}

	ca := c_.(*Cache)
	//go c.ToSave(ca.DB)
	ca.L.Add(c)
}
type InsOrder struct {

	InsInfo map[string]string

	Dis bool
	Open *Candle
	OpenPrice float64
	//OpenRef string

	Close *Candle
	ClosePrice float64
	//CloseRef string

	State int

}
func (self *InsOrder)Update(state int,v ...interface{}) {

	//return
	if (self.State+1) != state {
		fmt.Println(self.InsInfo["InstrumentID"],state,self.State)
		self.State = 0
		return
	}
	self.State = state
	switch state {
	case 1:
		self.OpenOrder(v[1].(*Candle),v[0].(bool))
	case 2:
		if !v[0].(bool){
			self.State = 0
		}else{
			self.OpenPrice = v[1].(float64)
		}
	case 3:
		self.CloseOrder(v[0].(*Candle))
	case 4:
		diff := self.Close.Val() - self.Open.Val()
		self.ClosePrice = v[0].(float64)
		diff_:= self.ClosePrice - self.OpenPrice
		ff := diff>0
		if (diff_>0) == ff{
			OrderCount[1]++

		}else{
			OrderCount[0]++
		}
		if self.Dis == ff {
			//OrderCount[2]+=math.Abs(diff)
			//OrderCount[3]+=math.Abs(diff_)
			OrderCount[3]++
		}else{
			OrderCount[2]++
			//OrderCount[2]+=math.Abs(diff)
			//OrderCount[3]+=math.Abs(diff_)
		}
		fmt.Println(OrderCount)
		self.State = 0
	}
	return

}

func (self *InsOrder)OpenOrder(open *Candle,_dir bool){
	//ins := self.Open.Name()
	self.Open = open
	self.Dis = _dir
	self.OpenPrice = 0
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
		fmt.Sprintf("OrderInsert %s %s %d 0 %s %.5f",
		self.Open.Name(),
		self.InsInfo["ExchangeID"],
		self.Open.Time(),
		dis,
		price,
		//stop,
	)))
}

func (self *InsOrder)CloseOrder(c *Candle){
	self.Close = c
	var dis string
	var f float64
	if self.Dis {
		dis = "1"
		//f = self.Open.GetUpperLimitPrice()
		f = self.Close.GetLowerLimitPrice()
		//f = self.Close.Bid
	}else{
		dis = "0"
		f = self.Close.GetUpperLimitPrice()
		//f = self.Open.GetLowerLimitPrice()
		//f = self.Close.Ask
	}
	config.Conf.DefUser().SendTr(
		[]byte(
			fmt.Sprintf("OrderInsert %s %s %d 3 %s %.5f",
			self.Open.Name(),
			self.InsInfo["ExchangeID"],
			self.Close.Time(),
			dis,
			f),
		),
	)
}

type Cache struct {
	L *Layer
	//Info map[string]string
	Order InsOrder
	DB *bolt.DB
}
func (self *Cache)GetLast() interface{} {
	return self.L.getLast()
}
//func (self *Cache)Open(_c interface{},dir bool){
//	c:= _c.(*Candle)
//}
func Show(ins string) *Cache {

	c_,ok := CacheMap.Load(ins)
	if !ok {
		return nil
	}
	return c_.(*Cache)

}
func StoreCache(info map[string]string){
	ins := info["InstrumentID"]
	_ , ok := CacheMap.Load(ins)
	if ok{
		return
	}
	c := &Cache{
		//Info:info,
		Order:InsOrder{InsInfo:info},
	}
	c.L=NewLayer(c)
	p := filepath.Join(
			config.Conf.GetDbPath(),
			ins,
		)
	//fmt.Println(p)
	var err error
	c.DB,err =  bolt.Open(
		p,
		0600,nil)
	if err != nil {
		panic(err)
	}
	CacheMap.Store(ins,c)
	return
}

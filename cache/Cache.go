package cache
import(
	"fmt"
	//"strings"
	//"time"
	"sync"
	"github.com/zaddone/ctpSystem/config"
	"github.com/boltdb/bolt"
	"path/filepath"
	//"os"
)
var (
	Count [5][4]float64
	OrderCount [6]float64
	//KeyChan = make(chan *MsgKey,100)
	//InsInfoMap sync.Map
	//InsOrderMap sync.Map
	//Order int = 1
	CacheMap sync.Map
)
func AddCandle(c *Candle) {
	c_,ok := CacheMap.Load(c.Name())
	if !ok {
		panic(c.Name())
	}

	ca := c_.(*Cache)
	if config.Conf.IsTrader{
		go c.ToSave(ca.DB)
	}
	if ca.L != nil {
		ca.L.Add(c)
	}
}
type InsOrder struct {

	InsInfo map[string]string
	DB *bolt.DB

	Dis bool
	Open *Candle
	OpenP float64
	OpenPrice float64
	OpenRef string

	Stop float64


	Close *Candle
	ClosePrice float64
	CloseRef string
	//Order int

	State int

}
func (self *InsOrder) SaveDB(p float64) error {
	return self.DB.Update(func(t *bolt.Tx)error{
		b,err := t.CreateBucketIfNotExists([]byte(self.InsInfo["InstrumentID"]))
		if err != nil {
			return err
		}
		k := []byte(fmt.Sprintf("%d",self.Open.Time()))
		v := b.Get(k)
		if v != nil {
			return fmt.Errorf("repeat %s",v)
		}
		return b.Put(k,[]byte(fmt.Sprintf("%.2f",p)))
	})
}

func (self *InsOrder) DeleteDB() error {
	return self.DB.Update(func(t *bolt.Tx)error{
		b  := t.Bucket([]byte(self.InsInfo["InstrumentID"]))
		if b == nil {
			return fmt.Errorf("bucker is nil")
		}
		return b.Delete([]byte(fmt.Sprintf("%d",self.Open.Time())))
	})
}

func (self *InsOrder) UpdateDB(p float64) error {
	return self.DB.Update(func(t *bolt.Tx)error{
		b  := t.Bucket([]byte(self.InsInfo["InstrumentID"]))
		if b == nil {
			return fmt.Errorf("bucker is nil")
		}
		k := []byte(fmt.Sprintf("%d",self.Open.Time()))
		v := b.Get(k)
		if v == nil {
			return fmt.Errorf("%s is Not Found",k)
		}
		return b.Put(k,append(v,[]byte(fmt.Sprintf(" %.2f",p))...))
	})
}

func (self *InsOrder)Update(state int,v ...interface{}) {

	switch state {
	case 1:
		if self.State!= 0 {
			return
		}
		self.State = state
		self.OpenOrder(v[1].(*Candle),v[0].(bool))
		self.SaveDB(self.OpenP)
	case 2:
		if (self.State!=1) || (self.State!=2) {
			return
		}
		self.State = state
		switch val := v[0].(type){
		case bool:
			self.DeleteDB()
			self.State = 0
		case float64:
			self.OpenPrice = val
			self.UpdateDB(self.OpenPrice)
		case string:
			self.OpenRef = val

		}
	case 3:
		if (self.State!=2) {
			return
		}
		if self.OpenPrice == 0 {
			self.ActionCancel()
			return
		}
		self.State = state
		self.CloseOrder(v[0].(*Candle))
		self.UpdateDB(self.Close.Val())
	case 4:
		diff := self.Close.Val() - self.Open.Val()
		self.ClosePrice = v[0].(float64)
		self.UpdateDB(self.ClosePrice)

		diff_:= self.ClosePrice - self.OpenPrice
		ff := diff>0
		ff_ := diff_>0
		if ff_ == ff{
			OrderCount[1]++

		}else{
			OrderCount[0]++
		}
		if self.Dis == ff {
			OrderCount[3]++
		}else{
			OrderCount[2]++
		}
		if self.Dis == ff_ {
			OrderCount[5]++
		}else{
			OrderCount[4]++
		}
		fmt.Println(OrderCount)
		self.State = 0
	}
	return

}

func (self *InsOrder)ActionCancel(){
	if len(self.OpenRef)==0 {
		return
	}
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("OrderAction %s %s %s",
		self.Open.Name(),
		self.InsInfo["ExchangeID"],
		self.OpenRef,
	)))
}
func (self *InsOrder)OpenOrder(open *Candle,_dir bool){
	//ins := self.Open.Name()
	self.Open = open
	self.Dis = _dir
	self.OpenPrice = 0
	var dis string
	//var stop float64
	if self.Dis {
		dis = "0"
		self.Stop  = self.Open.Ask
		self.OpenP = self.Open.Bid
	}else{
		dis = "1"
		self.Stop  = self.Open.Bid
		self.OpenP = self.Open.Ask
	}
	//Order++
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("OrderInsert %s %s 0 %s %.5f %.5f",
		self.Open.Name(),
		self.InsInfo["ExchangeID"],
		dis,
		self.OpenP,
		self.Stop,
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
	//Order++
	config.Conf.DefUser().SendTr(
		[]byte(
			fmt.Sprintf("OrderInsert %s %s 3 %s %.5f 0",
			self.Open.Name(),
			self.InsInfo["ExchangeID"],
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
	//IsAdd bool
	//DBT *bolt.DB
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

	DB,err :=  bolt.Open(
		filepath.Join(
			config.Conf.GetDbPath(),
			ins,
		),
		0600,nil)
	if err != nil {
		panic(err)
	}
	c := &Cache{
		DB:DB,
		//Info:info,
		//Order:InsOrder{InsInfo:info},
	}
	CacheMap.Store(ins,c)

	if config.Conf.IsTrader{
		isAdd := false
		for _,e := range config.Conf.Example{
			isAdd =  ins == e
			if isAdd {
				break
			}
		}
		if !isAdd {
			return
		}
	}
	c.L=NewLayer(c)
	//fmt.Println(p)
	//var err error

	c.Order=InsOrder{InsInfo:info}
	c.Order.DB,err = bolt.Open(
		filepath.Join(
			config.Conf.GetTPath(),
			ins,
		),
		0600,nil)
	if err != nil {
		panic(err)
	}
	return
}

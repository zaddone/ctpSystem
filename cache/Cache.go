package cache
import(
	"fmt"
	"log"
	//"strings"
	//"time"
	"sync"
	"github.com/zaddone/ctpSystem/config"
	"github.com/boltdb/bolt"
	"path/filepath"
	"encoding/gob"
	"bytes"
	//"os"
)
var (
	Count [5][4]float64
	OrderCount [6]float64

	Order int = 1
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

	insInfo map[string]string
	db *bolt.DB

	Dis bool
	Open *Candle
	OpenP float64
	OpenPrice float64
	OpenRef string
	OpenSys string

	Stop float64

	Close *Candle
	ClosePrice float64
	CloseRef string
	//Order int

	State int

	//par *InsOrder
	//children *InsOrder
}

func (self *InsOrder) LoadByte(data []byte) error {

	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(self)

}

func (self *InsOrder) ToByte() []byte {

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(self)
	if err != nil {
		log.Fatal("encode:", err)
	}
	return network.Bytes()
}

func (self *InsOrder) SetOpenPrice(p float64){
	self.OpenPrice = p
}

func (self *InsOrder) Wait(OrderSys string){

	if self.OpenSys == "" {
		self.OpenSys = OrderSys
	}else if self.OpenSys != OrderSys {
		panic(fmt.Errorf("is not Same orderSys :%s %s",self.OpenSys,OrderSys))
	}

}

func (self *InsOrder) SaveDB() error {
	if self.db == nil {
		return fmt.Errorf("db is nil")
	}
	return self.db.Batch(func(t *bolt.Tx)error{
		b,err := t.CreateBucketIfNotExists([]byte(self.insInfo["InstrumentID"]+"_order"))
		if err != nil {
			return err
		}
		k := []byte(fmt.Sprintf("%d",self.Open.Time()))
		v := b.Get(k)
		if v != nil {
			return fmt.Errorf("repeat %s",v)
		}
		return b.Put(k,self.ToByte())
	})
}


func (self *InsOrder) SendCloseOrder(c *Candle,ca *Cache){
	self.Close = c
	if self.OpenPrice == 0 {
		self.ActionCancel()
	}else{
		self.CloseOrder(c)
		ca.DelOrder(self.OpenRef)
		ca.LoadOrder(self.CloseRef,self)
	}
}
func (self *InsOrder)EndOrder(p float64){

	diff := self.Close.Val() - self.Open.Val()
	self.ClosePrice = p
	diff_ := self.ClosePrice - self.OpenPrice
	ff := diff>0
	ff_ := diff_>0
	if ff_ == ff {
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
	err := self.SaveDB()
	if err != nil {
		panic(err)
	}
}

func (self *InsOrder)ActionCancel(){
	if len(self.OpenRef)==0 {
		return
	}
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("OrderAction,%s,%s,%s,%s",
		self.Open.Name(),
		self.insInfo["ExchangeID"],
		self.OpenRef,
		self.OpenSys,
	)))
}
func (self *InsOrder)OpenOrder(open *Candle,_dir bool){
	//ins := self.Open.Name()
	self.Open = open
	self.Dis = _dir
	self.OpenPrice = 0
	self.CloseRef = ""
	//self.OpenRef = ""
	self.OpenSys = ""
	var dis string
	//var stop float64
	if self.Dis {
		dis = "0"
		//self.Stop  = self.Open.Ask
		self.OpenP = self.Open.Bid

		//self.OpenP = self.Open.Ask
	}else{
		dis = "1"
		//self.Stop  = self.Open.Bid
		self.OpenP = self.Open.Ask

		//self.OpenP = self.Open.Bid
	}
	Order++
	self.OpenRef =fmt.Sprintf("%012d", Order);
	config.Conf.DefUser().SendTr([]byte(
		fmt.Sprintf("OrderInsert,%s,%s,%s,0,%s,%.5f,%.5f",
		self.Open.Name(),
		self.insInfo["ExchangeID"],
		self.OpenRef,
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
	Order++
	self.CloseRef = fmt.Sprintf("%012d", Order)
	config.Conf.DefUser().SendTr(
		[]byte(
			fmt.Sprintf("OrderInsert,%s,%s,%s,3,%s,%.5f,0",
			self.Open.Name(),
			self.insInfo["ExchangeID"],
			self.CloseRef,
			dis,
			f),
		),
	)
}

type Cache struct {
	L *Layer
	Info map[string]string
	Order *InsOrder
	Orders map[string]*InsOrder
	DB *bolt.DB
	sync.Mutex
	//IsAdd bool
	//DBT *bolt.DB
}
func (self *Cache) AddOrder(dis bool,stop Element){
	self.Order = &InsOrder{
		insInfo:self.Info,
		Stop:func()float64{
			if dis{
				return stop.Max()
			}else{
				return stop.Min()
			}
		}(),
	}
	self.Order.OpenOrder(self.GetLast().(*Candle),dis)
	self.Order.db = self.DB
	self.Lock()
	self.Orders[self.Order.OpenRef] = self.Order
	self.Unlock()
}
func (self *Cache)DelOrder(orderRef string){
	self.Lock()
	delete(self.Orders,orderRef)
	self.Unlock()
	log.Println("map orders len:",len(self.Orders))
}
func (self *Cache)GetOrder(orderRef string)(o *InsOrder) {
	self.Lock()
	o = self.Orders[orderRef]
	self.Unlock()
	if o == nil{
		fmt.Println(orderRef)
		return
		panic("map Order is nil")
	}
	return
}
func (self *Cache)LoadOrder(k string,o *InsOrder){
	self.Lock()
	self.Orders[k] = o
	self.Unlock()
}
func (self *Cache) EachOrder(h func(string,*InsOrder)bool){
	self.Lock()
	for k,v := range self.Orders{
		if !h(k,v){
			break
		}
	}
	self.Unlock()
}
func (self *Cache)GetLast() interface{} {
	return self.L.getLast()
}
func (self *Cache)GetOrderList() (ios []*InsOrder,err error) {
	err = self.DB.View(func(t *bolt.Tx)error{
		b := t.Bucket([]byte(self.Info["InstrumentID"]+"_order"))
		if b == nil {
			return fmt.Errorf("b == nil")
		}
		return b.ForEach(func(k,v []byte)error{
			ino := &InsOrder{}
			ios = append(ios,ino)
			return ino.LoadByte(v)
		})

	})
	return
}
//func (self *Cache)Open(_c interface{},dir bool){
//	c:= _c.(*Candle)
//}
func ShowAll() (ca []map[string]string) {
	CacheMap.Range(func(k,v interface{})bool{
		ca = append(ca,v.(*Cache).Info)
		return true
	})
	return
}
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
		fmt.Println("store",ins)
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
		//Order:InsOrder{InsInfo:info},
		Info:info,
		Orders:map[string]*InsOrder{},
	}
	CacheMap.Store(ins,c)

	if len(config.Conf.Example)>0 {
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
	//c.Order=InsOrder{InsInfo:info}
	return

}

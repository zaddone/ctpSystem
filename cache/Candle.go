package cache
import(
	"fmt"
	"strconv"
	"time"
	"encoding/binary"
	"encoding/gob"
	"bytes"
	"strings"
	"github.com/boltdb/bolt"
)
var (
	Farmat = "20060102T15:04:05"
)
type Candle struct{
	ins string
	date int64
	Ask float64
	Bid float64
	v float64
	d float64
	upperLimitPrice float64
	lowerLimitPrice float64
}
func NewCandle(ins string,date int64,db []byte) (c *Candle) {
	c = &Candle{
		ins:ins,
		date:date,
	}
	err := c.decode(db)
	if err != nil {
		panic(err)
	}
	return c
}
func (self *Candle) GetUpperLimitPrice() float64 {
	return self.upperLimitPrice
}
func (self *Candle) GetLowerLimitPrice() float64 {
	return self.lowerLimitPrice
}
func (self *Candle) Name() string{
	return self.ins
}
func (self *Candle) Each(fn func(e Element)error)error{
	return fn(self)
}
func (self *Candle) Time() int64 {
	return self.date
}
func (self *Candle) LastTime() int64 {
	return self.date
}
func (self *Candle) Dur() int64 {
	return 1
}
func (self *Candle) SetDiff(d float64) {
	self.d = d
}
func (self *Candle) Diff() float64 {
	//if self.d == 0 {
	//	return self.Ask - self.Bid
	//}else{
	return self.d
	//}
}
func (self *Candle) Val() float64 {
	//if self.d != 0 {
	//	if self.d>0{
	//		return self.Ask
	//	}else{
	//		return self.Bid
	//	}
	//}
	if self.v == 0 {
		self.v = (self.Ask + self.Bid) /2
	}
	return self.v
}
func (self *Candle) encode() ([]byte,error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(self)
	return buf.Bytes(),err
}
func (self *Candle)Decode(data []byte) error {
	return self.decode(data)
}
func (self *Candle)decode(data []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(data)).Decode(self)
}
func (self *Candle)Load(db string) error {
	return self.load(db)
}
func (self *Candle) load(db string)(err error){
	db_ := strings.Split(db,",")
	self.ins = db_[0]
	d,err := time.Parse(Farmat,db_[1])
	if err != nil {
		return err
	}
	self.date = d.Unix()
	if len(db_[2])>30 || len(db_[3])>30{
		//fmt.Println(db)
		return fmt.Errorf("too long")
	}

	//fmt.Println(db_)
	self.Ask,err = strconv.ParseFloat(db_[2],64)
	if err != nil {
		return err
	}
	self.Bid,err = strconv.ParseFloat(db_[3],64)
	if err != nil {
		return err
	}
	self.upperLimitPrice,err = strconv.ParseFloat(db_[4],64)
	if err != nil {
		return err
	}
	self.lowerLimitPrice,err = strconv.ParseFloat(db_[5],64)
	if err != nil {
		return err
	}
	//if self.Ask == self.Bid {
	//	return fmt.Errorf("ask bid is same")
	//}
	return nil
}

func (self *Candle) ToSave(db *bolt.DB)error{
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

package cache
import(
	//"fmt"
	"sync"
)
var (
	Count [3]float64
	KeyChan = make(chan *MsgKey,100)
	//MkMap sync.Map
)
type MsgKey struct{
	T bool
	DB []byte
	Ins string
}
func Send(h func(*MsgKey)){
	for{
		//mk := <-KeyChan
		//_m := MkMap.Load(mk.Ins)
		h(<-KeyChan)
		//MkMap.Store(mk.Ins,mk)
	}
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

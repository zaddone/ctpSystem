package cache
import(
	//"fmt"
	"sync"
)
var (
	Count [3]float64
)
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

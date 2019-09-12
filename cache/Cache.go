package cache
import(
	//"fmt"
	//"sync"
)


type Cache struct {
	BaseLayer map[string]*Layer
	candleChan chan *Candle
}

func NewCache() *Cache {

	c := &Cache{
		BaseLayer:make(map[string]*Layer),
		candleChan:make(chan *Candle,100),
	}
	go c.syncAdd()
	return c

}

func (self *Cache) syncAdd(){
	for{
		self.add(<-self.candleChan)
	}
}

func (self *Cache) Add(c *Candle){
	self.candleChan <- c
}

func (self *Cache) add(c *Candle){

	L := self.BaseLayer[c.ins]
	if L == nil {
		L = &Layer{}
		self.BaseLayer[c.ins] = L
	}
	go L.add(c)

}

package cache
import(
	"math"
	"fmt"
	//"sync"
)

type Layer struct{
	cans []Element
	direction float64
	par *Layer
	tag int
	child *Layer
	canChan chan Element
	lastEl Element
	//sync.Mutex
}

func NewLayer() (L *Layer) {
	L = &Layer{
		canChan:make(chan Element,100),
	}
	go L.runChan()
	return L
}

//func (self *Layer) Add(e Element){
//
//	defer func(){
//		self.lastEl = e
//		self.cans = append(self.cans,e)
//		//self.direction = e.Val() - self.cans[0]
//	}()
//	if self.lastEl == nil {
//		return
//	}
//	dl := e.LastTime() -  self.lastEl.LastTime()
//	if dl <0 || dl > 60 {
//		self.cans = nil
//		self.par = nil
//		return
//	}
//	var absMax,absDiff float64
//	var mid int
//	for i,c := range self.cans {
//		absDiff = math.Abs(e.Val() - c.Val())
//		if absMax < absDiff {
//			absMax = absDiff
//			mid = i
//		}
//		//sum += c.Val()
//	}
//	if mid == 0 {
//		return
//	}
//	if self.par == nil {
//		self.par = &Layer{}
//		self.par.tag = self.tag+1
//	}
//	self.par.add(NewNode(self.cans[:mid+1]))
//	self.cans = self.cans[mid:]
//	return
//
//}

func (self *Layer) add(e Element){
	if e== nil {
		self.par = nil
		self.cans = nil
		return
	}

	le := len(self.cans)
	self.cans = append(self.cans,e)
	//if le<3 {
	//	return
	//}
	var absMax,absDiff float64
	var mid int
	for i,c := range self.cans[:le] {
		absDiff = math.Abs(e.Val() - c.Val())
		if absMax < absDiff {
			absMax = absDiff
			mid = i
		}
		//sum += c.Val()
	}
	if mid == 0 {
		return
	}
	if self.par == nil {
		self.par = &Layer{}
		self.par.tag = self.tag+1
		self.par.child = self
	}else{
		fmt.Println(self.tag,len(self.cans[:mid+1]),len(self.par.cans))
		self.par.add(NewNode(self.cans[:mid+1]))
	}
	self.cans = self.cans[mid:]
	return
}


func (self *Layer) runChan(){
	for{
		//self.Lock()
		self.add(<-self.canChan)
		//self.Unlock()
	}

}
func (self *Layer) Add(e Element){
	if self.lastEl !=nil {
		dl := e.LastTime() -  self.lastEl.LastTime()
		if dl <0 || dl > 60 {
			fmt.Println(dl)
			self.canChan<-nil
		}
	}
	self.canChan <- e
	self.lastEl = e
}

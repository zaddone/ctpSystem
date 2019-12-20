package cache
import(
	"math"
	"fmt"
	//"io"
	//"github.com/zaddone/analog/fitting"
	"github.com/zaddone/ctpSystem/config"
	//"github.com/boltdb/bolt"
	//"encoding/binary"
	//"encoding/gob"
	//"bytes"
	//"time"
	//"sync"
)
type Temple struct{
	can Element
	stop Element
	long float64
	lcan Element
	XMin,XMax,YMin,YMax float64
	Wei []float64
	Stats int
	Dis bool
}

func (self *Temple) SetStats(){
	self.Stats++
}
type Layer struct{
	cans []Element
	direction float64
	sum float64
	par *Layer
	tag int
	child *Layer
	canChan chan Element
	lastEl Element
	ca *Cache
	tem  *Temple
	splitID int
	isF bool
	//sync.Mutex
}

func NewLayer(ca *Cache) (L *Layer) {
	L = &Layer{
		ca:ca,
		canChan:make(chan Element,100),
	}
	go L.runChan()
	return L
}

func (self *Layer) checkTem() (isok bool) {

	c_ := self.getLast()
	t  := 0
	var Diff float64
	if self.tem.Dis {
		Diff = c_.Min() - self.tem.can.Max()
	}else{
		Diff = c_.Max() - self.tem.can.Min()
	}
	if self.tem.Dis == (Diff>0){
		Count[t][4] += math.Abs(Diff)
		Count[t][3]++
	}else{
		Count[t][4] -= math.Abs(Diff)
		Count[t][2]++
	}
	dis_:= c_.Val() - self.tem.can.Val()

	if  (dis_>0) == self.tem.Dis {
		Count[t][1]++
		Count[t][5] += math.Abs(dis_)
	}else{
		Count[t][0]++
		Count[t][5] -= math.Abs(dis_)
	}
	fmt.Println(Count[t],c_.Time() - self.tem.can.Time())


	self.tem = nil
	if self.ca.Order != nil {
		self.ca.Order.SendCloseOrder(c_.(*Candle),self.ca)
	}
	return
}

func (self *Layer) getLast() Element {
	if self.child != nil {
		return self.child.getLast()
	}else{
		return self.cans[len(self.cans)-1]
	}
}
func (self *Layer) isTem() *Layer {
	if self.tem != nil {
		return self
	}
	if self.par== nil {
		return nil
	}
	return self.par.isTem()
}

func (self *Layer) getTemplate(dis bool){

	if self.ca.L == nil {
		return
	}
	L := self.ca.L.isTem()
	if L != nil{
		if L.tag < self.tag {
			self.tem = L.tem
			L.tem = nil
		}
		return
	}
	//self.getNormalization(dis)
	self.tem = &Temple{Dis:dis}
	self.tem.can  = self.getLast()
	self.tem.lcan = self.cans[0]
	if config.Conf.IsTrader{
		self.ca.AddOrder(self.tem.Dis,self.tem.stop)
	}

}

func (self *Layer) runChan(){
	for c := range self.canChan{
	//for{
		self.baseAdd(c)
		//self.add_(<-self.canChan)
	}

}
func (self *Layer) baseAdd(e Element){
	if e == nil {
		if len(self.cans)>1{
			if self.par == nil {
				self.setPar()
			}
			self.par.add_1(NewNode(self.cans))
			//self.par.baseAdd(NewNode(self.cans))
			self.cans = nil
		}
		//self.cans = nil
		//self.par = nil
		return
	}
	le := len(self.cans)
	self.cans = append(self.cans,e)
	if le == 0 {
		return
	}
	last := self.cans[le-1]
	e.SetDiff(e.Val() - last.Val())
	e.SetDur(e.LastTime() - last.Time())

	if e.Diff() == 0 {
		e.SetDiff(last.Diff())
		return
	}
	if last.Diff()==0 {
		return
	}
	if (last.Diff() >0) == (e.Diff()>0){
		return
	}
	if self.par == nil {
		self.setPar()
	}
	self.par.add_1(NewNode(self.cans[:le]))
	self.cans = []Element{e}
	e.SetDiff(0)
	//self.cans=nil

}
func (self *Layer) Add(e Element){

	if e.Val() == 0 {
		return
	}
	if e.Max()==e.Min() {
		panic(0)
		return
	}
	if self.lastEl !=nil {
		dl := e.LastTime() -  self.lastEl.LastTime()
		if (dl < 0)  || (dl>300) {
			self.canChan<-nil
		}
	}
	self.canChan <- e
	self.lastEl = e

}
func (self *Layer) setPar() {

	if self.tag>10{
		return
	}
	self.par = &Layer{
		ca:self.ca,
		child:self,
		tag:self.tag+1,
	}
	fmt.Println(self.ca.Info["InstrumentID"],self.par.tag)
}
func (self *Layer) add_1(c Element) {

	self.cans = append(self.cans,c)
	//if len(self.cans)<3{
	//	return
	//}

	//if self.tag == 1 {
	self.Check(c)
	//self.CheckEnd()
	//}
	n1 := NewNode(self.cans)
	if math.Abs(self.direction) <= math.Abs(n1.Diff()){
		self.direction = n1.Diff()
		self.splitID = len(self.cans)-1
		//if self.tag >1 {
		//	self.Check()
		//}
		return
	}
	if n1.Diff()>0{
		if c.Val() > n1.Val(){
			return
		}
	}else{
		if c.Val() < n1.Val(){
			return
		}
	}
	C:=self.cans[0]
	var diffAbs,maxAbs float64
	var I int
	for i,c_ := range self.cans[1:]{
		diffAbs = math.Abs(c_.Val() - C.Val())
		if diffAbs > maxAbs {
			I = i+1
			maxAbs = diffAbs
		}
	}

	if self.par == nil {
		self.setPar()
	}
	if self.Par != nil{
		n_0 := NewNode(self.cans[:I+1])
		self.par.add_1(n_0)
	}
	self.cans = self.cans[I:]
	self.direction = c.Val() - self.cans[0].Val()
	self.CheckEnd()
	//if self.tem != nil {
	//	if self.tem.Stats==1{
	//		self.checkTem()
	//	}else{
	//		self.tem.SetStats()
	//	}
	//	//return
	//}


	//fmt.Printf("%d %10.2f %5d %5d %5d %10.2f %10.2f\r\n",self.tag,self.par.direction,len(self.par.cans),len(self.cans),len(n_0.Eles),self.direction,n_0.Diff())

}
func (self *Layer) CheckEnd(){
	if self.tem == nil {
		return
	}

	if self.tem.Stats==0{
		self.tem.SetStats()
		return
	}

	if math.Abs(self.direction)>math.Abs(self.par.cans[len(self.par.cans)-1].Diff()) != self.tem.Dis {
		self.checkTem()
	}

}
func (self *Layer) Check_1(e Element)bool{
	no := NewNode(self.cans)
	diff := math.Abs(e.Val() - no.Val())
	var sum,dur float64
	for _,c := range self.cans {
		d := float64(c.Dur())
		sum += (c.Max() -  c.Min())*d
		dur += d
	}
	sum/=dur
	if diff>sum{
		//fmt.Println(diff,sum)
		return true
	}
	return false
}
func (self *Layer) Check(e Element){
	if self.tem != nil {
		if self.tem.Stats <1{
			return
		}
		c_:= self.getLast()
		var Diff float64
		if self.tem.Dis {
			Diff = c_.Min() - self.tem.can.Max()
		}else{
			Diff = c_.Max() - self.tem.can.Min()
		}
		if self.tem.Dis == (Diff>0){
			self.checkTem()
		}

		//if (self.getLast().Val() > self.tem.can.Val()) == self.tem.Dis{
		//	fmt.Println("_______")
		//	self.checkTem()
		//}
		return
		//self.tem.SetStats()
	}
	if self.par == nil {
		return
	}
	if self.par.direction==0{
		return
	}
	//if (self.direction>0) == (self.par.direction>0){
	//	return
	//}
	//c := self.cans[len(self.cans)-1]
	//if (c.Max()-c.Min())*2 > math.Abs(self.direction) {
	//	return
	//}
	if self.Check_1(e){
	self.getTemplate(self.direction<0)
	}
}

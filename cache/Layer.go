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
	splitID_old int
	isF bool
	//sync.Mutex
	parCov float64
	parSum float64
	last *Node
}
func (self *Layer) getSum() {

	self.parSum = 0
	self.parCov = 0
	if self.par == nil {
		return
	}
	if self.par.direction == 0 {
		return
	}
	//if self.parCov != 0 {
	//	return self.parCov
	//}
	dis := self.direction>0
	var dur,d,f float64
	var cs  []float64
	for _,c := range self.par.cans {
		if c.Diff()<0 == dis {
			continue
		}
		d = float64(c.Dur())
		dur += d
		f = math.Abs(c.Diff())
		self.parSum += f*d
		cs = append(cs,f)
	}
	self.parSum /= dur
	for _,c := range cs {
		self.parCov += math.Pow(self.parSum - c,2)
	}
	self.parCov =math.Sqrt(self.parCov/float64(len(cs)))
	//fmt.Println(self.tag,self.parCov,self.parSum,len(cs),self.direction,self.par.direction)

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
		if L.tem.Dis != dis {
			L.checkTem()
			return
		}
		if L.tag < self.tag {
			//lcan := NewNode(self.par.cans)
			//if L.tem.can.Val()> lcan.Val() != dis {
			//	return
			//}
			self.tem = L.tem
			//self.tem.lcan = lcan
			L.tem = nil
		}
		return
	}
	//self.getNormalization(dis)
	tem := &Temple{Dis:dis}
	tem.can  = self.getLast()
	//tem.lcan = NewNode(self.par.cans)
	////tem.lcan = self.par.cans[len(self.par.cans)-1]
	//if tem.can.Val()>tem.lcan.Val() != dis{
	//	return
	//}
	//fmt.Println(self.tem.can.Val(),self.tem.lcan.Val(),dis)
	//self.tem.lcan = self.par.cans[len(self.par.cans)-1].Each(
	self.tem = tem
	if config.Conf.IsTrader{
		self.ca.AddOrder(self.tem.Dis,self.tem.stop)
	}

}

func (self *Layer) runChan(){
	for c := range self.canChan{
	//for{
		//self.baseAdd(c)
		self.add_1(c)
	}
	self.par = nil
	self.cans = nil

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

	//L := self.ca.L.isTem()
	//if L != nil {
	//	L.CheckEnd(e)
	//}

	if le == 0 {
		return
	}
	last := self.cans[le-1]
	//e.SetDiff(e.Val() - last.Val())
	//e.SetDur(e.LastTime() - last.Time())

	if e.Diff() == 0 {
		//e.SetDiff(last.Diff())
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
	//e.SetDiff(0)
	//self.cans=nil

}
func (self *Layer) Add(e Element){

	if e.Val() == 0 {
		return
	}
	if e.Max()==e.Min() {
		fmt.Println(e.Max(),e.Min())
		//panic(0)
		return
	}
	if self.lastEl !=nil {
		dl := e.LastTime() -  self.lastEl.LastTime()
		if (dl < 0)  || (dl>300) {
			//fmt.Println("-------------",float64(dl)/3600)
			self.canChan<-nil
			self.canChan <- e
			self.lastEl = e
			return
		}
	}
	if self.lastEl != nil {
		e.SetDiff(e.Val() - self.lastEl.Val())
		e.SetDur(e.LastTime() - self.lastEl.Time())
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

func (self *Layer) add_1(c Element) (o bool) {
	o = false
	self.isF = false
	if c == nil {
		if len(self.cans)>1{
			if self.par == nil {
				self.setPar()
			}
			self.par.add_1(NewNode(self.cans))
			self.cans = nil

		}
		self.splitID = 0
		self.splitID_old = 0
		self.last = nil
		return
	}
	self.cans = append(self.cans,c)
	d := c.Val() - self.cans[0].Val()
	if math.Abs(d) >= math.Abs(self.direction){
		self.direction = d
		self.splitID = len(self.cans)-1
		return
	}
	self.Split(c,func(n *Node,isU bool){
		o = true
		if self.tag==0 {
			return
		}
		//if n.Dur()>1000{
		fmt.Printf(
			"%d %5d %10.2f %5d %5d %10.2f\r\n",
			self.tag,
			len(self.par.cans),
			self.par.direction,
			n.Dur(),
			len(n.Eles),
			n.Diff(),
		)
		//}
	})
	//fmt.Println(o)

	//n1 := NewNode(self.cans)
	//if ((n1.Diff()>0) == (c.Val() > n1.Val())) {
	//	return
	//}

	//if self.par == nil {
	//	self.setPar()
	//}
	//n_0 := NewNode(self.cans[:self.splitID+1])
	//if self.par!= nil{
	//	self.isF = self.par.add_1(n_0)
	//}
	//if self.last != nil {
	//	if (self.last.Diff()>0) == (n_0.Diff()>0){
	//		self.last = NewNode(append(self.last.Eles,n_0.Eles[1:]...))
	//	}else{
	//		self.last = n_0
	//	}
	//}else{
	//	self.last = n_0
	//}
	//self.cans = self.cans[self.splitID:]
	////self.cans = n_1.Eles

	//o = true
	//C:=self.cans[0]
	//var diffAbs,maxAbs,diff float64
	//for i,c_ := range self.cans[1:]{
	//	diff = c_.Val() - C.Val()
	//	diffAbs = math.Abs(diff)
	//	if diffAbs > maxAbs {
	//		self.splitID = i+1
	//		self.direction = diff
	//		maxAbs = diffAbs
	//	}
	//}
	//self.splitID_old = self.splitID

	//if self.tag == 0 {
	//	return
	//}

	//if self.isF  {
	//	fmt.Println("-------")
	//	//&& self.par.direction != 0{
	//	//if self.tem == nil {
	//	//	self.getTemplate(self.par.direction>0)
	//	//}
	//}
	//n_2 := NewNode(self.cans)
	//fmt.Printf(
	//	"%d %5d %10.2f %5d %5d %10.2f %5d %10.2f %10.2f\r\n",
	//	self.tag,
	//	len(self.par.cans),
	//	self.par.direction,
	//	len(n_0.Eles),
	//	n_0.Dur(),
	//	n_0.Diff(),
	//	len(self.cans),
	//	self.direction,
	//	//n1.Val(),
	//	n_0.Val(),
	//	//n_2.Val() - n_0.Val(),
	//)
	return
}
func (self *Layer)setSplitID(){
	C:=self.cans[0]
	var diffAbs,maxAbs,diff float64
	self.splitID = 0
	for i,c_ := range self.cans[1:]{
		diff = c_.Val() - C.Val()
		diffAbs = math.Abs(diff)
		if diffAbs > maxAbs {
			self.splitID = i+1
			self.direction = diff
			maxAbs = diffAbs
		}
	}

}

func (self *Layer) Split(c Element,hand func(*Node,bool)){
	n1 := NewNode(self.cans)
	if ((n1.Diff()<0) == (c.Val() > n1.Val())) {
		np := NewNode(self.cans[:self.splitID+1])
		hand(np,self.addPar(np))
		self.cans = self.cans[self.splitID:]
		self.setSplitID()
		return
	}
	if len(self.cans)-self.splitID == 1 {
		return
	}
	n2 := NewNode(self.cans[self.splitID:])
	if ((n2.Diff()>0) == (c.Val() > n2.Val())) {
		return
	}
	np := NewNode(self.cans[:self.splitID+1])
	hand(np,self.addPar(np))
	self.cans = self.cans[self.splitID:]
	self.setSplitID()

	np = NewNode(self.cans[:self.splitID+1])
	hand(np,self.addPar(np))
	self.cans = self.cans[self.splitID:]
	self.setSplitID()


}

func (self *Layer)addPar(c Element) bool{
	if self.par == nil {
		self.setPar()
		if self.par == nil {
			return false
		}
	}
	return self.par.add_1(c)
}


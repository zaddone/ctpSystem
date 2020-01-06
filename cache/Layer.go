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
	splitID_1 int
	//isF bool
	//sync.Mutex
	parCov float64
	parSum float64
	//last *Node
	//nodeL *Node
	level_2 float64

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
		self.add(c)
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
			self.par.add(NewNode(self.cans))
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
	self.par.add(NewNode(self.cans[:le]))
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
func (self *Layer) addCans(c Element)(l Element){
	l = c
	le := len(self.cans)
	if le == 0 {
		self.cans = []Element{c}
		return
	}
	if c.Diff()==0 {
		self.cans = append(self.cans,c)
		return
	}
	lc := self.cans[le-1]
	if (lc.Diff() == 0) ||
	((lc.Diff()>0) == (c.Diff()>0)){
		l = MergeElement(lc,c,false)
		self.cans[le-1] = l
		return
	}
	self.cans = append(self.cans,c)
	return
}
func (self *Layer) add(c Element) (o bool) {

	o = false
	if c == nil {
		if len(self.cans)>1{
			self.addPar(NewNode(self.cans))
			self.cans = nil
		}

		self.splitID = 0
		return
	}
	c = self.addCans(c)
	//if l != c {
	//	return
	//}
	//if len(self.cans) == 1 {
	//	return
	//}
	d := c.Val() - self.cans[0].Val()
	if math.Abs(d) >= math.Abs(self.direction){
		self.direction = d
		self.splitID = len(self.cans)-1
		return
	}
	if (self.splitID+1) == len(self.cans){
		return
	}

	n0 := NewNode(self.cans)
	if (n0.Val()<c.Val()) == (self.direction>0) {
		return
	}
	o = true
	if self.addPar(NewNode(self.cans[:self.splitID+1])){
		self.setSplitID()

	}else{
		self.setSplitID()
	}


	if self.direction != 0 {
		return
	}
	n1 := NewNode(self.par.cans)
	fmt.Printf("%d %10d %10d %10.2f %10.2f %10d %10d %10d\r\n",
		self.par.tag,
		n1.Dur(),
		len(n1.Eles),
		self.par.direction,
		self.direction,
		len(n0.Eles),
		self.splitID,
		len(self.cans),
	)
	return

}
func (self *Layer)addPar(c Element) bool{
	if self.par == nil {
		self.setPar()
		if self.par == nil {
			return false
		}
	}
	return self.par.add(c)
}

func (self *Layer)setSplitID(){
	self.cans = self.cans[self.splitID:]
	self.splitID = 0
	self.direction = 0
	if len(self.cans) == 1 {
		fmt.Println("------")
		return
	}
	C:=self.cans[0]
	var diffAbs,maxAbs,diff float64
	for i,c_ := range self.cans[1:]{
		diff = c_.Val() - C.Val()
		diffAbs = math.Abs(diff)
		if diffAbs >= maxAbs {
			self.splitID = i+1
			self.direction = diff
			maxAbs = diffAbs
		}
	}

}

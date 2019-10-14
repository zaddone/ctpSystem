package cache
import(
	"math"
	"fmt"
	"io"
	//"github.com/zaddone/analog/fitting"
	//"github.com/zaddone/ctpSystem/config"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"encoding/gob"
	"bytes"
	//"sync"
)
type Temple struct{
	can Element
	lcan Element
	XMin,XMax,YMin,YMax float64
	Wei []float64
	Stats int
	Dis bool
}
func (self *Temple) Save(){
	return
	na := self.can.Name()
	k := make([]byte,8)
	binary.BigEndian.PutUint64(k,uint64(self.can.LastTime()))
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(self)
	if err != nil {
		panic(err)
	}

	_db,err := bolt.Open(fmt.Sprintf("temple_%s.db",na),0600,nil)
	if err != nil {
		panic(err)
	}
	err = _db.Update(func(t *bolt.Tx)error{
		b,err := t.CreateBucketIfNotExists([]byte(na))
		if err != nil {
			return err
		}
		return b.Put(k,buf.Bytes())
	})
	if err != nil {
		panic(err)
	}
	_db.Close()

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
	//ca *Cache
	tem  *Temple
	//sync.Mutex
}

func NewLayer() (L *Layer) {
	L = &Layer{
		//ca:ca,
		canChan:make(chan Element,100),
	}
	go L.runChan()
	return L
}
func (self *Layer)checkTemStop(){
	if math.Abs(self.direction) < math.Abs(self.par.cans[len(self.par.cans)-1].Diff()){
		return
	}
	self.checkTem()
}
func (self *Layer)checkTem(){
	c_ := self.getLast()
	dis_:= c_.Val() - self.tem.can.Val()
	//fmt.Println(dis_,c_.LastTime()-self.tem.can.LastTime())
	//if dis != dis_ || dis != self.tem.Dis {
	absDis := math.Abs(dis_)
	t := self.tag-1
	if dis_>0 == self.tem.Dis {
	//if dis != self.tem.Dis {
		self.tem.Stats = 2
		Count[t][0]+=absDis
	}else{
		//fmt.Println(absDis)
		self.tem.Stats = 1
		Count[t][0]-=absDis
	}
	Count[t][self.tem.Stats]++
	//fmt.Println(Count[t],c_.LastTime()-self.tem.can.LastTime())
	self.tem.Save()
	self.tem = nil
}

func (self *Layer) getLast() Element {
	if self.child != nil {
		return self.child.getLast()
	}else{
		return self.cans[len(self.cans)-1]
		//self.Lock()
		//defer self.Unlock()
		//return self.lastEl
	}
}

func (self *Layer) readAll(h func(Element)error)error{

	for _,c := range self.cans {
		err := c.Each(h)
		if err != nil {
			return err
		}
	}
	if self.child != nil {
		return self.child.readAll(h)
	}
	return nil

}
func (self *Layer) getNormalization(dis bool)(X,Y []float64){
	self.tem = &Temple{Dis:dis}
	return
	var x,y float64
	me := map[Element]bool{}
	err := self.readAll(func(e Element)error{
		if me[e]{
			return nil
		}
		me[e] = true
		x = float64(e.Time())
		y = e.Val()
		X = append(X,x)
		Y = append(Y,y)
		if self.tem.YMin == 0 || self.tem.YMin>y {
			self.tem.YMin = y
		}
		if self.tem.YMax < y {
			self.tem.YMax = y
		}
		if self.tem.XMin == 0 || self.tem.XMin>x {
			//self.tem.can = e
			self.tem.XMin = x
		}
		if self.tem.XMax < x {
			self.tem.XMax = x
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	xLong := self.tem.XMax - self.tem.XMin
	yLong := self.tem.YMax - self.tem.YMin
	for i,x := range X {
		X[i] = (x-self.tem.XMin)/xLong
		Y[i] = (y-self.tem.YMin)/yLong
	}
	return

}
func (self *Layer) getTemplate(dis bool){
	self.getNormalization(dis)
	//X,Y := self.getNormalization()
	//self.tem.Wei = make([]float64,config.Conf.Weight+self.tag)
	//if !fitting.GetCurveFittingWeight(X,Y,self.tem.Wei) {
	//	panic(fmt.Errorf("height fitting err"))
	//}
	////fmt.Println(self.tem.Wei,len(X))
	//self.tem = &Temple{}

	self.tem.can = self.getLast()
	self.tem.lcan = self.cans[len(self.cans)-1]
	var stop Element
	self.cans[0].Each(func(e Element)error{
		stop = e
		return io.EOF
	})
	NewInsOrder(self.tem.can,stop)

	//self.tem.Save()
}

func (self *Layer) runChan(){
	for{
		//self.Lock()
		self.baseAdd(<-self.canChan)
		//self.Unlock()
	}

}
func (self *Layer) baseAdd(e Element){
	if e== nil {
		self.par = nil
		self.cans = nil
		return
	}
	le := len(self.cans)
	self.cans = append(self.cans,e)

	if le == 0 {
		return
	}
	last := self.cans[le-1]
	e.SetDiff(e.Val() - last.Val())
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
		self.par = &Layer{
			//ca:self.ca
			child:self,
		}
		self.par.tag = self.tag+1
	}
	//fmt.Println("in----------->",len(self.cans))
	self.par.add(NewNode(self.cans[:le]))
	self.cans = []Element{e}

}
func (self *Layer) Add(e Element){
	if self.lastEl !=nil {
		dl := e.LastTime() -  self.lastEl.LastTime()
		if dl <0 || dl > 10 {
			//fmt.Println(dl)
			self.canChan<-nil
		}
	}
	self.canChan <- e
	//self.Lock()
	self.lastEl = e
	//self.Unlock()
}
func (self *Layer) initAdd (c Element){

	le := len(self.cans)
	self.cans = append(self.cans,c)
	if le < 3 {
		return
	}
	var maxD,absMaxD float64
	self.sum  += math.Abs(c.Diff())
	var splitID int
	for i,_c := range self.cans[:le] {
		//sum += math.Abs(_c.Diff())
		d := c.Val() - _c.Val()
		absD := math.Abs(d)
		if (absD > absMaxD){
			maxD = d
			absMaxD = absD
			splitID = i
		}
	}
	if (self.sum/float64(len(self.cans))) > absMaxD {
		return
	}
	self.direction = maxD
	self.cans = self.cans[splitID:]
	for _,_c := range self.cans{
		self.sum += math.Abs(_c.Diff())
	}

}

func (self *Layer) add(c Element) bool {

	if self.direction == 0 {
		self.initAdd(c)
		return false
	}
	le := len(self.cans)
	self.cans = append(self.cans,c)
	var absMaxD, maxD float64
	self.sum  += math.Abs(c.Diff())
	var splitID int
	for i,_c := range self.cans[:le] {
		//sum += math.Abs(_c.Diff())
		d := c.Val() - _c.Val()
		absD := math.Abs(d)
		if ((d>0) == (self.direction>0)) {
			if math.Abs(self.direction) < absD {
				self.direction = d
			}
			continue
		}
		if (absD > absMaxD){
			maxD = d
			absMaxD = absD
			splitID = i
		}
	}
	sumv := self.sum/float64(len(self.cans))
	//fmt.Println(sumv,self.direction)
	if splitID == 0 ||
	sumv > absMaxD {
		//if self.tem != nil {
		//	self.checkTemStop()
		//}
		return false
	}
	//fmt.Println(maxD)
	//dir := self.direction
	self.direction = maxD
	if self.par == nil {
		self.par = &Layer{
			//ca:self.ca
			child:self,
		}
		self.par.tag = self.tag+1
	}
	self.par.add(NewNode(self.cans[:splitID+1]))



	self.cans = self.cans[splitID:]
	self.sum=0
	for _,_c := range self.cans{
		self.sum += math.Abs(_c.Diff())
	}

	if self.tem != nil {
		self.checkTem()
	}
	self.getTemplate(self.direction>0)

	//}else{
		//if self.par != nil &&
		//self.par.direction != 0 &&
		//(self.par.direction>0) == !(self.direction>0){
		////if math.Abs(dir) < absMaxD {
			//self.getTemplate(self.direction>0)
		////}else if math.Abs(dir) < absMaxD{
		//	//self.getTemplate(self.direction>0)
		//}
	//}
	//fmt.Println(self.sum/float64(len(self.cans)),self.par.sum/float64(len(self.par.cans)))
	return true

}

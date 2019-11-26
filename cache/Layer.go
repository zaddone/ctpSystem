package cache
import(
	"math"
	"fmt"
	//"io"
	//"github.com/zaddone/analog/fitting"
	"github.com/zaddone/ctpSystem/config"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"encoding/gob"
	"bytes"
	"time"
	//"sync"
)
type Temple struct{
	can Element
	stop Element
	lcan Element
	XMin,XMax,YMin,YMax float64
	Wei []float64
	Stats int
	Dis bool
}
func (self *Temple)Check(l *Layer){

	//if self.Stats<0{
	//	return
	//}
	//la := l.cans[len(l.cans)-1]
	//if (la.Val()>self.stop.Val()) == self.Dis {
	//	//fmt.Println(Count[0])
	//	l.checkTem()
	//	//fmt.Println(self.Stats,"sl",Count[0])
	//	//fmt.Println(Count[0])
	//	return
	//}
	if self.Stats < 1 {
		//if l.splitID == 0 {
		//	return
		//}
		//e := NewNode(l.cans[:(l.splitID+1)])
		//e_ := NewNode(l.cans[l.splitID:])
		//dis := e_.Val() > e.Val()
		//if self.Dis != dis {
		//	l.checkTem()
		//}
		return
	}


	if (l.direction>0) != self.Dis {
		//if (la.Val()>self.can.Val()) == self.Dis {
		//	fmt.Println(Count[0])
		//	l.checkTem()
		//	fmt.Println(self.Stats,"sl",Count[0])
		//	return
		//}
		if l.splitID != 0 {
			return
		}
		e := l.par.cans[len(l.par.cans)-1]
		e_:= NewNode(l.cans)
		//dis := e_.Val() > e.Val()
		//fmt.Println(self.Dis,e.Diff(),e_.Diff(),e.Val()<e_.Val())
		if math.Abs(e_.Diff()) > math.Abs(e.Diff()) {
			l.checkTem()
		}
		return
	}else{
		if l.splitID == 0 {
			return
		}
		e := NewNode(l.cans[:(l.splitID+1)])
		e_ := NewNode(l.cans[l.splitID:])
		//dis := e_.Val() > e.Val()
		if  math.Abs(e_.Diff()) >math.Abs(e.Diff()) {
		//if (e_.Val() > e.Val()) != (self.Dis) {
			//fmt.Println(self.Stats,self.Dis,e_.Diff(),e.Diff())
			l.checkTem()
		}
		return
	}

}
func (self *Temple) SetStats(){
	self.Stats++
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
	ca *Cache
	tem  *Temple
	splitID int
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

func (self *Layer)checkTemStop(){
	if math.Abs(self.direction) < math.Abs(self.par.cans[len(self.par.cans)-1].Diff()){
		return
	}
	self.checkTem()
}

func (self *Layer) checkTem() (isok bool) {
	c_ := self.getLast()

	//if self.tag ==1 {
		//if self.par!=nil &&
		//self.par.tem != nil &&
		//self.par.tem.Dis != self.tem.Dis {
			t := self.tag-1
			var Diff float64
			if self.tem.Dis {
				Diff = c_.Min() - self.tem.can.Max()
			}else{
				Diff = c_.Max() - self.tem.can.Min()
			}
			if self.tem.Dis == (Diff>0){
			//if isok_ {
				Count[t][3]++
			}else{
				Count[t][2]++
			}
			dis_:= c_.Val() - self.tem.can.Val()
			//absDis := math.Abs(dis_)
			if  (dis_>0) == self.tem.Dis {
			//if isok {
				Count[t][1]++
				//self.tem.Stats = 1
				//Count[t][0] += absDis
			}else{
				Count[t][0]++
				//self.tem.Stats = 0
				//Count[t][0] -= absDis
			}
			//Count[t][self.tem.Stats]++
			//fmt.Println(Count[t],c_.Time() - self.tem.can.Time())
		//}
	//}


	self.tem = nil
	//self.tem.Stats = -1
	if config.Conf.IsTrader{
		self.ca.Order.SendCloseOrder(c_.(*Candle),self.ca)
	}
	return
}

func (self *Layer) getLast() Element {
	if self.child != nil {
		return self.child.getLast()
	}else{
		return self.lastEl
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
func (self *Layer) getTemplate(dis bool,e Element){
	//if self.tem.Stats <0{
	//	return
	//}
	diff := e.Val() - self.ca.GetLast().(Element).Val()
	if (diff>0) != dis {
		return
		panic(9)
	}
	//fmt.Println(diff,NewNode(self.cans).Diff())

	//var diff1,diff2,n1,n2 float64
	//for _,c := range self.cans{
	//	diff1 += math.Abs(c.Diff())
	//	n1++
	//	c.Each(func(e Element)error{
	//		diff2+=(e.Max()-e.Min())
	//		n2++
	//		return nil
	//	})
	//}
	//if ((diff2/n2) / math.Abs(self.direction)) > 0.5{
	//	return
	//}

	self.getNormalization(dis)
	self.tem.stop = e

	//X,Y := self.getNormalization()
	//self.tem.Wei = make([]float64,config.Conf.Weight+self.tag)
	//if !fitting.GetCurveFittingWeight(X,Y,self.tem.Wei) {
	//	panic(fmt.Errorf("height fitting err"))
	//}
	////fmt.Println(self.tem.Wei,len(X))
	//self.tem = &Temple{}

	self.tem.can = self.getLast()
	self.tem.lcan = self.cans[len(self.cans)-1]

	//var stop Element
	//self.cans[0].Each(func(e Element)error{
	//	stop = e
	//	return io.EOF
	//})
	//dif := self.tem.can.Val() - stop.Val()
	//if dif==0 ||
	// (dif>0)!=(self.direction>0){
	//	self.tem = nil
	//	return
	//}
	if config.Conf.IsTrader{
		//self.ca.Order.Update(1,self.tem.Dis,self.tem.can)
		self.ca.AddOrder(self.tem.Dis,self.tem.stop)
		//OpenInsOrder(self.tem.can,self.tem.Dis)
	}

	//self.tem.Save()
}

func (self *Layer) runChan(){
	for{
		//self.Lock()
		self.baseAdd(<-self.canChan)
		//self.add(<-self.canChan)
		//self.Unlock()
	}

}

func (self *Layer) _baseAdd(e Element){
	if e== nil {
		self.par = nil
		self.cans = nil
		return
	}
	self.add(e)
}
func (self *Layer) baseAdd(e Element){
	if e == nil {
		//CloseInsOrder(self.getLast())
		//self.ca.Order.Update(3,self.getLast())
		c:= self.ca.GetLast().(*Candle)
		self.ca.EachOrder(func(k string,o *InsOrder)bool{
			o.SendCloseOrder(c,self.ca)
			return true
		})
		//self.ca.Order.SendCloseOrder(self.getLast())
		self.par = nil
		self.cans = nil
		return
	}

	if !config.Conf.IsTrader{
		self.CheckPL(e)
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
			ca:self.ca,
			child:self,
		}
		self.par.tag = self.tag+1
	}
	//fmt.Println("in----------->",le)
	if le==1 {
		self.par.add(self.cans[0])
	}else{
		self.par.add(NewNode(self.cans[:le]))
	}
	self.cans = []Element{e}
	e.SetDiff(0)
	//self.cans=nil

}
func (self *Layer) Add(e Element){
	if self.lastEl !=nil {
		//end := time.Unix(e.LastTime()).Day() != 
		//begin := time.Unix(self.lastEl.LastTime())

		dl := e.LastTime() -  self.lastEl.LastTime()
		if dl <0 || (time.Unix(e.LastTime(),0).Day() != time.Unix(self.lastEl.LastTime(),0).Day()) {
			//fmt.Println("timeOut",dl,self.lastEl.Val(),e.Val())
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

	//self.cans = append(self.cans,c)
	//if le < 3 {
	//	return
	//}
	if le >0 {
		last := self.cans[le-1]
		if (last.Diff()>0)== (c.Diff()>0){
			c = MergeElement(last,c)
			le--
			self.cans[le] = c
			self.sum =self.sum - (last.Max()-last.Min()) + (c.Max()-c.Min())
			//self.sum  += c.Max()-c.Min()
		}else{
			self.sum  += c.Max()-c.Min()
			self.cans = append(self.cans,c)
		}
		//le--
	}else{
		self.sum  += c.Max()-c.Min()
		self.cans = append(self.cans,c)
	}
	//self.sum  += math.Abs(c.Diff())
	var maxD,absMaxD float64
	self.splitID = 0
	for i,_c := range self.cans[:le] {
		//sum += math.Abs(_c.Diff())
		d := c.Val() - _c.Val()
		absD := math.Abs(d)
		if (absD > absMaxD){
			maxD = d
			absMaxD = absD
			self.splitID = i
		}
	}
	if (self.sum/float64(len(self.cans))) > absMaxD {
		return
	}
	self.direction = maxD
	self.cans = self.cans[self.splitID:]
	self.sum = 0
	for _,_c := range self.cans{
		self.sum  += _c.Max()-_c.Min()
		//self.sum  += math.Abs(_c.Diff())
	}

}
func (self *Layer)CheckPL(e Element) {
	if self.tem == nil {
		if self.par != nil {
			self.par.CheckPL(e)
		}
	}else{
		//if (self.tem.Stats < 0) {
		//	return
		//}
		if (e.Val()>self.tem.stop.Val()) == self.tem.Dis {
			//fmt.Println(e.Val(),self.tem.stop.Val(),self.tem.Dis)
			//fmt.Println(Count[0])
			self.checkTem()
			//fmt.Println(Count[0])
			return
		}
	}
	//e  := NewNode(self.cans[:(self.splitID+1)])
	//e_ := NewNode(self.cans[self.splitID:])
	//return (e_.Val() > e.Val()) == (self.direction>0)
}
func (self *Layer)CheckCansLong(dis bool) bool {

	var sum1,n1 float64
	var sum2,n2 float64
	for _,c := range self.cans{
		if (c.Diff()>0) == dis{
			sum1 += math.Abs(c.Diff())
			n1++
		}else{
			sum2 += math.Abs(c.Diff())
			n2++
		}
	}
	if n1==0 || n2==0 {
		return false
	}
	return (sum1/n1) > (sum2/n2)

}

func (self *Layer) GetAmplitude(dis bool) float64 {
	var sum,n float64
	for _,c := range self.cans{
		if (c.Diff()>0) == dis{
			sum += math.Abs(c.Diff())
			n++
		}
	}
	return sum/n
}

func (self *Layer) add(c Element) bool {

	if self.direction == 0 {
		self.initAdd(c)
		return false
	}
	le := len(self.cans)
	//self.cans = append(self.cans,c)
	//self.sum  += math.Abs(c.Diff())

	if le >0 {
		last := self.cans[le-1]
		if (last.Diff()>0)== (c.Diff()>0){
			c = MergeElement(last,c)
			le--
			self.cans[le] = c
			self.sum =self.sum- (last.Max()-last.Min())+(c.Max()-c.Min())
		}else{
			self.sum  += c.Max()-c.Min()
			self.cans = append(self.cans,c)
		}
	}else{
		self.sum  += c.Max()-c.Min()
		self.cans = append(self.cans,c)
	}

	var absMaxD, maxD float64
	self.splitID = 0
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
			self.splitID = i
		}
	}

	//if (self.splitID != 0) &&
	if (self.tag == 1) &&
	(self.par != nil) &&
	//(self.par.direction !=0 ) &&
	//(self.par.splitID != 0 ) &&
	//(self.par.direction >0 ) == (self.direction>0) &&
	(self.tem == nil) &&
	//self.CheckSplit() &&
	//self.par.CheckCansLong(self.direction<0) &&
	math.Abs(self.par.cans[len(self.par.cans)-1].Diff())>math.Abs(self.direction)&&
	(math.Abs(self.direction) > self.par.GetAmplitude(self.direction>0)) {
		self.getTemplate(self.direction<0,self.cans[0])
	}

	if self.tem != nil {
		self.tem.Check(self)
	}

	if self.splitID == 0 {
		return false
	}
	sumv := self.sum/float64(len(self.cans))
	if sumv > absMaxD {
		return false
	}

	//fmt.Println(sumv,absMaxD)
	//fmt.Println(maxD)
	//dir := self.direction

	self.direction = maxD
	//var e1 Element = nil
	if self.par == nil {
		self.par = &Layer{
			ca:self.ca,
			child:self,
		}
		self.par.tag = self.tag+1
	//}else{
	//	e1 = self.par.cans[len(self.par.cans)-1]
	}
	e := NewNode(self.cans[:self.splitID+1])
	self.par.initAdd(e)
	self.cans = self.cans[self.splitID:]
	self.sum = 0
	for _,_c := range self.cans{
		//self.sum += math.Abs(_c.Diff())
		self.sum  += _c.Max()-_c.Min()
	}
	//if self.tag == 1 {
		//isU:= true
		if self.tem != nil  {
			//isU = self.checkTem()
			//if self.tem.Stats >=0{
				self.checkTem()
				//self.tem.SetStats()
			//}else{
			//	self.tem = nil
			//}
		}
		//if isU{
		//if self.par.direction!=0 &&
		//(self.par.direction>0) == (self.direction>0)&&
		//(self.par.direction>0) != (self.direction<0){
		//if math.Abs(dir) > absMaxD {
		//if (e1 != nil) &&
		//(e.Val() > e1.Val()) == (self.direction<0){
		//}
	//}



	return true

}

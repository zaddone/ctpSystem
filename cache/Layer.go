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
	"time"
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
	//if self.tem.Stats <1 {
	//	return
	//}
	//if self.splitID==0{
	//	return
	//}
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
				Count[t][4] += math.Abs(Diff)
				Count[t][3]++
			}else{
				Count[t][4] -= math.Abs(Diff)
				Count[t][2]++
			}
			dis_:= c_.Val() - self.tem.can.Val()

			//absDis := math.Abs(dis_)
			if  (dis_>0) == self.tem.Dis {
			//if isok {
				Count[t][1]++
				//self.tem.Stats = 1
				//Count[t][0] += absDis
				Count[t][5] += math.Abs(dis_)
			}else{
				Count[t][0]++
				//self.tem.Stats = 0
				//Count[t][0] -= absDis
				Count[t][5] -= math.Abs(dis_)
			}
			//Count[t][self.tem.Stats]++
			fmt.Println(Count[t],c_.Time() - self.tem.can.Time())
		//}
	//}

	self.tem = nil
	//self.tem.Stats = -1
	//if config.Conf.IsTrader {
	if self.ca.Order != nil {
		self.ca.Order.SendCloseOrder(c_.(*Candle),self.ca)
	}
	//}
	return
}

func (self *Layer) getLast() Element {
	if self.child != nil {
		return self.child.getLast()
	}else{
		return self.cans[len(self.cans)-1]
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
	//if (self.par == nil) {
	//	return
	//}
	//if (self.par.direction ==0 ) {
	//	return
	//}
	//if (self.par.direction>0) == (self.direction<0){
	//	return
	//}
	//if (self.par.splitID!=0){
	//	return
	//}

	//diff_ := self.par.GetAmplitude(self.direction>0)
	//var sum float64
	//for _,c := range self.cans {
	//	sum += c.Max() - c.Min()
	//}
	//if diff_ < sum/float64(len(self.cans)){
	//	return
	//}

	//if math.Abs(self.direction)/(e.Max()-e.Min())<2{
	//	return
	//}
	//diff := e.Val() - self.ca.GetLast().(Element).Val()
	//if (diff>0) != dis {
	//	return
	//}

	//if ( math.Abs(self.direction) < diff_ ) {
	//	return
	//}
	//diff_ = math.Abs(self.direction)
	L := self.ca.L.isTem()
	if L != nil{
		//fmt.Println(L.tag)
		if L.tag < self.tag {
			self.tem = L.tem
			L.tem = nil
		//}else{
		//	if L.tem.Dis == (self.direction>0)
		}
		return
	}


	//if self.par.splitID ==0 {
	//	return
	//}

	//if math.Abs(self.par.cans[len(self.par.cans)-1].Diff()) < math.Abs(self.direction) {
	//	return
	//}




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
	//self.par.cans[len(self.par.cans)-1].Each(func(el Element)error{
	//	self.tem.stop = el
	//	return io.EOF
	//})

	//X,Y := self.getNormalization()
	//self.tem.Wei = make([]float64,config.Conf.Weight+self.tag)
	//if !fitting.GetCurveFittingWeight(X,Y,self.tem.Wei) {
	//	panic(fmt.Errorf("height fitting err"))
	//}
	////fmt.Println(self.tem.Wei,len(X))
	//self.tem = &Temple{}

	self.tem.can  = self.getLast()
	//if self.tem.Dis {
	//	self.tem.long = self.tem.can.Min() - diff_
	//}else{
	//	self.tem.long = self.tem.can.Max() + diff_
	//}
	//self.tem.lcan = self.cans[len(self.cans)-1]

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
		//self.add_(<-self.canChan)
		//self.Unlock()
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
			//fmt.Println("___________________")
		}
		//self.cans = nil
		//self.par = nil
		return
	}
	//if !config.Conf.IsTrader{
		//self.CheckPL(e)
		//self.CheckSL(e)
	//}
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
		self.setPar()
	}
	self.par.add_1(NewNode(self.cans[:le]))
	self.cans = []Element{e}
	e.SetDiff(0)
	//self.cans=nil

}
func (self *Layer) Add(e Element){

	if e.Val() == 0 {
		//panic(0)
		//fmt.Println(e)
		return
	}
	//if e.Min() == e.Max() {
	//	fmt.Println(e)
	//	panic(0)
	//	return
	//}
	if self.lastEl !=nil {
		//end := time.Unix(e.LastTime(),0).Day()
		//begin := time.Unix(self.lastEl.LastTime(),0).Day()
		dl := e.LastTime() -  self.lastEl.LastTime()
		//if dl <0 || dl>60 {
		//if (dl < 0)  || (end!=begin) {
		if (dl < 0)  || (dl>60) {
			//fmt.Println("timeOut",dl,self.lastEl.Val(),e.Val())
			//if self.par == nil || self.par.tem == nil {
			if self.ca.Order == nil {
				self.canChan<-nil
			}
		}
	}

	self.canChan <- e
	//self.Lock()
	self.lastEl = e
	//self.Unlock()
}
func (self *Layer) initAdd (c Element){

	le := len(self.cans)
	if le >0 {
		last := self.cans[le-1]
		if (last.Diff()>0)== (c.Diff()>0){
			c = MergeElement(last,c)
			le--
			self.cans[le] = c
			//self.sum =self.sum- (last.Max()-last.Min())+(c.Max()-c.Min())
			self.sum=self.sum-math.Abs(last.Diff()) + math.Abs(c.Diff())
			//self.sum  += c.Max()-c.Min()
		}else{
			//self.sum  += c.Max()-c.Min()
			self.sum  += math.Abs(c.Diff())
			self.cans = append(self.cans,c)
		}
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
		//self.sum  += _c.Max()-_c.Min()
		self.sum  += math.Abs(_c.Diff())
	}

}
func (self *Layer)Check(){
	var sum,n int64
	var cov,vcov,val float64
	for _,c := range self.par.cans{
		if (c.Diff()>0) == (self.direction>0) {
			continue
		}
		sum += c.Dur()
		val += math.Abs(c.Diff())
		n++
	}
	sum /= n
	val /= float64(n)
	for _,c := range self.par.cans{
		if (c.Diff()>0) == (self.direction>0) {
			continue
		}
		cov += math.Pow(float64(c.Dur()-sum),2)
		vcov += math.Pow(math.Abs(c.Diff())-val,2)
	}
	cov = math.Sqrt(cov/float64(n))
	vcov = math.Sqrt(vcov/float64(n))
	fmt.Println(cov,sum,vcov,val,self.direction,self.par.direction,n,len(self.par.cans))
}

func (self *Layer)CheckSL(e Element) {
	if self.tem == nil {
		if self.par != nil {
			self.par.CheckSL(e)
		}
	}else{
		//if self.tem.Stats != -1 {
		if (e.Val()<self.tem.long) == self.tem.Dis {
			self.checkTem()
			return
		}
		//}
	}
}
func (self *Layer)CheckPL(e Element) {
	//return
	if self.tem == nil {
		if self.par != nil {
			self.par.CheckPL(e)
		}
	}else{
		//if (self.tem.Stats < 0) {
		//	return
		//}
		if (e.Val()>self.tem.long) == self.tem.Dis {
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
func (self *Layer) add_init(c Element) {

	le := len(self.cans)
	self.cans = append(self.cans,c)
	if le == 0 {
		return
	}
	//self.direction  = c.Val() - self.cans[0].Val()
	//self.splitID = le

	var sum,maxAbs,max,dif,difAbs float64
	sum = c.Max() - c.Min()
	var I int
	for i,c_ := range self.cans[:le]{
		sum += c_.Max()-c_.Min()
		dif = c.Val() - c_.Val()
		difAbs = math.Abs(dif)
		if maxAbs < difAbs {
			max = dif
			maxAbs = difAbs
			I = i
		}
	}
	if sum/float64(le+1) > maxAbs {
		return
	}
	self.direction = max
	self.cans = self.cans[I:]
	self.splitID = len(self.cans)-1

	//fmt.Println(self.tag,I,self.splitID)

}
func (self *Layer)split() bool {

	if len(self.cans) == 1 {
		return false
	}
	if self.par == nil {
		self.setPar()
	}
	n := NewNode(self.cans[:self.splitID+1])
	self.par.add_(n)
	self.cans = self.cans[self.splitID:]
	c := self.cans[0]
	var maxAbs,d,dAbs float64
	self.splitID = 0
	dir:= self.direction
	self.direction = 0
	for i,c_ := range self.cans[1:]{
		d = c_.Val() - c.Val()
		dAbs = math.Abs(d)
		if maxAbs <= dAbs{
			self.direction = d
			maxAbs = dAbs
			self.splitID = i+1
		}
	}
	fmt.Printf(
		"%s %d %d %d %d %d %t %.0f %.0f %.0f\r\n",
		time.Unix(c.Time(),0),
		len(n.Eles),
		self.splitID,
		len(self.cans),
		n.Dur(),
		self.tag,
		self.isF,
		dir,
		self.direction,
		self.par.direction,

	)
	return true

	//N_ := NewNode(self.par.cans)
	//sum := math.Abs(N_.Diff()/float64(N_.Dur()))
	//var cov float64
	//for _,c_ := range self.par.cans {
	//	o := math.Pow(sum - math.Abs(c_.Diff()/float64(c_.Dur())),2)
	//	fmt.Println(o,sum)
	//	cov+=o
	//}
	//cov = math.Sqrt( cov/float64(len(self.par.cans)))
	//fmt.Println(time.Unix(c.Time(),0),self.splitID,len(self.cans),self.tag,self.isF,dir,self.direction,self.par.direction,cov,sum)

}

func (self *Layer)getDir() float64 {
	if self.par == nil {
		return self.direction
	}
	return self.par.getDir()
}
func (self *Layer) setPar(){
	self.par = &Layer{
		ca:self.ca,
		child:self,
		tag:self.tag+1,
	}
	fmt.Println(self.par.tag)

}
func (self *Layer) add_(c Element) {
	if c == nil {
		if len(self.cans)>0{
			for{
				if !self.split(){
					break
				}
			}
			self.cans = nil
		}
		self.direction = 0

		//self.par = nil
		//self.cans = nil
		//self.direction = 0
		fmt.Println("_________")
		return
	}
	if self.direction == 0 {
		self.add_init(c)
		return
	}
	le := len(self.cans)
	self.cans = append(self.cans,c)

	d := c.Val() - self.cans[0].Val()
	self.isF = false
	if ((d>0)==(self.direction>0)) {
		if (math.Abs(d)>math.Abs(self.direction)){
			self.splitID = le
			self.direction = d
			return
		}else{
			N1 := NewNode(self.cans[:self.splitID+1])
			N2 := NewNode(self.cans[self.splitID:])
			t1 := float64(N1.Dur())
			t2 := float64(N2.Dur())
			t_ := t1 + t2
			t1 /= t_
			t2 /= t_
			d1 := math.Abs(N1.Diff())
			d2 := math.Abs(N2.Diff())
			d_ := d1 + d2
			d1 /= d_
			d2 /= d_
			p1 := math.Sqrt(math.Pow(t1,2)+math.Pow(d1,2))
			p2 := math.Sqrt(math.Pow(t2,2)+math.Pow(d2,2))
			if  p1 > p2 {
				return
			}else{
				//fmt.Println(t1,t2,d1,d2,p1,p2)
				fmt.Println(N1.Dur(),N2.Dur(),N1.Diff(),N2.Diff())
				self.isF = true
			}
		}
	}

	self.split()

}

func (self *Layer) add_1(c Element) {
	self.cans = append(self.cans,c)
	n1 := NewNode(self.cans)
	if math.Abs(self.direction) <= math.Abs(n1.Diff()){
		self.direction = n1.Diff()
		self.splitID = len(self.cans)-1
		return
	}
	//var v float64
	//fmt.Println(c.Val(),n1.Val())
	if n1.Diff()>0{
		//if c.Min() > n1.Val(){
		if c.Val() > n1.Val(){
			return
		}
	}else{
		//if c.Max() < n1.Val(){
		if c.Val() < n1.Val(){
			return
		}
	}
	//fmt.Println(c)
	//self.split()
	C:=self.cans[0]
	var diffAbs,maxAbs float64
	var I int
	for i,c_ := range self.cans[1:]{
		//diff = c_.Val() - C.Val()
		diffAbs = math.Abs(c_.Val() - C.Val())
		if diffAbs > maxAbs {
			I = i+1
			//max = diff
			//self.direction = diff
			maxAbs = diffAbs
		}
	}
	//if I==0 {
	//	fmt.Println("--",len(self.cans))
	//	panic(0)
	//	return
	//}
	if self.par == nil {
		self.setPar()
	}
	n_0 := NewNode(self.cans[:I+1])
	self.par.add_1(n_0)
	self.cans = self.cans[I:]
	self.direction = c.Val() - self.cans[0].Val()
	if self.tag == 1 {
		self.Check_1()
	}

	//if len(self.cans)==1{
	//	panic(1)
	//}
	//if len(self.cans)>1{
	//n_1 := NewNode(self.cans)
	//if self.par.direction>0{
	//if (self.direction>0) == (n_0.Diff()>0){
	//	panic(0)
	//}
	//fmt.Printf("%d %10.2f %5d %5d %5d %10.2f %10.2f\r\n",self.tag,self.par.direction,len(self.par.cans),len(self.cans),len(n_0.Eles),self.direction,n_0.Diff())
	//}
	//}

	//if len(self.par.cans)>1{
	//	st := ""
	//	for _,n := range self.par.cans {
	//		st =fmt.Sprintf("%s %.2f",st, n.Diff())
	//	}
	//	//fmt.Println(st,self.tag,self.par.direction)
	//}

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
			//self.sum =self.sum- (last.Max()-last.Min())+(c.Max()-c.Min())
			self.sum=self.sum-math.Abs(last.Diff()) + math.Abs(c.Diff())
		}else{
			//self.sum  += c.Max()-c.Min()
			self.sum  += math.Abs(c.Diff())
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




	//if (self.tem == nil) {
	//	self.getTemplate(self.direction<0)
	//}else{
	//	self.checkTem()
	//}


	if self.splitID == 0 {
		return false
	}
	dir := math.Abs(self.direction)
	isP := false
	if absMaxD < dir {
		p1 := math.Sqrt(math.Pow(float64(self.cans[self.splitID].LastTime()-self.cans[0].Time()),2) +math.Pow(dir,2))
		p2 := math.Sqrt(math.Pow(float64(c.LastTime() - self.cans[self.splitID].Time()),2)+math.Pow(absMaxD,2))
		if  p1 > p2 {
			return false
		}else{
		//if self.par != nil {
			isP = true
			//if (self.tem == nil) {
			//	self.getTemplate(self.direction<0)
		//	}else{
		//		self.checkTem()
			//}
			//fmt.Println(self.direction,maxD,self.par.direction,len(self.cans),self.tag)
		//}
		}
	}
	//sumv := self.sum/float64(len(self.cans))
	//if sumv > absMaxD {
	//	return false
	//}
	//if self.par != nil {
	//	fmt.Println(self.direction,maxD,self.par.direction,len(self.cans),self.tag)
	//}
	self.direction = maxD
	if self.par == nil {
		tag := self.tag+1
		//if tag <4{
		//fmt.Println(tag)
		self.par = &Layer{
			ca:self.ca,
			child:self,
			tag:tag,
		}
		//}
	}

	self.par.add(NewNode(self.cans[:self.splitID+1]))
	self.cans = self.cans[self.splitID:]
	if self.tem != nil {
		//self.tem.SetStats()
		self.checkTem()
	}else if isP {
		self.getTemplate(self.direction<0)
	}
	self.sum = 0
	for _,_c := range self.cans{
		self.sum += math.Abs(_c.Diff())
		//self.sum  += _c.Max()-_c.Min()
	}

	return true

}
func (self *Layer)CheckPar(isF bool) bool {
	if (self.direction>0) != isF {
		return false
	}
	if self.splitID != len(self.cans)-1 {
		return false
	}
	if self.par == nil {
		return true
	}
	return self.par.CheckPar(isF)
}

func (self *Layer)Check_1(){
	if self.tem != nil {
		//if (self.par.direction>0) != self.tem.Dis  {
		if (len(self.par.cans) - self.par.splitID >2) {
			self.checkTem()
		}
		return
	}
	if self.par == nil {
		return
	}
	isF := self.direction<0
	if !self.par.CheckPar(isF){
		return
	}
	fmt.Printf("%d %10.2f %5d %5d %10.2f\r\n",self.tag,self.par.direction,len(self.par.cans),len(self.cans),self.direction)
	self.getTemplate(isF)
}

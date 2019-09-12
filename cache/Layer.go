package cache
import(
	"math"
)

type Layer struct{
	cans []Element
	direction float64
	par *Layer
	child *Layer
}

func (self *Layer) initAdd (c Element){

	le := len(self.cans)
	self.cans = append(self.cans,c)
	var maxD,absMaxD float64
	var sum  = c.Diff()
	var splitID int
	for i,_c := range self.cans[:le] {
		sum += _c.Diff()
		d := c.Val() - _c.Val()
		absD := math.Abs(d)
		if (absD > absMaxD){
			maxD = d
			absMaxD = absD
			splitID = i
		}
	}
	sum /= float64(len(self.cans))
	if sum > absMaxD {
		return
	}
	self.direction = maxD
	self.cans = self.cans[splitID:]

}

func (self *Layer) add(c Element){

	if self.direction == 0 {
		self.initAdd(c)
		return
	}
	le := len(self.cans)
	self.cans = append(self.cans,c)
	var absMaxD float64
	var sum  = c.Diff()
	var splitID int
	for i,_c := range self.cans[:le] {
		sum += _c.Diff()
		d := c.Val() - _c.Val()
		absD := math.Abs(d)
		if (absD > absMaxD){
			//maxD = d
			absMaxD = absD
			if ((d>0) != (self.direction>0)) {
				splitID = i
			}
		}
	}
	sum /= float64(len(self.cans))
	if sum > absMaxD {
		return
	}
	if self.par == nil {
		self.par = &Layer{}
	}
	self.par.add(NewNode(self.cans[:splitID+1]))
	self.cans = self.cans[splitID:]

}

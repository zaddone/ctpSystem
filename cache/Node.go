package cache
import(
	//"fmt"
)
type Element interface{
	Diff() float64
	Val() float64
	Time() int64
	LastTime() int64
	Dur() int64
}
type Node struct {
	Eles []Element
	Diff_ float64
	Val_ float64
	Time_ int64
	LastTime_ int64
	Dur_ int64
}
func NewNode(eles []Element) (n *Node) {
	b:=eles[0]
	e:=eles[len(eles)-1]
	n  = &Node{
		Eles:eles,
		Diff_:e.Val() - b.Val(),
		Time_ :b.Time(),
		LastTime_:e.LastTime(),
	}
	for _,e := range eles{
		n.Val_ += e.Val()
	}
	n.Val_ /= float64(len(eles))
	n.Dur_ = n.LastTime_ - n.Time_
	return n
}

func (self *Node)Val() float64{
	return self.Val_
}
func (self *Node)Time() int64{
	return self.Time_
}
func (self *Node)LastTime() int64{
	return self.LastTime_
}
func (self *Node)Dur() int64{
	return self.Dur_
}

func (self *Node)Diff() float64 {
	return self.Diff_
}
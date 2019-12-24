package cache
import(
	//"io"
	//"fmt"
)
type Element interface{
	Diff() float64
	SetDiff(float64)
	Val() float64
	Time() int64
	LastTime() int64
	Dur() int64
	Each(func(Element)error) error
	Name() string
	Max() float64
	Min() float64
	SetDur(int64)
	GetEle(int) Element
}
type Node struct {
	Eles []Element
	Diff_ float64
	Val_ float64
	Time_ int64
	LastTime_ int64
	Dur_ int64
}
func MergeElement(a,b Element) (c *Node) {

	var eles []Element
	a.Each(func(c_ Element)error{
		eles = append(eles,c_)
		return nil
	})
	b.Each(func(c_ Element)error{
		eles = append(eles,c_)
		return nil
	})
	return NewNode(eles)

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
		n.Val_ += e.Val()*float64(e.Dur())
		n.Dur_ += e.Dur()
	}

	n.Val_ /= float64(n.Dur())
	//fmt.Println(n.Val_)
	//fmt.Println(n.Dur_,n.Val_)
	//n.Dur_ = n.LastTime_ - n.Time_
	return n
}
func (self *Node) GetEle(i int)Element {
	return self.Eles[i]
}

func (self *Node) SetDur(d int64) {
	self.Dur_ = d
}
func (self *Node) Min() float64{
	if self.Diff_>0 {
		return self.Eles[0].Min()
	}else{
		return self.Eles[len(self.Eles)-1].Min()
	}
}
func (self *Node) Max() float64{
	if self.Diff_<0 {
		return self.Eles[0].Max()
	}else{
		return self.Eles[len(self.Eles)-1].Max()
	}
}
func (self *Node) Name() (n string){
	return self.Eles[0].Name()
	//self.Each(func(e Element)error{
	//	n = e.Name()
	//	return io.EOF
	//})
	//return
}

func (self *Node) Each(fn func(Element)error)error{

	//fn(self)
	for _,e := range self.Eles {
		err := e.Each(fn)
		if err != nil {
			return err
		}
	}
	return nil

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
func (self *Node) SetDiff(d float64){
	self.Diff_ = d
}
func (self *Node)Diff() float64 {
	return self.Diff_
}

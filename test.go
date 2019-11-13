package main
import(
	//"os"
	//"os/exec"
	"fmt"
	//"bufio"
	//"time"
)
type InsOrder struct {
	State int
	par *InsOrder
}
func (self *InsOrder)cop(){
	self.par = &(*self)
	self.par.State = self.State+1
}

func main(){
	O := &InsOrder{}
	O.cop()
	fmt.Println(O.State,O.par.State)
}

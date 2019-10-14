package main
import(
	"fmt"
	"strings"
)
func main(){
	str := "ins abc adfa sfd"
	d := strings.SplitN(str," ",2)
	for _,s := range d{
		fmt.Println(s)
	}
}

package main
import(
	//"os"
	//"os/exec"
	"fmt"
	//"bufio"
	"time"
)
func main(){
	str := fmt.Sprintf("%d",time.Now().UnixNano())[:13]
	fmt.Println(str,len(str))
}

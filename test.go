package main
import(
	//"os"
	"os/exec"
	"fmt"
	"bufio"
)
func main(){
	cmd := exec.Command(
		"bin/mdServer",
		"9999",
		"150797",
		"abc2019",
		"Dimon2019",
		"tcp://218.202.237.33:10112",
		"/tmp/md")
	out,err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(out)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	for {
		buf.ReadLine()
		l,p,e := buf.ReadLine()
		fmt.Println(string(l),p,e)
		if e != nil {
			panic(e)
		}
	}
	//cmd.Wait()
	//fmt.Println(out.String())
	//var line [8196]byte
	//var n int
	//for{
	//	n,err = out.Read(line[:])
	//	fmt.Println(string(line[:n]))
	//	if err != nil {
	//		fmt.Println(err)
	//		//panic(err)
	//	}
	//}
}

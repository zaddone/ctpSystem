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
	//buf := bufio.NewReader(out)
	//var out bytes.Buffer
	//cmd.Stdout = &out
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	var line [8196]byte
	for {
		//buf.ReadLine()
		n,e := out.Read(line[:])
		//l,e := buf.ReadBytes('\n')
		//l,p,e := buf.ReadLine()
		if e != nil {
			panic(e)
		}
		fmt.Println(string(l[:len(l)-1]),e)
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

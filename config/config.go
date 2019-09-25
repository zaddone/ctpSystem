package config
import(
	//"net/http"
	"github.com/BurntSushi/toml"
	"flag"
	"os"
)
var(
	LogFileName   = flag.String("c", "conf.log", "config log")
	Conf *Config
)
func init(){
	//EntryList = make(chan *Entry,1000)
	//flag.Parse()
	Conf = NewConfig(*LogFileName)
}
type Config struct {
	Port string
	Static string
	Templates string
	RunServer bool

	Weight int

	BrokerID string
	UserID string
	Password string
	PasswordBak string
	Taddr []string
	Maddr []string
	DefAdd int

}
func (self *Config)Save(fileName string){
	fi,err := os.OpenFile(fileName,os.O_CREATE|os.O_WRONLY,0777)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	e := toml.NewEncoder(fi)
	err = e.Encode(self)
	if err != nil {
		panic(err)
	}
}
func NewConfig(fileName string)  *Config {
	var c Config
	_,err := os.Stat(fileName)
	if err != nil {
		c.Port=":8080"
		c.Static = "static"
		c.Templates = "./templates/*"
		c.BrokerID = "9999"
		c.UserID = "150797"
		c.Password = "Dimon2019"
		c.PasswordBak = "abc2019"
		c.Weight = 4
		c.DefAdd = 2
		c.RunServer = false
		c.Taddr = []string{
			"tcp://180.168.146.187:10100",
			"tcp://180.168.146.187:10101",
			"tcp://218.202.237.33:10102",
			"tcp://180.168.146.187:10130",
		}
		c.Maddr = []string{
			"tcp://180.168.146.187:10110",
			"tcp://180.168.146.187:10111",
			"tcp://218.202.237.33:10112",
			"tcp://180.168.146.187:10131",
		}
		c.Save(fileName)
	}else{
		if _,err := toml.DecodeFile(fileName,&c);err != nil {
			panic(err)
		}
	}
	return &c
}

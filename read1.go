package main

import(
	"github.com/zaddone/ctpSystem/config"
	"github.com/zaddone/ctpSystem/cache"
	//"github.com/boltdb/bolt"
	"path/filepath"
	"encoding/binary"
	"os"
	"fmt"
	"time"

)
func loadCache(){
	filepath.Walk(config.Conf.DbPath,func(name string,f os.FileInfo,er error)error{
		if f.IsDir(){
			return nil
		}
		ca := cache.StoreCache(map[string]string{"InstrumentID":f.Name()})
		fmt.Println(name)
		runRead(ca)
		return nil
	})
}

func runRead(ca *cache.Cache) {
	//db = ca.DB
	t,err := ca.DB.Begin(false)
	if err != nil {
		panic(err)
	}
	b := t.Bucket([]byte(ca.Info["InstrumentID"]))
	if b == nil {
		return
		panic(fmt.Errorf("b == nil"))
	}
	err = b.ForEach(func(k,v []byte)error{
		c := cache.NewCandle(
			ca.Info["InstrumentID"],
			int64(binary.BigEndian.Uint64(k)),
			v)
		//fmt.Println(c)
		if time.Unix(c.Time(),0).Month()<time.December{
			return nil
		}
		if ca.L != nil{
			ca.L.Add(c)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}

func main(){
	loadCache()
}

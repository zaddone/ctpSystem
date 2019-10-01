package main

import(
	"fmt"
	"flag"
	"github.com/boltdb/bolt"
	"github.com/zaddone/ctpSystem/cache"
	"encoding/binary"
)
var (
	dbName = flag.String("db","ins.db","db name")
	DB *bolt.DB
	Cache = cache.NewCache()
)
func init(){
	flag.Parse()
	var err error
	DB,err = bolt.Open(*dbName,0600,nil)
	if err != nil {
		panic(err)
	}
}

func main(){

	fmt.Println("start")
	t,err := DB.Begin(false)
	if err != nil {
		panic(err)
	}
	err = t.ForEach(func(name []byte,b *bolt.Bucket)error{

		if b== nil {
			return nil
		}
		cache.Count = [3]float64{0,0,0}
		err := b.ForEach(func(k,v []byte)error{
			c := cache.NewCandle(
				string(name),
				int64(binary.BigEndian.Uint64(k)),
				v)
			if c.Ask == c.Bid {
				b.Delete(k)
				return nil
			}
			//fmt.Println(c)
			Cache.Add(c)
			return nil
		})
		if err != nil {
			return err
		}
		if cache.Count[1]+cache.Count[2]>0 {
			fmt.Println(string(name),cache.Count,cache.Count[1]/cache.Count[2])
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}
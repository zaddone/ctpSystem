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
	//Cache = cache.NewCache()
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
		for i,c := range cache.Count[:4]{
			for j,_c := range c {
				cache.Count[4][j]+=_c
			}
			cache.Count[i] = [4]float64{0,0,0,0}
		}
		ca := cache.StoreCache(map[string]string{"InstrumentID":string(name)})
		if ca.L == nil {
			return nil
		}
		//fmt.Println(string(name))
		err := b.ForEach(func(k,v []byte)error{
			c := cache.NewCandle(
				string(name),
				int64(binary.BigEndian.Uint64(k)),
				v)
			//if c.Ask == c.Bid {
			//	b.Delete(k)
			//	return nil
			//}
			//fmt.Println(c)
			//Cache.Add(c)
			//cache.AddCandle(c)
			ca.L.Add(c)
			return nil
		})
		if err != nil {
			return err
		}
		var c0,c1,c2,c3 float64
		for _,c := range cache.Count[:4]{
			c0+=c[0]
			c1+=c[1]
			c2+=c[2]
			c3+=c[3]
			//sum += c[0]
			//if (c[1]+c[2]) >0 {
			//	fmt.Println(string(name),i,c,c[1]/c[2])
			//}
		}
		c_ := cache.Count[4]
		if c1+c2 >0 {
			fmt.Println(string(name),c0,c1,c2,c3,c_[0]/c_[1],c_[2]/c_[3])
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

}

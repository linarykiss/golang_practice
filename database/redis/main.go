package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	ExampleClient()
}

func ExampleClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, "err: ", err)

	err = client.Set("feekey", "examples", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("feekey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("feekey", val)

	val2, err := client.Get("feekey2").Result()
	if err == redis.Nil {
		fmt.Println("feekey2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("feekey2", val2)
	}

	// 不存在时设置
	set, err := client.SetNX("feekey", "value", 10*time.Second).Result()
	fmt.Println(set, err)

	//直接执行一个命令
	res, err := client.Do("set", "dotest", "testdo").Result()
	fmt.Println(res, err)

	// 在末尾添加
	res2, err := client.Append("feekey", "_add").Result()
	fmt.Println(res2, err)

	val, err = client.Get("feekey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("feekey", val)

	Get := func(redisdb *redis.Client, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd("get", key)
		redisdb.Process(cmd)
		return cmd
	}

	v, err := Get(client, "key_does_not_exist").Result()
	fmt.Printf("%q %s", v, err)

}

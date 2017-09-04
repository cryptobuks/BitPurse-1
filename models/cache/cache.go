package cache

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

// {"key":"collectionName","conn":":6039","dbNum":"0","password":"thePassWord"}
func New(_key string) IMyCache {
	config := fmt.Sprintf(`{"conn": "127.0.0.1:6379", "key": "%s"}`, _key)

	c, err := cache.NewCache("my", config)
	if err != nil {
		beego.Error(err)
		return nil
	}
	return c.(IMyCache)
}

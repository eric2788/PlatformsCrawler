package crawling

import "fmt"

type Cache struct {
	Module string
	Prefix string
	Local map[string]interface{}
}

func NewCache(module, prefix string) *Cache {
	return &Cache{
		Module: module,
		Prefix: prefix,
		Local: make(map[string]interface{}),
	}
}

func (c *Cache) key(id string) string {
	return fmt.Sprintf("%s:%s:%s", c.Module, c.Prefix, id)
}

func (c *Cache) SetStruct(id string, arg interface{}) error {
	err := Store(c.key(id), arg)
	if err != nil {
		logger.Errorf("儲存 %s 到 redis 時出現錯誤: %v, 將嘗試使用本地快取", c.Prefix, err)
		c.Local[id] = arg
	}
	return err
}

func (c *Cache) GetStruct(id string, res interface{}) (bool, error) {
	return GetStruct(c.key(id), res)
}

func (c *Cache) GetString(id string) (string, bool) {
	result, err := GetString(c.key(id))
	if err != nil {
		logger.Errorf("從 redis 獲取 %s 時出現錯誤: %v, 將嘗試使用本地快取", c.Prefix, err)
		if str, ok := c.Local[id]; ok {
			return str.(string), ok
		}
	}
	return result, result != ""
}

func (c *Cache) SetString(id, value string) error {
	err := SetString(c.key(id), value)
	if err != nil {
		logger.Errorf("儲存 %s 到 redis 時出現錯誤: %v, 將嘗試使用本地快取", c.Prefix, err)
		c.Local[id] = value
	}
	return err
}
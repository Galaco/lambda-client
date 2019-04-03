package cache

import (
	"github.com/galaco/Lambda-Core/core/entity"
	"github.com/galaco/Lambda-Core/core/model"
	entity2 "github.com/galaco/Lambda-Core/game/entity"
)

type PropCache struct {
	cacheList []Entry
}

func (c *PropCache) NeedsRecache() bool {
	return len(c.cacheList) == 0
}

func (c *PropCache) Add(props ...entity.IEntity) {
	for _, prop := range props {
		c.cacheList = append(c.cacheList, Entry{
			Transform: prop.Transform(),
			Model:     prop.(entity2.IProp).Model(),
		})
	}
}

func (c *PropCache) All() *[]Entry {
	return &c.cacheList
}

type Entry struct {
	Transform *entity.Transform
	Model     *model.Model
}

func NewPropCache(props ...entity.IEntity) *PropCache {
	c := &PropCache{}

	if len(props) > 0 {
		c.Add(props...)
	}

	return c
}

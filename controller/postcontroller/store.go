package postcontroller

import "github.com/L-oris/yabb/models/post"

func (c Controller) addPost(p post.Post) postControllerStore {
	c.store[p.ID] = p
	return c.store
}

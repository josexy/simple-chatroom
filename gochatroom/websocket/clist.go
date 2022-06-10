package websocket

type ClientList []*wsClient

func newClientList() *ClientList {
	return &ClientList{}
}

func (list *ClientList) PushBack(client *wsClient) {
	*list = append(*list, client)
}

func (list *ClientList) Remove(client *wsClient) bool {
	index := list.Index(client)
	if index == -1 {
		return false
	}
	*list = append((*list)[:index], (*list)[index+1:]...)
	return true
}

func (list *ClientList) Contain(client *wsClient) bool {
	return list.Index(client) != -1
}

func (list *ClientList) Index(client *wsClient) int {
	for index, c := range *list {
		if c == client {
			return index
		}
	}
	return -1
}

func (list *ClientList) Range(fn func(client *wsClient)) {
	for _, client := range *list {
		if fn != nil {
			fn(client)
		}
	}
}

func (list *ClientList) Size() int {
	return len(*list)
}

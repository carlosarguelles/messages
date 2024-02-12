package internal

type Pool struct {
	clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	unregister chan *Client
}

func NewPool() *Pool {
	return &Pool{
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (p *Pool) Run() {
	for {
		select {
		case client := <-p.Register:
			p.clients[client] = true
		case client := <-p.unregister:
			if _, ok := p.clients[client]; ok {
				delete(p.clients, client)
				close(client.Send)
			}
		case message := <-p.Broadcast:
			for client := range p.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(p.clients, client)
				}
			}
		}
	}
}

func (p *Pool) Close() {
	close(p.Register)
	close(p.Broadcast)
	close(p.unregister)
}

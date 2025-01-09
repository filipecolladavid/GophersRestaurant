package restaurant

import "fmt"

type Stats struct {
	CustomersArrived int
	OrdersAccepted   int
	OrdersPrepared   int
	OrdersDelivered  int
	OrdersCancelled  int
	CustomersLeft    int
	CustomersServed  int
}

type Event struct {
	ID   int
	Name string
}

type Logger struct {
	Events chan Event
	Stats  Stats
}

func NewLogger() *Logger {
	return &Logger{
		Events: make(chan Event),
		Stats:  Stats{},
	}
}

func (logger *Logger) StreamEvent(event Event) {
	logger.Events <- event
}

func (logger *Logger) OpenRestaurantLog() {
	Event := Event{
		ID:   0,
		Name: "[Restaurant] - Opened",
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) CloseRestaurantLog() {
	Event := Event{
		ID:   1,
		Name: "[Restaurant] - Closed",
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) NewChefLog(chef *Chef) {
	Event := Event{
		ID:   2,
		Name: fmt.Sprintf("[Chef] - %s - Recruited", chef.Name),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) NewWaiterLog(waiter *Waiter) {
	Event := Event{
		ID:   3,
		Name: fmt.Sprintf("[Waiter] - %s - Recruited", waiter.Name),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) NewCustomerLog(customer *Customer) {
	Event := Event{
		ID:   4,
		Name: fmt.Sprintf("[Customer] - %d - Arrived", customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ChefAvailableLog(chef *Chef) {
	Event := Event{
		ID:   5,
		Name: fmt.Sprintf("[Chef] - %s - Available", chef.Name),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ChefAcceptedCustomerLog(chef *Chef, customer *Customer) {
	Event := Event{
		ID:   6,
		Name: fmt.Sprintf("[Chef] - %s - Accepted Customer %d", chef.Name, customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ChefCancelledCustomerLeftLog(chef *Chef, customer *Customer) {
	Event := Event{
		ID:   7,
		Name: fmt.Sprintf("[Chef] - %s - Cancelled Order for Customer %d", chef.Name, customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ChefFinishedOrderLog(chef *Chef, customer *Customer) {
	Event := Event{
		ID:   8,
		Name: fmt.Sprintf("[Chef] - %s - Finished Order for Customer %d", chef.Name, customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ChefPreparedButCustomerLeftLog(chef *Chef, customer *Customer) {
	Event := Event{
		ID:   9,
		Name: fmt.Sprintf("[Chef] - %s - Prepared Order for Customer %d, but Customer %d left", chef.Name, customer.ID, customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) CustomerWaitingLeftLog(order *Order) {
	Event := Event{
		ID:   10,
		Name: fmt.Sprintf("[Customer] - %d - Left after Waiting for Order %s", order.Customer.ID, order.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) WaiterAvailableLog(waiter *Waiter) {
	Event := Event{
		ID:   11,
		Name: fmt.Sprintf("[Waiter] - %s - Available", waiter.Name),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) WaiterWasToDeliverButCustomerLeftLog(waiter *Waiter, order *Order) {
	Event := Event{
		ID:   12,
		Name: fmt.Sprintf("[Waiter] - %s - Was to Deliver Order %s, but Customer left", waiter.Name, order.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) WaiterServedOrderLog(waiter *Waiter, order *Order) {
	Event := Event{
		ID:   13,
		Name: fmt.Sprintf("[Waiter] - %s - Served Order for Customer %d", waiter.Name, order.Customer.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) CustomerAcceptedOrderLog(order *Order) {
	Event := Event{
		ID:   14,
		Name: fmt.Sprintf("[Customer] - %d - Accepted Order %s", order.Customer.ID, order.ID),
	}
	fmt.Println(Event.Name)
	go logger.StreamEvent(Event)
}

func (logger *Logger) ResetLogger() {
	logger.Stats = Stats{}
}

func (logger *Logger) CreateStats() {
	for event := range logger.Events {
		if event.ID == 4 {
			logger.Stats.CustomersArrived++
		} else if event.ID == 6 {
			logger.Stats.OrdersAccepted++
		} else if event.ID == 8 {
			logger.Stats.OrdersPrepared++
		} else if event.ID == 13 {
			logger.Stats.OrdersDelivered++
		} else if event.ID == 7 || event.ID == 9 || event.ID == 10 || event.ID == 12 {
			logger.Stats.CustomersLeft++
		} else if event.ID == 14 {
			logger.Stats.CustomersServed++
		}
	}
}

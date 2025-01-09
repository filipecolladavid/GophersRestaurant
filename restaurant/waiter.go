package restaurant

import (
	"context"
	"math/rand"
)

// Struct to represent the Waiter Entity
type Waiter struct {
	ID    int
	Name  string
	Order *Order
}

// Responsible for creating a limited number of Waiters (Worker pool pattern)
func GenerateWaiters(waiterChannel chan<- *Waiter, amountOfWaiters int) {
	for i := 0; i < amountOfWaiters; i++ {
		newWaiter := &Waiter{
			ID:   i,
			Name: WaitersNames[rand.Intn(len(WaitersNames))],
		}

		logger.NewWaiterLog(newWaiter)
		waiterChannel <- newWaiter
	}
}

// Responsible for handling the delivery of an order to a customer.
//
// Handles restaurant closing and customer cancelation before the order is delivered.
func (waiter *Waiter) ServeOrder(restCtx context.Context, customerCtx context.Context, order *Order, waiterChannel chan *Waiter) {
	select {
	case <-restCtx.Done():
		return
	case <-customerCtx.Done():
		waiter.Order = nil
		waiterChannel <- waiter
		logger.WaiterWasToDeliverButCustomerLeftLog(waiter, order)
		return

	default:
		logger.WaiterServedOrderLog(waiter, order)
		customerCtx.Done()
		waiterChannel <- waiter
		waiter.Order = nil
		logger.CustomerAcceptedOrderLog(order)
		return
	}
}

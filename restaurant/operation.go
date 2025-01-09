package restaurant

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var logger = NewLogger()

// Initiates the Restuarant operation.
//
// Sets up the customer waiting line, the chef and waiter worker pools.
//
// In te end prints statistics of the operation.
func OpenRestaurant(operationTimeout time.Duration) {
	logger.OpenRestaurantLog()

	loggerCtx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	restCtx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	customerChannel := make(chan *Customer)
	chefChannel := make(chan *Chef, 5)
	waiterChannel := make(chan *Waiter, 5)

	GenerateChefs(chefChannel, 5)
	GenerateWaiters(waiterChannel, 5)

	go GenerateCustomers(restCtx, customerChannel)
	go StartOperation(restCtx, customerChannel, chefChannel, waiterChannel, -1)
	go logger.CreateStats()

	<-restCtx.Done()
	<-loggerCtx.Done()

	fmt.Printf("\nRestaurant Work-Day Stats:\n")
	fmt.Printf("Customers Arrived: %d\n", logger.Stats.CustomersArrived)
	fmt.Printf("Orders Accepted: %d\n", logger.Stats.OrdersAccepted)
	fmt.Printf("Orders Prepared: %d\n", logger.Stats.OrdersPrepared)
	fmt.Printf("Orders Delivered: %d\n", logger.Stats.OrdersDelivered)
	fmt.Printf("Customers Left: %d\n", logger.Stats.CustomersLeft)
	fmt.Printf("Customers Served: %d\n", logger.Stats.CustomersServed)

	logger.CloseRestaurantLog()
}

// Starts the main operation with the configurations OpenRestaurant.
//
// Receives a time_to_prepare variable to improve unit testing - if is -1 then it's random, otherwise its the value passed
func StartOperation(restCtx context.Context, customerChannel <-chan *Customer, chefChannel chan *Chef, waiterChannel chan *Waiter, time_to_prepare int) {
	for {
		select {
		case <-restCtx.Done():
			return
		case newCustomer := <-customerChannel:
			select {
			case <-restCtx.Done():
				return
			case availableChef := <-chefChannel:
				customerCtx, cancel := context.WithTimeout(context.Background(), time.Duration(newCustomer.TimeToWait)*time.Second)
				defer cancel()
				time.Sleep(1 * time.Second) // Time it takes to process the request
				logger.ChefAvailableLog(availableChef)
				order_idx := rand.Intn(len(OrdersNames))
				if time_to_prepare == -1 {
					time_to_prepare = rand.Intn(5) + 1
				}
				newCustomer.Order = &Order{
					ID:            "CustID_" + strconv.Itoa(newCustomer.ID) + "_" + OrdersNames[order_idx],
					Name:          OrdersNames[order_idx],
					TimeToPrepare: time_to_prepare,
					Customer:      newCustomer,
				}

				availableChef.ActiveOrder = newCustomer.Order
				logger.ChefAcceptedCustomerLog(availableChef, newCustomer)

				go availableChef.PrepareOrder(restCtx, customerCtx, chefChannel, waiterChannel)
			}
		}
	}
}

// This function is called when the chef finishes preparing an order.
//
// Waits for a waiter to pick up and deliver to the customer.
//
// Handles restaurant closing and customer cancelation before the order is picked up by the waiter.
func DeliveryOperation(restCtx context.Context, customerCtx context.Context, order *Order, waiterChannel chan *Waiter) {
	select {
	case <-restCtx.Done():
		return
	case <-customerCtx.Done():
		logger.CustomerWaitingLeftLog(order)
		return
	case availableWaiter := <-waiterChannel:
		availableWaiter.Order = order
		logger.WaiterAvailableLog(availableWaiter)
		go availableWaiter.ServeOrder(restCtx, customerCtx, order, waiterChannel)
	}

}

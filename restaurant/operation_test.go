package restaurant

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func mockGenerateCustomers(restCtx context.Context, customerChannel chan<- *Customer, number_customers int, timeToWait []int) {
	for i := range number_customers {
		select {
		case <-restCtx.Done():
			return
		default:
			newCustomer := &Customer{
				ID:         i,
				TimeToWait: timeToWait[i],
			}
			logger.NewCustomerLog(newCustomer)
			customerChannel <- newCustomer
		}
	}
}

func mockOperation(number_chef int, number_waiters int, timeTowait []int, wg *sync.WaitGroup, time_to_prepare int, operation_time int) {
	number_customers := len(timeTowait)
	logger.ResetLogger()
	fmt.Printf("\nRestaurant Work-Day Stats:\n")
	fmt.Printf("Customers Arrived: %d\n", logger.Stats.CustomersArrived)
	fmt.Printf("Orders Accepted: %d\n", logger.Stats.OrdersAccepted)
	fmt.Printf("Orders Prepared: %d\n", logger.Stats.OrdersPrepared)
	fmt.Printf("Orders Delivered: %d\n", logger.Stats.OrdersDelivered)
	fmt.Printf("Customers Left: %d\n", logger.Stats.CustomersLeft)
	fmt.Printf("Customers Served: %d\n", logger.Stats.CustomersServed)

	loggerCtx, cancelLogger := context.WithTimeout(context.Background(), time.Duration(operation_time)*time.Second)
	defer cancelLogger()
	restCtx, cancelRest := context.WithTimeout(context.Background(), time.Duration(operation_time)*time.Second)
	defer cancelRest()

	customerChannel := make(chan *Customer)
	chefChannel := make(chan *Chef, number_chef)
	waiterChannel := make(chan *Waiter, number_waiters)

	GenerateChefs(chefChannel, 5)
	GenerateWaiters(waiterChannel, 5)
	go mockGenerateCustomers(restCtx, customerChannel, number_customers, timeTowait)
	go StartOperation(restCtx, customerChannel, chefChannel, waiterChannel, time_to_prepare)
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
	wg.Done()
}

// Tests normal operation - where no client cancells
func TestNormalOperation(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	timeToWait := []int{60, 60, 60, 60, 60}
	go mockOperation(5, 5, timeToWait, &wg, -1, 20)
	wg.Wait()

	if logger.Stats.CustomersArrived != 5 {
		t.Errorf("Customers Arrived: Got %d, wanted %d", logger.Stats.CustomersArrived, 5)
	}
	if logger.Stats.CustomersServed != 5 {
		t.Errorf("Customers Served: Got %d, wanted %d", logger.Stats.CustomersServed, 5)
	}
	if logger.Stats.CustomersLeft != 0 {
		t.Errorf("Customers Left: Got %d, wanted %d", logger.Stats.CustomersLeft, 0)
	}
	if logger.Stats.OrdersAccepted != 5 {
		t.Errorf("Orders Accepted: Got %d, wanted %d", logger.Stats.OrdersAccepted, 5)
	}
	if logger.Stats.OrdersPrepared != 5 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersPrepared, 5)
	}
	if logger.Stats.OrdersDelivered != 5 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersDelivered, 5)
	}

	logger.CloseRestaurantLog()
}

// Test all clients leaving before getting served
func TestImpatientClients(t *testing.T) {
	wg := sync.WaitGroup{}
	timeToWait := []int{0, 0, 0, 0, 0}
	wg.Add(1)
	go mockOperation(5, 5, timeToWait, &wg, -1, 20)
	wg.Wait()

	if logger.Stats.CustomersArrived != 5 {
		t.Errorf("Customers Arrived: Got %d, wanted %d", logger.Stats.CustomersArrived, 5)
	}
	if logger.Stats.CustomersServed != 0 {
		t.Errorf("Customers Served: Got %d, wanted %d", logger.Stats.CustomersServed, 0)
	}
	if logger.Stats.CustomersLeft != 5 {
		t.Errorf("Customers Left: Got %d, wanted %d", logger.Stats.CustomersLeft, 5)
	}
	if logger.Stats.OrdersAccepted != 5 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersAccepted, 5)
	}
	if logger.Stats.OrdersPrepared != 0 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersPrepared, 0)
	}
	if logger.Stats.OrdersDelivered != 0 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersDelivered, 0)
	}
	logger.CloseRestaurantLog()
}

// Test 1 client leaving before getting served
func TestOneImpatientClients(t *testing.T) {
	wg := sync.WaitGroup{}
	timeToWait := []int{60, 60, 0, 60, 60}
	wg.Add(1)
	go mockOperation(5, 5, timeToWait, &wg, -1, 20)
	wg.Wait()

	if logger.Stats.CustomersArrived != 5 {
		t.Errorf("Customers Arrived: Got %d, wanted %d", logger.Stats.CustomersArrived, 5)
	}
	if logger.Stats.CustomersServed != 4 {
		t.Errorf("Customers Served: Got %d, wanted %d", logger.Stats.CustomersServed, 4)
	}
	if logger.Stats.CustomersLeft != 1 {
		t.Errorf("Customers Left: Got %d, wanted %d", logger.Stats.CustomersLeft, 1)
	}
	if logger.Stats.OrdersAccepted != 5 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersAccepted, 5)
	}
	if logger.Stats.OrdersPrepared != 4 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersPrepared, 4)
	}
	if logger.Stats.OrdersDelivered != 4 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersDelivered, 4)
	}

	logger.CloseRestaurantLog()
}

// Test more clients than workers
func TestMoreCustomersThanWorkers(t *testing.T) {
	wg := sync.WaitGroup{}
	timeToWait := []int{60, 60, 60, 60, 60, 60, 60, 60, 60}
	wg.Add(1)
	go mockOperation(5, 5, timeToWait, &wg, -1, 20)
	wg.Wait()

	if logger.Stats.CustomersArrived != 9 {
		t.Errorf("Customers Arrived: Got %d, wanted %d", logger.Stats.CustomersArrived, 9)
	}
	if logger.Stats.CustomersServed != 9 {
		t.Errorf("Customers Served: Got %d, wanted %d", logger.Stats.CustomersServed, 9)
	}
	if logger.Stats.CustomersLeft != 0 {
		t.Errorf("Customers Left: Got %d, wanted %d", logger.Stats.CustomersLeft, 0)
	}
	if logger.Stats.OrdersAccepted != 9 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersAccepted, 9)
	}
	if logger.Stats.OrdersPrepared != 9 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersPrepared, 9)
	}
	if logger.Stats.OrdersDelivered != 9 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersDelivered, 9)
	}

	logger.CloseRestaurantLog()
}

// Test orders droped after being Accepted
func TestLongPrepareOrders(t *testing.T) {
	wg := sync.WaitGroup{}
	timeToWait := []int{10, 10, 10, 10, 10, 10, 10, 10, 10}
	wg.Add(1)
	go mockOperation(5, 5, timeToWait, &wg, 15, 40)
	wg.Wait()

	if logger.Stats.CustomersArrived != 9 {
		t.Errorf("Customers Arrived: Got %d, wanted %d", logger.Stats.CustomersArrived, 9)
	}
	if logger.Stats.CustomersServed != 0 {
		t.Errorf("Customers Served: Got %d, wanted %d", logger.Stats.CustomersServed, 0)
	}
	if logger.Stats.CustomersLeft != 9 {
		t.Errorf("Customers Left: Got %d, wanted %d", logger.Stats.CustomersLeft, 9)
	}
	if logger.Stats.OrdersAccepted != 9 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersAccepted, 9)
	}
	if logger.Stats.OrdersPrepared != 0 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersPrepared, 0)
	}
	if logger.Stats.OrdersDelivered != 0 {
		t.Errorf("Orders Prepared: Got %d, wanted %d", logger.Stats.OrdersDelivered, 0)
	}
	logger.CloseRestaurantLog()
}

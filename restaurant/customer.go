package restaurant

import (
	"context"
	"math/rand"
	"time"
)

// Struct to represent the Customer Entity
type Customer struct {
	ID         int
	TimeToWait int
	Order      *Order
}

// Stream that creates customers and places them in line (channel).
func GenerateCustomers(restCtx context.Context, customerChannel chan<- *Customer) {
	customerIdCounter := 0

	for {
		select {
		case <-restCtx.Done():
			return
		default:
			newCustomer := &Customer{
				ID:         customerIdCounter,
				TimeToWait: int(time.Duration(rand.Intn(5)*int(time.Second) + 1)),
			}
			logger.NewCustomerLog(newCustomer)
			customerChannel <- newCustomer
		}
		customerIdCounter++
	}
}

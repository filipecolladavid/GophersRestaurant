package restaurant

import (
	"context"
	"math/rand"
	"time"
)

// Struct to represent the Chef Entity
type Chef struct {
	ID          int
	Name        string
	ActiveOrder *Order
}

// Responsible for creating a limited number of chefs (Worker pool pattern)
func GenerateChefs(chefChannel chan<- *Chef, amountOfChefs int) {
	for i := 0; i < amountOfChefs; i++ {
		newChef := &Chef{
			ID:   i,
			Name: ChefNames[rand.Intn(len(ChefNames))],
		}

		logger.NewChefLog(newChef)
		chefChannel <- newChef
	}
}

// Responsible for handling the preparation of an order.
//
// Handles restaurant closing and customer cancelation before and after the order is prepared.
//
// Once the order is prepared, is placed on operation.DeliveryOperation for a waiter to pickup.
func (chef *Chef) PrepareOrder(restCtx context.Context, customerCtx context.Context, chefChannel chan<- *Chef, waiterChannel chan *Waiter) {
	select {
	case <-restCtx.Done():
		return
	case <-customerCtx.Done():
		logger.ChefCancelledCustomerLeftLog(chef, chef.ActiveOrder.Customer)
		chef.ActiveOrder = nil
		chefChannel <- chef
		return
	case <-time.After(time.Duration(chef.ActiveOrder.TimeToPrepare) * time.Second):
		select {
		case <-restCtx.Done():
			return
		case <-customerCtx.Done():
			logger.ChefPreparedButCustomerLeftLog(chef, chef.ActiveOrder.Customer)
			chef.ActiveOrder = nil
			chefChannel <- chef
			return
		default:
			logger.ChefFinishedOrderLog(chef, chef.ActiveOrder.Customer)
			go DeliveryOperation(restCtx, customerCtx, chef.ActiveOrder, waiterChannel)
			chef.ActiveOrder = nil
			chefChannel <- chef
		}
	}
}

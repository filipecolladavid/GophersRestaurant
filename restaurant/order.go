package restaurant

// Struct to represent the Order entity
type Order struct {
	ID            string
	Name          string
	TimeToPrepare int
	Customer      *Customer
}

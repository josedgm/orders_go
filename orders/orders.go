package orders

import (
	"fmt"
)

// Incoming order item structure
type IncomingOrderItem struct {
	ItemId  string  `json:"itemId"`
	CostEur float64 `json:"costEur"`
}

// Incoming order structure (as received in batch)
type IncomingOrder struct {
	CustomerId string              `json:"customerId"`
	OrderId    string              `json:"orderId"`
	Timestamp  int64               `json:"timestamp"`
	Items      []IncomingOrderItem `json:"items"`
}

// Output item structure for customer items list
type CustomerItem struct {
	CustomerId string  `json:"customerId"`
	ItemId     string  `json:"itemId"`
	CostEur    float64 `json:"costEur"`
}

// Output per customer
type CustomerItems struct {
	Items []CustomerItem `json:"items"`
}

// ValidateIncomingOrder checks the required keys and data validity
func ValidateIncomingOrder(order IncomingOrder) error {
	if order.CustomerId == "" || order.OrderId == "" || order.Timestamp == 0 || len(order.Items) == 0 {
		return fmt.Errorf("order %s is missing required order fields", order.OrderId) // What if the orderId is also missing?
	}
	for _, item := range order.Items {
		if item.ItemId == "" {
			return fmt.Errorf("at least one item missing its itemId in order %s", order.OrderId)
		}
		if item.CostEur < 0 {
			return fmt.Errorf("itemId %s has negative costEur", item.ItemId)
		}
	}
	return nil
}

// Check if item is already in CustomerItem list
func containsItem(items []CustomerItem, itemId string) bool {
	for _, item := range items {
		if item.ItemId == itemId {
			return true
		}
	}
	return false
}

// TransformOrders processes a batch of orders and returns grouped customer items
func TransformOrders(orders []IncomingOrder) ([]CustomerItems, []error) {
	customerMap := make(map[string][]CustomerItem)
	var errs []error

	for _, order := range orders {
		if err := ValidateIncomingOrder(order); err != nil {
			// Collect validation error but continue processing other orders
			errs = append(errs, fmt.Errorf("order %s: %w", order.OrderId, err))
			continue
		}
		for _, item := range order.Items {
			if !containsItem(customerMap[order.CustomerId], item.ItemId) {
				customerMap[order.CustomerId] = append(customerMap[order.CustomerId], CustomerItem{ // only append if the item is not already in the CustomerItems slice.
					CustomerId: order.CustomerId,
					ItemId:     item.ItemId,
					CostEur:    item.CostEur,
				})
			}
		}
	}

	var result []CustomerItems
	for _, items := range customerMap {
		result = append(result, CustomerItems{Items: items})
	}

	return result, errs
}

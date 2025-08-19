package orders

import (
	"reflect"
	"testing"
)

func TestValidateIncomingOrder(t *testing.T) {
	tests := []struct {
		name    string
		order   IncomingOrder
		wantErr bool
	}{
		{
			name: "valid order",
			order: IncomingOrder{
				CustomerId: "cust1",
				OrderId:    "ord1",
				Timestamp:  1234567890,
				Items: []IncomingOrderItem{
					{ItemId: "item1", CostEur: 10},
				},
			},
			wantErr: false,
		},
		{
			name: "missing fields",
			order: IncomingOrder{
				CustomerId: "",
				OrderId:    "",
				Timestamp:  0,
				Items:      nil,
			},
			wantErr: true,
		},
		{
			name: "item missing ItemId",
			order: IncomingOrder{
				CustomerId: "cust1",
				OrderId:    "ord1",
				Timestamp:  1234,
				Items: []IncomingOrderItem{
					{ItemId: "", CostEur: 5},
				},
			},
			wantErr: true,
		},
		{
			name: "item negative cost",
			order: IncomingOrder{
				CustomerId: "cust1",
				OrderId:    "ord1",
				Timestamp:  1234,
				Items: []IncomingOrderItem{
					{ItemId: "item1", CostEur: -1},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIncomingOrder(tt.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIncomingOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainsItem(t *testing.T) {
	items := []CustomerItem{
		{CustomerId: "c1", ItemId: "i1", CostEur: 10},
		{CustomerId: "c1", ItemId: "i2", CostEur: 20},
	}

	tests := []struct {
		itemId string
		want   bool
	}{
		{"i1", true},
		{"i2", true},
		{"i3", false},
	}

	for _, tt := range tests {
		if got := containsItem(items, tt.itemId); got != tt.want {
			t.Errorf("containsItem() = %v, want %v for itemId %s", got, tt.want, tt.itemId)
		}
	}
}

func TestTransformOrders(t *testing.T) {
	orders := []IncomingOrder{
		{
			CustomerId: "cust1",
			OrderId:    "ord1",
			Timestamp:  1000,
			Items: []IncomingOrderItem{
				{ItemId: "item1", CostEur: 10},
				{ItemId: "item2", CostEur: 20},
				{ItemId: "item1", CostEur: 10}, // duplicate item to test deduplication
			},
		},
		{
			CustomerId: "cust2",
			OrderId:    "ord2",
			Timestamp:  1001,
			Items: []IncomingOrderItem{
				{ItemId: "item3", CostEur: 30},
			},
		},
		{
			// invalid order to test error accumulation
			CustomerId: "",
			OrderId:    "ord3",
			Timestamp:  0,
			Items:      nil,
		},
	}

	wantResult := []CustomerItems{
		{
			Items: []CustomerItem{
				{CustomerId: "cust1", ItemId: "item1", CostEur: 10},
				{CustomerId: "cust1", ItemId: "item2", CostEur: 20},
			},
		},
		{
			Items: []CustomerItem{
				{CustomerId: "cust2", ItemId: "item3", CostEur: 30},
			},
		},
	}

	gotResult, errs := TransformOrders(orders)

	if len(errs) == 0 {
		t.Errorf("TransformOrders() expected errors for invalid orders, got none")
	}

	// Compare the results ignoring order in array (simple check)
	if len(gotResult) != len(wantResult) {
		t.Fatalf("TransformOrders() got %d customer groups, want %d", len(gotResult), len(wantResult))
	}

	for _, wantGroup := range wantResult {
		found := false
		for _, gotGroup := range gotResult {
			if reflect.DeepEqual(wantGroup, gotGroup) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("TransformOrders() missing expected group: %+v", wantGroup)
		}
	}
}

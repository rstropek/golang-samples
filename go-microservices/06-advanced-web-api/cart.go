package main

import (
	"regexp"
	"sort"
)

type Product struct {
	ID           int    `json:"id"`
	Description  string `json:"desc"`
	PricePerItem uint32 `json:"pricePerItem"`
}

type CartItem struct {
	ID       int     `json:"id"`
	Product  Product `json:"prod"`
	Quantity uint16  `json:"quantity"`
}

type ReadCartOp struct {
	Resp chan []CartItem
}

type UpdateQuantityOp struct {
	ID       int
	Quantity uint16
	Resp     chan bool
}

type ApplyCouponOp struct {
	Coupon string
	Resp   chan bool
}

var products = []Product{
	{ID: 1, Description: "Bike", PricePerItem: 499},
	{ID: 2, Description: "Tire", PricePerItem: 19},
	{ID: 3, Description: "Sport Shoes", PricePerItem: 129},
	{ID: 4, Description: "Cap", PricePerItem: 13},
	{ID: 5, Description: "Tools", PricePerItem: 99},
}

func generateCart() map[int]CartItem {
	cart := make(map[int]CartItem, len(products))
	for i, prod := range products {
		p := Product{ID: prod.ID, Description: prod.Description, PricePerItem: prod.PricePerItem}
		cart[i] = CartItem{ID: i, Product: p, Quantity: 1}
	}

	return cart
}

func cloneCart(source map[int]CartItem) []CartItem {
	ids := make([]int, len(source))
	ix := 0
	for k := range source {
		ids[ix] = k
		ix++
	}

	sort.Ints(ids)

	copy := make([]CartItem, len(source))
	ix = 0
	for _, k := range ids {
		copy[ix] = source[k]
		ix++
	}

	return copy
}

type CartChannels struct {
	reads   chan ReadCartOp
	updates chan UpdateQuantityOp
	resets  chan ReadCartOp
	coupons chan ApplyCouponOp
	quit    chan bool
}

func CreateCartChannels() CartChannels {
	return CartChannels{
		reads:   make(chan ReadCartOp),
		updates: make(chan UpdateQuantityOp),
		resets:  make(chan ReadCartOp),
		coupons: make(chan ApplyCouponOp),
		quit:    make(chan bool),
	}
}

func Cart(channels CartChannels) {
	couponRegex, err := regexp.Compile("^[A-Z]{3}-[0-9]{4}$")
	if err != nil {
		panic(err)
	}

	cartState := generateCart()
	for {
		select {
		case read := <-channels.reads:
			read.Resp <- cloneCart(cartState)
		case update := <-channels.updates:
			if item, ok := cartState[update.ID]; ok {
				if update.Quantity == 0 {
					delete(cartState, update.ID)
				} else {
					item.Quantity = update.Quantity
					cartState[update.ID] = item
				}
				update.Resp <- true
			} else {
				update.Resp <- false
			}
		case reset := <-channels.resets:
			cartState = generateCart()
			reset.Resp <- cloneCart(cartState)
		case coupon := <-channels.coupons:
			if !couponRegex.MatchString(coupon.Coupon) {
				coupon.Resp <- false
			} else {
				for k, v := range cartState {
					v.Product.PricePerItem = v.Product.PricePerItem * 9 / 10
					cartState[k] = v
				}

				coupon.Resp <- true
			}
		case <-channels.quit:
			return
		}
	}
}

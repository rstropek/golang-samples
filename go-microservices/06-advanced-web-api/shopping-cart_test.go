package main

import (
	"testing"
)

func TestFilledAfterCreating(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	read := ReadCartOp{Resp: make(chan []CartItem)}
	channels.reads <- read
	cart := <-read.Resp

	if len(cart) == 0 {
		t.Error("Response is empty")
	}

	channels.quit <- true
}

func TestUpdateQuantity(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	update := UpdateQuantityOp{ID: 0, Quantity: 99, Resp: make(chan bool)}
	channels.updates <- update
	success := <-update.Resp

	if !success {
		t.Error("ID wasn't found")
	}

	read := ReadCartOp{Resp: make(chan []CartItem)}
	channels.reads <- read
	cart := <-read.Resp
	found := false
	for _, item := range cart {
		if item.ID == 0 && item.Quantity == 99 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Quantity wasn't updated")
	}

	channels.quit <- true
}

func TestSetToZero(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	read := ReadCartOp{Resp: make(chan []CartItem)}
	channels.reads <- read
	lenBeforeUpdate := len(<-read.Resp)

	update := UpdateQuantityOp{ID: 1, Quantity: 0, Resp: make(chan bool)}
	channels.updates <- update
	<-update.Resp

	channels.reads <- read
	if lenBeforeUpdate != len(<-read.Resp)+1 {
		t.Error("Item should be gone")
	}

	channels.quit <- true
}
func TestFailingUpdateQuantity(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	update := UpdateQuantityOp{ID: 99, Quantity: 99, Resp: make(chan bool)}
	channels.updates <- update
	success := <-update.Resp

	if success {
		t.Error("ID shouldn't exist")
	}

	channels.quit <- true
}

func TestReset(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	update := UpdateQuantityOp{ID: 1, Quantity: 99, Resp: make(chan bool)}
	channels.updates <- update
	<-update.Resp

	reset := ReadCartOp{Resp: make(chan []CartItem)}
	channels.resets <- reset
	cart := <-reset.Resp
	if cart[1].Quantity != 1 {
		t.Error("Quantity wasn't reset")
	}

	channels.quit <- true
}

func TestCoupon(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	read := ReadCartOp{Resp: make(chan []CartItem)}
	channels.reads <- read
	cartOriginal := <-read.Resp

	coupon := ApplyCouponOp{Coupon: "ABC-1234", Resp: make(chan bool)}
	channels.coupons <- coupon
	success := <-coupon.Resp

	if !success {
		t.Error("Error while applying coupon")
	}

	read = ReadCartOp{Resp: make(chan []CartItem)}
	channels.reads <- read
	newCart := <-read.Resp
	if cartOriginal[0].Product.PricePerItem*9/10 != newCart[0].Product.PricePerItem {
		t.Error("Price wasn't changed")
	}

	channels.quit <- true
}

func TestInvalidCoupon(t *testing.T) {
	channels := CreateCartChannels()
	go Cart(channels)

	coupon := ApplyCouponOp{Coupon: "AB-1234", Resp: make(chan bool)}
	channels.coupons <- coupon
	success := <-coupon.Resp

	if success {
		t.Error("Coupon should be invalid")
	}

	channels.quit <- true
}

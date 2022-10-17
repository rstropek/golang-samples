package main

import "fmt"

var cartChannels = CreateCartChannels()

func main() {
	go Cart(cartChannels)

	fmt.Println("Reading 📖")
	cart := readItems()
	for _, item := range cart {
		fmt.Printf("%d items of %s are in the cart for %d\n", item.Quantity, item.Product.Description, item.Product.PricePerItem)
	}

	fmt.Println("\nWriting 📝")
	update := UpdateQuantityOp{ID: cart[0].ID, Quantity: cart[0].Quantity * 2, Resp: make(chan bool)}
	cartChannels.updates <- update
	success := <-update.Resp
	if success {
		cart = readItems()
		fmt.Printf("Now %d items of %s are in the cart\n", cart[0].Quantity, cart[0].Product.Description)
	}
	
	fmt.Println("\nCoupon 🎫")
	coupon := ApplyCouponOp{Coupon: "ABC-1234", Resp: make(chan bool)}
	cartChannels.coupons <- coupon
	success = <-coupon.Resp
	if (success) {
		cart = readItems()
		for _, item := range cart {
			fmt.Printf("%d items of %s are in the cart for %d\n", item.Quantity, item.Product.Description, item.Product.PricePerItem)
		}
	}
}

func readItems() []CartItem {
	// Read all items from the cart
	read := ReadCartOp{Resp: make(chan []CartItem)}
	cartChannels.reads <- read
	cart := <-read.Resp

	return cart
}

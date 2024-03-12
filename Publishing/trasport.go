package Publishing

import (
	"ProductionOrders/order"
	"context"
	"fmt"
	"sync"
)

func SearchChanges(ctx context.Context, s Service, e ServiceElastic, delta string) error {
	var wg sync.WaitGroup
	added := 0
	updated := 0

	// Split the string into two parts
	hrs := delta[:8] // First 8 characters *date
	tms := delta[8:] // Remaining characters *hours

	data, err := s.GetOrders(ctx, hrs, tms)

	if err != nil {
		fmt.Println("Error ", err)
		return err
	}

	for _, orden := range data {
		wg.Add(1)

		go func(orden order.Orden) {
			defer wg.Done()
			ord, err := e.FindOne(ctx, orden.Order)

			if err != nil {
				fmt.Println(err)
				err = e.Insert(ctx, orden)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Orden " + orden.Order + " agregado")
					added++
				}

			}
			if ord.Order == orden.Order {
				err = e.Update(ctx, orden)
				if err != nil {
					fmt.Println(err)

				} else {
					fmt.Println("Orden " + orden.Order + " actualizado")
					updated++
				}

			}
		}(orden)
	}
	wg.Wait()

	fmt.Println("Ordenes Añadidad: ", added)
	fmt.Println("Ordenes Actualizadas: ", updated)

	return nil
}

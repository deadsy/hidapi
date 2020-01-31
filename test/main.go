//-----------------------------------------------------------------------------
/*

Test program to exercise the libhidapi Go wrapper.

*/
//-----------------------------------------------------------------------------

package main

import (
	"fmt"
	"os"

	"github.com/deadsy/hidapi"
)

//-----------------------------------------------------------------------------

func libTest() error {

	err := hidapi.Init()
	if err != nil {
		return err
	}
	defer hidapi.Exit()

	devs := hidapi.Enumerate(0, 0)
	for i, d := range devs {
		fmt.Printf("Device %d\n", i)
		fmt.Printf("%s\n\n", d)
	}

	for _, di := range devs {

		dev, err := hidapi.Open(di.VendorID, di.ProductID, "")
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}

		fmt.Printf("%s\n", dev)

		dev.Close()
	}

	return nil
}

//-----------------------------------------------------------------------------

func main() {

	err := libTest()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

//-----------------------------------------------------------------------------

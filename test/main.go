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

	devs := hidapi.Enumerate(0, 0)

	for i, d := range devs {
		fmt.Printf("Device %d\n", i)
		fmt.Printf("%s\n", d)
	}

	err = hidapi.Exit()
	if err != nil {
		return err
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

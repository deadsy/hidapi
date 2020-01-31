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

	//devs := hidapi.Enumerate(0, 0)
	//for i, d := range devs {
	//	fmt.Printf("Device %d\n", i)
	//	fmt.Printf("%s\n", d)
	//}

	// Open the device using the VID, PID,
	// and optionally the Serial number.

	vid := uint16(0x046d)
	pid := uint16(0xc52b)

	dev := hidapi.Open(vid, pid, "")

	s, err := dev.GetManufacturerString()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)

	s, err = dev.GetProductString()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)

	s, err = dev.GetSerialNumberString()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", s)

	dev.Close()

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

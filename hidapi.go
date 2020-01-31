//-----------------------------------------------------------------------------
/*

Go bindings for the libhidapi library.

See: https://github.com/deadsy/hidapi
See: https://github.com/signal11/hidapi

*/
//-----------------------------------------------------------------------------

package hidapi

/*
#cgo pkg-config: hidapi-libusb
//#cgo pkg-config: hidapi-hidraw
#include <hidapi/hidapi.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"strings"
	"unsafe"

	wchar "github.com/vitaminwater/cgo.wchar"
)

//-----------------------------------------------------------------------------
// utility functions

// go2cBuffer creates a C uint8_t buffer from a Go []byte buffer.
// Call freeBuffer on the returned C buffer.
func go2cBuffer(buf []byte) *C.uint8_t {
	return (*C.uint8_t)(unsafe.Pointer(C.CString(string(buf))))
}

// c2goCopy copies a C uint8_t buffer into a Go []byte slice.
func c2goCopy(s []byte, buf *C.uint8_t) []byte {
	x := (*[1 << 30]byte)(unsafe.Pointer(buf))
	copy(s, x[:])
	return s
}

// c2goSlice creates a Go []byte slice and copies in a C uint8_t buffer.
func c2goSlice(buf *C.uint8_t, n int) []byte {
	s := make([]byte, n)
	return c2goCopy(s, buf)
}

// allocBuffer allocates a C uint8_t buffer of length n bytes.
// Call freeBuffer on the returned C buffer.
func allocBuffer(n int) *C.uint8_t {
	return (*C.uint8_t)(C.malloc(C.size_t(n)))
}

// freeBuffer frees a C uint8_t buffer.
func freeBuffer(buf *C.uint8_t) {
	C.free(unsafe.Pointer(buf))
}

// Set a value within a C uint8_t buffer.
func (buf *C.uint8_t) Set(idx int, val byte) {
	x := (*[1 << 30]byte)(unsafe.Pointer(buf))
	x[idx] = val
}

// boolToInt converts a boolean to an int.
func boolToInt(x bool) int {
	if x {
		return 1
	}
	return 0
}

//-----------------------------------------------------------------------------

// Device is a HID device.
type Device struct {
	dev *C.struct_hid_device_
}

func (d *Device) String() string {
	s := []string{}
	x, err := d.GetManufacturerString()
	if err == nil {
		s = append(s, fmt.Sprintf("manufacturer: %s", x))
	}
	x, err = d.GetProductString()
	if err == nil {
		s = append(s, fmt.Sprintf("product: %s", x))
	}
	x, err = d.GetSerialNumberString()
	if err == nil {
		s = append(s, fmt.Sprintf("serial number: %s", x))
	}
	return strings.Join(s, "\n")
}

//-----------------------------------------------------------------------------

// DeviceInfo stores information about the HID device.
type DeviceInfo struct {
	Path            string // Platform-specific device path
	VendorID        uint16 // Vendor ID
	ProductID       uint16 // Product ID
	SerialNumber    string // Serial Number
	ReleaseNumber   uint16 // Device Release Number in BCD
	Manufacturer    string // Manufacturer String
	Product         string // Product String
	UsagePage       uint16 // Usage Page for this Device/Interface (Windows/Mac only)
	Usage           uint16 // Usage for this Device/Interface (Windows/Mac only)
	InterfaceNumber int    // USB interface
}

func (di *DeviceInfo) String() string {
	s := []string{}
	s = append(s, fmt.Sprintf("path: %s", di.Path))
	s = append(s, fmt.Sprintf("vid.pid: %04x:%04x", di.VendorID, di.ProductID))
	s = append(s, fmt.Sprintf("serial number: %s", di.SerialNumber))
	s = append(s, fmt.Sprintf("release number: %d.%d", (di.ReleaseNumber>>4)&0xf, di.ReleaseNumber&0xf))
	s = append(s, fmt.Sprintf("manufacturer: %s", di.Manufacturer))
	s = append(s, fmt.Sprintf("product: %s", di.Product))
	//s = append(s, fmt.Sprintf("usage page: %04x", di.UsagePage))
	//s = append(s, fmt.Sprintf("usage: %04x", di.Usage))
	//s = append(s, fmt.Sprintf("interface number: %d", di.InterfaceNumber))
	return strings.Join(s, "\n")
}

const wsLength = 128

//-----------------------------------------------------------------------------
// Errors

// Error stores C-API error information.
type Error struct {
	FunctionName string // function name
	ErrorString  string // device error string
	ReturnCode   int    // return code
}

func (e *Error) Error() string {
	if e.ErrorString != "" {
		return fmt.Sprintf("%s() returned %d: %s", e.FunctionName, e.ReturnCode, e.ErrorString)
	}
	return fmt.Sprintf("%s() returned %d", e.FunctionName, e.ReturnCode)
}

// apiError returns a C-API error.
func apiError(name string, rc int) *Error {
	return &Error{
		FunctionName: name,
		ReturnCode:   rc,
	}
}

// devError returns a C-API error from a device call.
func (d *Device) devError(name string, rc int) error {
	s, err := d.ErrorString()
	if err != nil {
		s = "?"
	}
	return &Error{
		FunctionName: name,
		ErrorString:  s,
		ReturnCode:   rc,
	}
}

// ErrorString returns a string describing the last error which occurred.
func (d *Device) ErrorString() (string, error) {
	ws := wchar.FromWcharStringPtr(unsafe.Pointer(C.hid_error(d.dev)))
	return ws.GoString()
}

//-----------------------------------------------------------------------------

// Init initializes the HIDAPI library.
func Init() error {
	rc := int(C.hid_init())
	if rc != 0 {
		return apiError("hid_init", rc)
	}
	return nil
}

// Exit finalizes the HIDAPI library.
func Exit() error {
	rc := int(C.hid_exit())
	if rc != 0 {
		return apiError("hid_exit", rc)
	}
	return nil
}

// Enumerate returns a list of HID Device Information.
func Enumerate(vid, pid uint16) []*DeviceInfo {
	diList := C.hid_enumerate(C.uint16_t(vid), C.uint16_t(pid))
	if diList == nil {
		return nil
	}
	devs := []*DeviceInfo{}
	di := diList
	for di != nil {
		sString, _ := wchar.WcharStringPtrToGoString(unsafe.Pointer(di.serial_number))
		mString, _ := wchar.WcharStringPtrToGoString(unsafe.Pointer(di.manufacturer_string))
		pString, _ := wchar.WcharStringPtrToGoString(unsafe.Pointer(di.product_string))
		dev := &DeviceInfo{
			Path:            C.GoString(di.path),
			VendorID:        uint16(di.vendor_id),
			ProductID:       uint16(di.product_id),
			SerialNumber:    sString,
			ReleaseNumber:   uint16(di.release_number),
			Manufacturer:    mString,
			Product:         pString,
			UsagePage:       uint16(di.usage_page),
			Usage:           uint16(di.usage),
			InterfaceNumber: int(di.interface_number),
		}
		devs = append(devs, dev)
		di = di.next
	}
	C.hid_free_enumeration(diList)
	return devs
}

// Open a HID device using a vendor ID (VID), product ID (PID) and optionally a serial number.
func Open(vid, pid uint16, sn string) (*Device, error) {
	var dev *C.struct_hid_device_
	if sn == "" {
		dev = C.hid_open(C.uint16_t(vid), C.uint16_t(pid), nil)
	} else {
		ws, err := wchar.FromGoString(sn)
		if err != nil {
			return nil, errors.New("can't convert serial number")
		}
		dev = C.hid_open(C.uint16_t(vid), C.uint16_t(pid), (*C.wchar_t)(ws.Pointer()))
	}
	if dev == nil {
		return nil, errors.New("hid_open() returned NULL")
	}
	d := &Device{
		dev: dev,
	}
	return d, nil
}

// OpenPath opens a HID device by its path name.
func OpenPath(path string) (*Device, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	dev := C.hid_open_path(cPath)
	if dev == nil {
		return nil, errors.New("hid_open_path() returned NULL")
	}
	d := &Device{
		dev: dev,
	}
	return d, nil
}

// Write an output report to a HID device.
func (d *Device) Write(data []byte) (int, error) {
	cBuffer := go2cBuffer(data)
	defer freeBuffer(cBuffer)
	rc := int(C.hid_write(d.dev, cBuffer, C.ulong(len(data))))
	if rc < 0 {
		return 0, d.devError("hid_write", rc)
	}
	return rc, nil
}

// ReadTimeout reads an input report from a HID device with timeout.
func (d *Device) ReadTimeout(id byte, length, milliseconds int) ([]byte, error) {
	cBuffer := allocBuffer(length)
	defer freeBuffer(cBuffer)
	cBuffer.Set(0, id)
	rc := int(C.hid_read_timeout(d.dev, cBuffer, C.ulong(length), C.int(milliseconds)))
	if rc < 0 {
		return nil, d.devError("hid_read_timeout", rc)
	}
	if rc == 0 {
		return nil, nil
	}
	return c2goSlice(cBuffer, rc), nil
}

// Read an input report from a HID device.
func (d *Device) Read(id byte, length int) ([]byte, error) {
	cBuffer := allocBuffer(length)
	defer freeBuffer(cBuffer)
	cBuffer.Set(0, id)
	rc := int(C.hid_read(d.dev, cBuffer, C.ulong(length)))
	if rc < 0 {
		return nil, d.devError("hid_read", rc)
	}
	if rc == 0 {
		return nil, nil
	}
	return c2goSlice(cBuffer, rc), nil
}

// SetNonBlocking sets the device handle to be non-blocking.
func (d *Device) SetNonBlocking(nonblock bool) error {
	rc := int(C.hid_set_nonblocking(d.dev, C.int(boolToInt(nonblock))))
	if rc != 0 {
		return d.devError("hid_set_nonblocking", rc)
	}
	return nil
}

// SendFeatureReport sends a feature report to the device.
func (d *Device) SendFeatureReport(data []byte) (int, error) {
	cBuffer := go2cBuffer(data)
	defer freeBuffer(cBuffer)
	rc := int(C.hid_send_feature_report(d.dev, cBuffer, C.ulong(len(data))))
	if rc < 0 {
		return 0, d.devError("hid_send_feature_report", rc)
	}
	return rc, nil
}

// GetFeatureReport gets a feature report from a HID device.
func (d *Device) GetFeatureReport(id byte, length int) ([]byte, error) {
	cBuffer := allocBuffer(length + 1)
	defer freeBuffer(cBuffer)
	cBuffer.Set(0, id)
	rc := int(C.hid_get_feature_report(d.dev, cBuffer, C.ulong(length+1)))
	if rc < 0 {
		return nil, d.devError("hid_get_feature_report", rc)
	}
	if rc == 0 {
		return nil, nil
	}
	return c2goSlice(cBuffer, rc), nil
}

// Close closes a HID device.
func (d *Device) Close() {
	C.hid_close(d.dev)
}

// GetManufacturerString returns the manufacturer string from a HID device.
func (d *Device) GetManufacturerString() (string, error) {
	ws := wchar.NewWcharString(wsLength)
	rc := int(C.hid_get_manufacturer_string(d.dev, (*C.wchar_t)(ws.Pointer()), wsLength))
	if rc != 0 {
		return "", d.devError("hid_get_manufacturer_string", rc)
	}
	return ws.GoString()
}

// GetProductString returns the product string from a HID device.
func (d *Device) GetProductString() (string, error) {
	ws := wchar.NewWcharString(wsLength)
	rc := int(C.hid_get_product_string(d.dev, (*C.wchar_t)(ws.Pointer()), wsLength))
	if rc != 0 {
		return "", d.devError("hid_get_product_string", rc)
	}
	return ws.GoString()
}

// GetSerialNumberString returns the serial number string from a HID device.
func (d *Device) GetSerialNumberString() (string, error) {
	ws := wchar.NewWcharString(wsLength)
	rc := int(C.hid_get_serial_number_string(d.dev, (*C.wchar_t)(ws.Pointer()), wsLength))
	if rc != 0 {
		return "", d.devError("hid_get_serial_number_string", rc)
	}
	return ws.GoString()
}

// GetIndexedString gets a string from a HID device, based on its string index.
func (d *Device) GetIndexedString(idx int) (string, error) {
	ws := wchar.NewWcharString(wsLength)
	rc := int(C.hid_get_indexed_string(d.dev, C.int(idx), (*C.wchar_t)(ws.Pointer()), wsLength))
	if rc != 0 {
		return "", d.devError("hid_get_indexed_string", rc)
	}
	return ws.GoString()
}

//-----------------------------------------------------------------------------

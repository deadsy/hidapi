//-----------------------------------------------------------------------------
/*

Go bindings for the libhidapi library.

See: https://github.com/deadsy/hidapi
See: https://github.com/signal11/hidapi

*/
//-----------------------------------------------------------------------------

package hidapi

/*
//#cgo pkg-config: libusb-1.0
#cgo pkg-config: hidapi-libusb
#include <hidapi/hidapi.h>
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"fmt"
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

type Device struct {
	dev *C.struct_hid_device_
}

type DeviceInfo struct {
	Path               string // Platform-specific device path
	Vid                uint16 // Vendor ID
	Pid                uint16 // Product ID
	SerialNumber       string // Serial Number
	ReleaseNumber      uint16 // Device Release Number in BCD
	ManufacturerString string // Manufacturer String
	ProductString      string // Product String
	UsagePage          uint16 // Usage Page for this Device/Interface (Windows/Mac only)
	Usage              uint16 // Usage for this Device/Interface (Windows/Mac only)
	InterfaceNumber    int    // USB interface
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
		return fmt.Sprintf("%s failed %s (%d)", e.FunctionName, e.ErrorString, e.ReturnCode)
	}
	return fmt.Sprintf("%s failed (%d)", e.FunctionName, e.ReturnCode)
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

/** @brief Enumerate the HID Devices.

	This function returns a linked list of all the HID devices
	attached to the system which match vendor_id and product_id.
	If @p vendor_id is set to 0 then any vendor matches.
	If @p product_id is set to 0 then any product matches.
	If @p vendor_id and @p product_id are both set to 0, then
	all HID devices will be returned.

	@ingroup API
	@param vendor_id The Vendor ID (VID) of the types of device
		to open.
	@param product_id The Product ID (PID) of the types of
		device to open.

    @returns
    	This function returns a pointer to a linked list of type
    	struct #hid_device, containing information about the HID devices
    	attached to the system, or NULL in the case of failure. Free
    	this linked list by calling hid_free_enumeration().
*/

func Enumerate(vid, pid uint16) []*DeviceInfo {
	//struct hid_device_info HID_API_EXPORT * HID_API_CALL hid_enumerate(unsigned short vendor_id, unsigned short product_id);
	return nil
}

/** @brief Free an enumeration Linked List

    This function frees a linked list created by hid_enumerate().

	@ingroup API
    @param devs Pointer to a list of struct_device returned from
    	      hid_enumerate().
*/

func FreeEnumeration(devs []*DeviceInfo) {
	//void  HID_API_EXPORT HID_API_CALL hid_free_enumeration(struct hid_device_info *devs);
}

/** @brief Open a HID device using a Vendor ID (VID), Product ID
(PID) and optionally a serial number.

If @p serial_number is NULL, the first device with the
specified VID and PID is opened.

@ingroup API
@param vendor_id The Vendor ID (VID) of the device to open.
@param product_id The Product ID (PID) of the device to open.
@param serial_number The Serial Number of the device to open
	               (Optionally NULL).

@returns
	This function returns a pointer to a #hid_device object on
	success or NULL on failure.
*/

func Open(vid, pid uint16, sn string) *Device {
	//HID_API_EXPORT hid_device * HID_API_CALL hid_open(unsigned short vendor_id, unsigned short product_id, const wchar_t *serial_number);
	return nil
}

/** @brief Open a HID device by its path name.

	The path name be determined by calling hid_enumerate(), or a
	platform-specific path name can be used (eg: /dev/hidraw0 on
	Linux).

	@ingroup API
    @param path The path name of the device to open

	@returns
		This function returns a pointer to a #hid_device object on
		success or NULL on failure.
*/

func OpenPath(path string) *Device {
	//HID_API_EXPORT hid_device * HID_API_CALL hid_open_path(const char *path);
	return nil
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

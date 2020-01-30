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
//#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	wchar "github.com/vitaminwater/cgo.wchar"
)

//-----------------------------------------------------------------------------
// utility functions

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

/** @brief Write an Output report to a HID device.

The first byte of @p data[] must contain the Report ID. For
devices which only support a single report, this must be set
to 0x0. The remaining bytes contain the report data. Since
the Report ID is mandatory, calls to hid_write() will always
contain one more byte than the report contains. For example,
if a hid report is 16 bytes long, 17 bytes must be passed to
hid_write(), the Report ID (or 0x0, for devices with a
single report), followed by the report data (16 bytes). In
this example, the length passed in would be 17.

hid_write() will send the data on the first OUT endpoint, if
one exists. If it does not, it will send the data through
the Control Endpoint (Endpoint 0).

@ingroup API
@param device A device handle returned from hid_open().
@param data The data to send, including the report number as
	the first byte.
@param length The length in bytes of the data to send.

@returns
	This function returns the actual number of bytes written and
	-1 on error.
*/

func (d *Device) Write(data []byte, length uint) (int, error) {
	//int  HID_API_EXPORT HID_API_CALL hid_write(hid_device *device, const unsigned char *data, size_t length);
	return 0, nil
}

/** @brief Read an Input report from a HID device with timeout.

Input reports are returned
to the host through the INTERRUPT IN endpoint. The first byte will
contain the Report number if the device uses numbered reports.

@ingroup API
@param device A device handle returned from hid_open().
@param data A buffer to put the read data into.
@param length The number of bytes to read. For devices with
	multiple reports, make sure to read an extra byte for
	the report number.
@param milliseconds timeout in milliseconds or -1 for blocking wait.

@returns
	This function returns the actual number of bytes read and
	-1 on error. If no packet was available to be read within
	the timeout period, this function returns 0.
*/

func (d *Device) ReadTimeout(length uint, ms int) ([]byte, error) {
	//int HID_API_EXPORT HID_API_CALL hid_read_timeout(hid_device *dev, unsigned char *data, size_t length, int milliseconds);
	return nil, nil
}

/** @brief Read an Input report from a HID device.

	Input reports are returned
    to the host through the INTERRUPT IN endpoint. The first byte will
	contain the Report number if the device uses numbered reports.

	@ingroup API
	@param device A device handle returned from hid_open().
	@param data A buffer to put the read data into.
	@param length The number of bytes to read. For devices with
		multiple reports, make sure to read an extra byte for
		the report number.

	@returns
		This function returns the actual number of bytes read and
		-1 on error. If no packet was available to be read and
		the handle is in non-blocking mode, this function returns 0.
*/

func (d *Device) Read(length uint) ([]byte, error) {
	//int  HID_API_EXPORT HID_API_CALL hid_read(hid_device *device, unsigned char *data, size_t length);
	return nil, nil
}

// SetNonBlocking sets the device handle to be non-blocking.
func (d *Device) SetNonBlocking(nonblock bool) error {
	rc := int(C.hid_set_nonblocking(d.dev, C.int(boolToInt(nonblock))))
	if rc != 0 {
		return d.devError("hid_set_nonblocking", rc)
	}
	return nil
}

/** @brief Send a Feature report to the device.

Feature reports are sent over the Control endpoint as a
Set_Report transfer.  The first byte of @p data[] must
contain the Report ID. For devices which only support a
single report, this must be set to 0x0. The remaining bytes
contain the report data. Since the Report ID is mandatory,
calls to hid_send_feature_report() will always contain one
more byte than the report contains. For example, if a hid
report is 16 bytes long, 17 bytes must be passed to
hid_send_feature_report(): the Report ID (or 0x0, for
devices which do not use numbered reports), followed by the
report data (16 bytes). In this example, the length passed
in would be 17.

@ingroup API
@param device A device handle returned from hid_open().
@param data The data to send, including the report number as
	the first byte.
@param length The length in bytes of the data to send, including
	the report number.

@returns
	This function returns the actual number of bytes written and
	-1 on error.
*/

func (d *Device) SendFeatureReport(data []byte) (int, error) {
	//int HID_API_EXPORT HID_API_CALL hid_send_feature_report(hid_device *device, const unsigned char *data, size_t length);
	return 0, nil
}

/** @brief Get a feature report from a HID device.

Set the first byte of @p data[] to the Report ID of the
report to be read.  Make sure to allow space for this
extra byte in @p data[]. Upon return, the first byte will
still contain the Report ID, and the report data will
start in data[1].

@ingroup API
@param device A device handle returned from hid_open().
@param data A buffer to put the read data into, including
	the Report ID. Set the first byte of @p data[] to the
	Report ID of the report to be read, or set it to zero
	if your device does not use numbered reports.
@param length The number of bytes to read, including an
	extra byte for the report ID. The buffer can be longer
	than the actual report.

@returns
	This function returns the number of bytes read plus
	one for the report ID (which is still in the first
	byte), or -1 on error.
*/

func (d *Device) GetFeatureReport(length uint) ([]byte, error) {
	//int HID_API_EXPORT HID_API_CALL hid_get_feature_report(hid_device *device, unsigned char *data, size_t length);
	return nil, nil
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

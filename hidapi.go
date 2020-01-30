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
import "fmt"

//-----------------------------------------------------------------------------
// Errors

// Error stores error information.
type Error struct {
	Name string // function name
	Code int    // C return code
}

// apiError returns an C-API error.
func apiError(name string, rc int) *Error {
	return &Error{
		Name: name,
		Code: rc,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s failed (%d)", e.Name, e.Code)
}

//-----------------------------------------------------------------------------

type Device struct {
	dev *C.struct_hid_device
}

type DeviceInfo struct {
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

/** @brief Set the device handle to be non-blocking.

In non-blocking mode calls to hid_read() will return
immediately with a value of 0 if there is no data to be
read. In blocking mode, hid_read() will wait (block) until
there is data to read before returning.

Nonblocking can be turned on and off at any time.

@ingroup API
@param device A device handle returned from hid_open().
@param nonblock enable or not the nonblocking reads
 - 1 to enable nonblocking
 - 0 to disable nonblocking.

@returns
	This function returns 0 on success and -1 on error.
*/

func (d *Device) SetNonBlocking(nonblock bool) error {
	//int  HID_API_EXPORT HID_API_CALL hid_set_nonblocking(hid_device *device, int nonblock);
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

/** @brief Close a HID device.

@ingroup API
@param device A device handle returned from hid_open().
*/

func (d *Device) Close() {
	//void HID_API_EXPORT HID_API_CALL hid_close(hid_device *device);
}

/** @brief Get The Manufacturer String from a HID device.

@ingroup API
@param device A device handle returned from hid_open().
@param string A wide string buffer to put the data into.
@param maxlen The length of the buffer in multiples of wchar_t.

@returns
	This function returns 0 on success and -1 on error.
*/

func (d *Device) GetManufacturerString() (string, error) {
	//int HID_API_EXPORT_CALL hid_get_manufacturer_string(hid_device *device, wchar_t *string, size_t maxlen);
	return "", nil
}

/** @brief Get The Product String from a HID device.

@ingroup API
@param device A device handle returned from hid_open().
@param string A wide string buffer to put the data into.
@param maxlen The length of the buffer in multiples of wchar_t.

@returns
	This function returns 0 on success and -1 on error.
*/

func (d *Device) GetProductString() (string, error) {
	//int HID_API_EXPORT_CALL hid_get_product_string(hid_device *device, wchar_t *string, size_t maxlen);
	return "", nil
}

/** @brief Get The Serial Number String from a HID device.

@ingroup API
@param device A device handle returned from hid_open().
@param string A wide string buffer to put the data into.
@param maxlen The length of the buffer in multiples of wchar_t.

@returns
	This function returns 0 on success and -1 on error.
*/

func (d *Device) GetSerialNumberString() (string, error) {
	//int HID_API_EXPORT_CALL hid_get_serial_number_string(hid_device *device, wchar_t *string, size_t maxlen);
	return "", nil
}

/** @brief Get a string from a HID device, based on its string index.

@ingroup API
@param device A device handle returned from hid_open().
@param string_index The index of the string to get.
@param string A wide string buffer to put the data into.
@param maxlen The length of the buffer in multiples of wchar_t.

@returns
	This function returns 0 on success and -1 on error.
*/

func (d *Device) GetIndexedString(idx int) (string, error) {
	//int HID_API_EXPORT_CALL hid_get_indexed_string(hid_device *device, int string_index, wchar_t *string, size_t maxlen);
	return "", nil
}

/** @brief Get a string describing the last error which occurred.

@ingroup API
@param device A device handle returned from hid_open().

@returns
	This function returns a string containing the last error
	which occurred or NULL if none has occurred.
*/

func (d *Device) Error() string {
	//    HID_API_EXPORT const wchar_t* HID_API_CALL hid_error(hid_device *device);
	return ""
}

//-----------------------------------------------------------------------------

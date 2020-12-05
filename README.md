[![Go Report Card](https://goreportcard.com/badge/github.com/deadsy/hidapi)](https://goreportcard.com/report/github.com/deadsy/hidapi)
[![GoDoc](https://godoc.org/github.com/deadsy/hidapi?status.svg)](https://godoc.org/github.com/deadsy/hidapi)

# hidapi
Go bindings for the hidapi library.

## What Is It?

hidapi is a C-based library providing an API for controlling USB HID class devices.
This package provides a Go wrapper for the C-library API so the library can be called from Go programs.

## Dependencies

 * hidapi (https://github.com/libusb/hidapi)
 * libusb-1.0 (https://libusb.info/)

## Notes

All C-API functions have Go wrappers.
The public interface of this package is a 1-1 mapping from the C-API to a Go style function prototypes.
 
## Status
 
 * Some testing has been done, mostly using USB based CMSIS-DAP devices.
 * Version 0.10.1 of the hidapi library has been tested.

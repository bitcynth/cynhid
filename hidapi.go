package cynhid

// #cgo pkg-config: hidapi-hidraw
// #include <hidapi/hidapi.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type hidapiDevice *C.hid_device
type hidapiDeviceInfo *C.struct_hid_device_info

func hidapiInit() error {
	return errFromErrno(C.hid_init())
}

func hidapiExit() {
	C.hid_exit()
}

func hidapiEnumerate(vendorID int, productID int) []hidapiDeviceInfo {
	var devs []hidapiDeviceInfo
	var dev *C.struct_hid_device_info
	dev = C.hid_enumerate(C.ushort(vendorID), C.ushort(productID))
	for dev != nil {
		devs = append(devs, dev)
		dev = dev.next
	}
	C.hid_free_enumeration(dev)
	return devs
}

func hidapiOpenPath(path string) (hidapiDevice, error) {
	var pathCStr *C.char
	pathCStr = C.CString(path)

	var dh *C.hid_device
	dh = C.hid_open_path(pathCStr)
	if dh == nil {
		return nil, fmt.Errorf("failed to open path")
	}

	C.free(unsafe.Pointer(pathCStr))

	return (hidapiDevice)(dh), nil
}

func hidapiClose(dev hidapiDevice) {
	C.hid_close((*C.hid_device)(dev))
}

func hidapiSetNonblocking(dev hidapiDevice, nonblock bool) {
	nb := 0
	if nonblock == true {
		nb = 1
	}
	C.hid_set_nonblocking((*C.hid_device)(dev), C.int(nb))
}

func hidapiWrite(dev hidapiDevice, data []byte, length int) (int, error) {
	var n C.int
	n = C.hid_write((*C.hid_device)(dev), (*C.uchar)(unsafe.Pointer(&data[0])), C.ulong(length))
	if n < 0 {
		return -1, fmt.Errorf("failed to write via hidapi")
	}
	return int(n), nil
}

func hidapiRead(dev hidapiDevice, length int) ([]byte, error) {
	var n C.int
	buf := make([]byte, length)
	n = C.hid_read((*C.hid_device)(dev), (*C.uchar)(unsafe.Pointer(&buf[0])), C.ulong(length))
	if n < 0 {
		return nil, fmt.Errorf("failed to read from hidapi")
	}
	return buf, nil
}

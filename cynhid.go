package cynhid

// #include <hidapi/hidapi.h>
import "C"
import "fmt"

// DeviceInfo represents *C.struct_hid_device_info in a more Go-friendly way
type DeviceInfo struct {
	Path               string
	VendorID           int
	ProductID          int
	SerialNumber       string
	ReleaseNumber      int
	ManufacturerString string
	ProductString      string
	UsagePage          int
	Usage              int
	InterfaceNumber    int
}

// Device represents an HID device
type Device struct {
	handle hidapiDevice
}

// Init initializes the HIDAPI library
func Init() error {
	return hidapiInit()
}

// Exit cleans up the library
func Exit() {
	hidapiExit()
}

// Enumerate generates a slice of devices with a VID, PID combo
func Enumerate(vendorID int, productID int) ([]DeviceInfo, error) {
	var output []DeviceInfo

	devs := hidapiEnumerate(vendorID, productID)
	if len(devs) < 1 {
		return nil, fmt.Errorf("failed to enumerate devices with VID: %x PID: %x", vendorID, productID)
	}
	for _, d := range devs {
		// TODO: populate SerialNumber, ManufacturerString, ProductString
		devInfo := DeviceInfo{
			Path:            C.GoString(d.path),
			VendorID:        int(d.vendor_id),
			ProductID:       int(d.product_id),
			ReleaseNumber:   int(d.release_number),
			UsagePage:       int(d.usage_page),
			Usage:           int(d.usage),
			InterfaceNumber: int(d.interface_number),
		}
		output = append(output, devInfo)
	}

	return output, nil
}

// OpenPath opens a hid device by using a path
func OpenPath(path string) (*Device, error) {
	dev, err := hidapiOpenPath(path)

	if err != nil {
		return nil, err
	}

	return &Device{handle: dev}, nil
}

// SetNonblocking sets a specific device to non-blocking mode
func (dev *Device) SetNonblocking(nonblock bool) {
	hidapiSetNonblocking(dev.handle, nonblock)
}

// Close cleans up a device handle
func (dev *Device) Close() {
	hidapiClose(dev.handle)
}

func (dev *Device) Write(data []byte, length int) (int, error) {
	return hidapiWrite(dev.handle, data, length)
}

func (dev *Device) Read(length int) ([]byte, error) {
	return hidapiRead(dev.handle, length)
}

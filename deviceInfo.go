package main

/*
#cgo LDFLAGS: -ljabra
#include <stdlib.h>
#include "jabra/Common.h"
*/
import "C"
import "log"

func NewDeviceInfo(cDeviceInfo C.Jabra_DeviceInfo) *DeviceInfo {
	deviceInfo := DeviceInfo{
		DeviceID:               (uint16)(cDeviceInfo.deviceID),
		ProductID:              (uint16)(cDeviceInfo.productID),
		VendorID:               (uint16)(cDeviceInfo.vendorID),
		DeviceName:             C.GoString(cDeviceInfo.deviceName),
		USBDevicePath:          C.GoString(cDeviceInfo.usbDevicePath),
		ParentInstanceID:       C.GoString(cDeviceInfo.parentInstanceId),
		IsDongle:               (bool)(cDeviceInfo.isDongle),
		DongleName:             C.GoString(cDeviceInfo.dongleName),
		Variant:                C.GoString(cDeviceInfo.variant),
		SerialNumber:           C.GoString(cDeviceInfo.serialNumber),
		IsInFirmwareUpdateMode: (bool)(cDeviceInfo.isInFirmwareUpdateMode),
		ParentDeviceId:         (uint16)(cDeviceInfo.parentDeviceId),
	}

	deviceInfo.IsBusylightSupported = (bool)(C.Jabra_IsBusylightSupported(cDeviceInfo.deviceID))
	deviceInfo.BusylightStatus = (bool)(C.Jabra_GetBusylightStatus(cDeviceInfo.deviceID))

	return &deviceInfo
}

type DeviceInfo struct {
	DeviceID               uint16 // unsigned short deviceID
	ProductID              uint16 // unsigned short productID
	VendorID               uint16 // unsigned short vendorID
	DeviceName             string // char* deviceName
	USBDevicePath          string // char* usbDevicePath
	ParentInstanceID       string // char* parentInstanceId
	IsDongle               bool   // bool isDongle
	DongleName             string // char* dongleName
	Variant                string // char* variant
	SerialNumber           string // char* serialNumber
	IsInFirmwareUpdateMode bool   // bool isInFirmwareUpdateMode
	ParentDeviceId         uint16 // unsigned short parentDeviceId
	IsBusylightSupported   bool
	BusylightStatus        bool
}

func (d *DeviceInfo) SetBusylightStatus(value bool) {
	if !d.IsBusylightSupported {
		return
	}
	if d.BusylightStatus == value {
		return
	}
	d.BusylightStatus = value
	C.Jabra_SetBusylightStatus(C.ushort(d.DeviceID), C.bool(value))
	log.Printf("Set busy light on %s to %t", d.DeviceName, d.BusylightStatus)
}

func (d *DeviceInfo) EnableBusylightStatus() {
	d.SetBusylightStatus(true)
}

func (d *DeviceInfo) DisableBusylightStatus() {
	d.SetBusylightStatus(false)
}

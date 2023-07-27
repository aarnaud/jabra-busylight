package main

/*
#cgo LDFLAGS: -ljabra
#include <stdlib.h>
#include "jabra/Common.h"
#include "jabra/JabraDeviceConfig.h"
#include "jabra/JabraNativeHid.h"
extern void goFirstscanfordevicesdonefunc(void);
extern void goDeviceattachedfunc(Jabra_DeviceInfo deviceInfo);
extern void goDeviceremovedfunc(unsigned short deviceID);
extern void goButtonindatarawhidfunc(unsigned short deviceID, unsigned short usagePage, unsigned short usage, unsigned char buttonInData);

__attribute__((weak))
char* testC(char* val) {
	return val;
};
*/
import "C"
import (
	"log"
	"sync"
	"time"
)

var deviceList = make(map[uint16]*DeviceInfo, 0)
var deviceListLock = sync.Mutex{}
var mainDevice *DeviceInfo

func main() {
	log.Println(C.GoString(C.testC(C.CString("testing C binding: this line must be print"))))

	C.Jabra_SetAppID(C.CString("linux-busylight"))
	init := C.Jabra_InitializeV2((*[0]byte)(C.goFirstscanfordevicesdonefunc), (*[0]byte)(C.goDeviceattachedfunc), (*[0]byte)(C.goDeviceremovedfunc), (*[0]byte)(C.goButtonindatarawhidfunc), (*[0]byte)(nil), true, nil)
	if !init {
		log.Fatalln("failed to init jabra SDK")
	}
	for {
		time.Sleep(60 * time.Second)
	}
}

//export goFirstscanfordevicesdonefunc
func goFirstscanfordevicesdonefunc() {
	log.Println("first scan")
}

//export goDeviceattachedfunc
func goDeviceattachedfunc(cdeviceinfo C.Jabra_DeviceInfo) {
	deviceInfo := NewDeviceInfo(cdeviceinfo)
	log.Println("attach ", deviceInfo.DeviceName)
	addDevice(deviceInfo)
}

//export goDeviceremovedfunc
func goDeviceremovedfunc(deviceid uint16) {
	if _, ok := deviceList[deviceid]; !ok {
		return
	}
	log.Println("remove ", deviceList[deviceid].DeviceName)
	removeDevice(deviceList[deviceid])
}

//export goButtonindatarawhidfunc
func goButtonindatarawhidfunc(deviceid uint16, usagepage uint16, usage uint16, buttonindata bool) {
	log.Println(deviceid, usagepage, usage, buttonindata)

	if mainDevice != nil && usagepage == 0xff30 && usage == 0x002a {
		if buttonindata {
			mainDevice.EnableBusylightStatus()
		} else {
			mainDevice.DisableBusylightStatus()
		}
	}
}

func addDevice(device *DeviceInfo) {
	deviceListLock.Lock()
	deviceList[device.DeviceID] = device
	if device.IsBusylightSupported {
		mainDevice = device
	}
	defer deviceListLock.Unlock()
}

func removeDevice(device *DeviceInfo) {
	deviceListLock.Lock()
	if mainDevice != nil && device.SerialNumber == mainDevice.SerialNumber {
		mainDevice = nil
	}
	delete(deviceList, device.DeviceID)
	defer deviceListLock.Unlock()
}

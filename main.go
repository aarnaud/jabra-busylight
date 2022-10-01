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

__attribute__((weak))
char* testC(char* val) {
	return val;
};
*/
import "C"
import (
	"github.com/jfreymuth/pulse"
	"github.com/jfreymuth/pulse/proto"
	"log"
	"sync"
	"time"
)

var deviceList = make(map[uint16]*DeviceInfo, 0)
var deviceListLock = sync.Mutex{}
var mainDevice *DeviceInfo

func main() {
	log.Println(C.GoString(C.testC(C.CString("testing C binding: this line must be print"))))
	pulseClient, err := pulse.NewClient()
	if err != nil {
		log.Fatalln("failed to connect to pulseaudio server")
	}

	C.Jabra_SetAppID(C.CString("linux-busylight"))
	init := C.Jabra_InitializeV2((*[0]byte)(C.goFirstscanfordevicesdonefunc), (*[0]byte)(C.goDeviceattachedfunc), (*[0]byte)(C.goDeviceremovedfunc), (*[0]byte)(nil), (*[0]byte)(nil), true, nil)
	if !init {
		log.Fatalln("failed to init jabra SDK")
	}

	for {
		time.Sleep(5 * time.Second)
		if mainDevice == nil {
			continue
		}
		if isSourceOutputPresent(pulseClient) {
			mainDevice.EnableBusylightStatus()
		} else {
			mainDevice.DisableBusylightStatus()
		}
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

func isSourceOutputPresent(client *pulse.Client) bool {
	sourceOutputReq := proto.GetSourceOutputInfoList{}
	sourceOutputReply := proto.GetSourceOutputInfoListReply{}
	err := client.RawRequest(&sourceOutputReq, &sourceOutputReply)
	if err != nil {
		log.Println("error to get source-outputs")
	}
	return len(sourceOutputReply) > 0
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

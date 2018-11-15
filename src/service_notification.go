package main

import (
	"github.com/NaySoftware/go-fcm"
	"sync"
)

type PushService struct {
	FCMClient *fcm.FcmClient
	sendLock  sync.Mutex
}

func (r *PushService) PushToAll(topic string, data map[string]string) (status *fcm.FcmResponseStatus, err error) {

	r.sendLock.Lock()
	defer r.sendLock.Unlock()

	r.FCMClient.NewFcmMsgTo(topic, data)
	status, err = r.FCMClient.Send()

	return
}

func (r *PushService) PushToDevice(deviceList []string, data map[string]string) (status *fcm.FcmResponseStatus, err error) {

	r.sendLock.Lock()
	defer r.sendLock.Unlock()

	id, err := newUUID()

	if err != nil {
		return nil, err
	}

	r.FCMClient.NewFcmRegIdsMsg([]string{id}, data)
	r.FCMClient.AppendDevices(deviceList)
	status, err = r.FCMClient.Send()

	return
}

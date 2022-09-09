// SPDX-FileCopyrightText: 2022-present Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"github.com/thakurajayL/drsm"
	"log"
	"os"
	"time"
)

type drsmInterface struct {
	initDrsm bool
	Mode     drsm.DrsmMode
	d        *drsm.Drsm
	poolName string
}

var drsmIntf drsmInterface

func scanChunk(i int32) bool {
	log.Println("Received callback from module to scan Chunk resource ", i)
	return false
}

func initDrsm(resName string) {

	if drsmIntf.initDrsm == true {
		return
	}
	drsmIntf.initDrsm = true
	drsmIntf.poolName = resName

	podn := os.Getenv("HOSTNAME") // pod-name
	podi := os.Getenv("POD_IP")
	podId := drsm.PodId{PodName: podn, PodIp: podi}
	db := drsm.DbInfo{Url: "mongodb://mongodb-arbiter-headless", Name: "sdcore"}

	t := time.Now().UnixNano()
	opt := &drsm.Options{}
	if t%2 == 0 {
		log.Println("Running in Demux Mode")
		drsmIntf.Mode = drsm.ResourceDemux
	} else {
		opt.ResourceValidCb = scanChunk
		opt.IpPool = make(map[string]string)
		opt.IpPool["pool1"] = "192.168.1.0/24"
		opt.IpPool["pool2"] = "192.168.2.0/24"
	}
	drsmIntf.d, _ = drsm.InitDRSM(resName, podId, db, opt)
}

func AllocateInt32One(resName string) int32 {
	id, err := drsmIntf.d.AllocateInt32ID()
	if err != nil {
		log.Println("Id allocation error ", err)
		return 0
	}
	log.Printf("Received id %v ", id)
	return id
}

func AllocateInt32Many(resName string, number int32) []int32 {
	// code to acquire more than 1000 Ids
	var resIds []int32
	var count int32 = 0

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			id, _ := drsmIntf.d.AllocateInt32ID()
			if id != 0 {
				resIds = append(resIds, id)
			}
			log.Printf("Received id %v ", id)
			count++
			if count >= number {
				return resIds
			}
		}
	}
}

func ReleaseInt32One(resName string, resId int32) error {
	err := drsmIntf.d.ReleaseInt32ID(resId)
	if err != nil {
		log.Println("Id release error ", err)
		return err
	}
	return nil
}

func IpAddressAllocOne(pool string) (string, error) {
	ip, err := drsmIntf.d.AcquireIp(pool)
	if err != nil {
		log.Printf("%v : Ip allocation error ", pool, err)
		return "", err
	}
	log.Printf("%v : Received ip %v ", pool, ip)
	return ip, nil
}

func IpAddressAllocMany(pool string, number int32) []string {
	var resIds []string
	var count int32 = 0

	ticker := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			ip, err := drsmIntf.d.AcquireIp(pool)
			if err != nil {
				log.Printf("%v : Ip allocation error %v", pool, err)
			} else {
				log.Printf("%v : Received ip %v ", pool, ip)
				resIds = append(resIds, ip)
			}
			count++
			if count >= number {
				return resIds
			}
		}
	}
}

func IpAddressRelease(pool, ip string) error {
	err := drsmIntf.d.ReleaseIp(pool, ip)
	return err
}

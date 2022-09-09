// SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0


package drsm

import (
	"fmt"
	"log"
)

type DbInfo struct {
	Url  string
	Name string
}

type PodId struct {
	PodName string `bson:"podName,omitempty" json:"podName,omitempty"`
	PodIp   string `bson:"podIp,omitempty" json:"podIp,omitempty"`
}

type DrsmMode int

const (
	ResourceClient DrsmMode = iota + 0
	ResourceDemux
)

type Options struct {
	ResIdSize       int32 //size in bits e.g. 32 bit, 24 bit.
	Mode            DrsmMode
	ResourceValidCb func(int32) bool // return if ID is in use or not used
	IpPool          map[string]string
}

func InitDRSM(sharedPoolName string, myid PodId, db DbInfo, opt *Options) (*Drsm, error) {
	log.Println("CLIENT ID: ", myid)

	d := &Drsm{sharedPoolName: sharedPoolName,
		clientId: myid,
		db:       db,
		mode:     ResourceClient}

	d.ConstuctDrsm(opt)

	return d, nil
}

func (d *Drsm) AllocateInt32ID() (int32, error) {
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not allocate Resource index ")
		err := fmt.Errorf("Demux mode does not allow Resource Id allocation")
		return 0, err
	}
	for _, c := range d.localChunkTbl {
		if len(c.FreeIds) > 0 {
			return c.AllocateIntID(), nil
		}
	}
	c, err := d.GetNewChunk()
	if err != nil {
		err := fmt.Errorf("Ids not available")
		return 0, err
	}
	return c.AllocateIntID(), nil
}

func (d *Drsm) ReleaseInt32ID(id int32) error {
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not release Resource index ")
		err := fmt.Errorf("Demux mode does not allow Resource Id allocation")
		return err
	}

	chunkId := id >> 10
	chunk, found := d.localChunkTbl[chunkId]
	if found == true {
		chunk.ReleaseIntID(id)
		return nil
	} else {
		chunk, found := d.scanChunks[chunkId]
		if found == true {
			chunk.ReleaseIntID(id)
			return nil
		}
	}
	err := fmt.Errorf("Unknown Id")
	return err
}

func (d *Drsm) FindOwnerInt32ID(id int32) (*PodId, error) {
	chunkId := id >> 10
	chunk, found := d.globalChunkTbl[chunkId]
	if found == true {
		podId := chunk.GetOwner()
		return podId, nil
	}
	err := fmt.Errorf("Unknown Id")
	return nil, err
}

func (d *Drsm) AcquireIp(pool string) (string, error) {
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not allocate Ip ")
		err := fmt.Errorf("Demux mode does not allow Resource allocation")
		return "", err
	}
	return d.acquireIp(pool)
}

func (d *Drsm) ReleaseIp(pool, ip string) error {
	if d.mode == ResourceDemux {
		log.Println("Demux mode can not Release Resource ")
		err := fmt.Errorf("Demux mode does not allow Resource Release")
		return err
	}
	return d.releaseIp(pool, ip)
}

//add new api for add ip pool, remove ip pool

// SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
package drsm

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

func (d *Drsm) podDownDetected() {
	fmt.Println("Started Pod Down goroutine")
	for {
		select {
		case p := <-d.podDown:
			fmt.Println("Pod Down detected ", p)
			// Given Pod find out current Chunks owned by this POD
			pd := d.podMap[p]
			for k, _ := range pd.podChunks {
				c, found := d.globalChunkTbl[k]
				fmt.Printf("Found : %v chunk : %v ", found, c)
				go c.claimChunk(d)
			}
		}
	}
}

func (c *chunk) claimChunk(d *Drsm) {
	if d.mode != ResourceClient {
		fmt.Println("claimChunk ignored demux mode ")
		return
	}
	// try to claim. If success then notification will update owner.
	fmt.Println("claimChunk started")
	docId := fmt.Sprintf("chunkid-%d", c.Id)
	update := bson.M{"_id": docId, "type": "chunk", "podId": d.clientId.PodName}
	filter := bson.M{"_id": docId, "podId": c.Owner.PodName}
	updated := d.mongo.RestfulAPIPutOnly(d.sharedPoolName, filter, update)
	if updated == nil {
		// TODO : don't add to local pool yet. We can add it only if scan is done.
		fmt.Println("claimChunk success ")
		c.Owner.PodName = d.clientId.PodName
		c.Owner.PodIp = d.clientId.PodIp
		go c.scanChunk(d)
	} else {
		fmt.Println("claimChunk failure ")
	}
}

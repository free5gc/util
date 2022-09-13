// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"github.com/omec-project/util/mongoapi"
	"log"
)

var mongoHndl *mongoapi.MongoClient
// TODO : take DB name from helm chart
// TODO : inbuild shell commands to

func main() {
	log.Println("dbtestapp started")

	// connect to mongoDB
	mongoHndl, _ = mongoapi.SetMongoDB("sdcore", "mongodb://mongodb-arbiter-headless")

	initDrsm("resourceids")

	//blocking
	http_server()
}

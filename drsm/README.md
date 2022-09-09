<!--
SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>

SPDX-License-Identifier: Apache-2.0

-->

# drsm
Distributed Resource Sharing Module (DRSM)

Resources can be
    - integer numbers (TEID, SEID, NGAPIDs,TMSI,...)
    - IP address pool

Modes
    - demux mode : just listen and get mapping about PODS and their resource assignments
        * can be used by sctplb, upf-adapter
    - client mode : Learn about other clients and their resource mappings
        * can be used by AMF pods, SMF pod

Dependency
    - MongoDB should run in cluster(replicaset) Mode or sharded Mode

Limitation:
    -  MongoDBlib keeps global client variable. So only 1 DB connection and 1 database allowed.
    -  We can not have multiple Database connections with above limitation.

Testing
    -  All the DRSM clients discover other clients through pub/sub
    -  Allocate resource id ( indirectly chunk). Other Pods should get notification of newly allocated chunk
    -  POD down event should be detected
    -  Get candidate ORPHAN chunk list once POD down detected
    -  CLAIM chunk to change owner
    -  Through notification other PODS should detect if CHUNK is claimed
    -  Run large number of clients and bring down replicaset by 1..All other pod would try to claim chunks of crashed pod.
       we should see only 1 client claiming it successfully
    -  If some pod is started late and already there are number of documents in collections. Then does stream provide
       old docs as well ? No. Added code to read existing docs.
    -  Multiple Pods trying to allocate same Chunkid. dbInsert only succeeds for one client. Does DRSM handle error and retry other Chunk
    -  Clear Separation of demux API vs regular CLIENT API
    -  Callback should be available where chunk scanning (resource id usage) can be done with help of application
    -  Pod identity is IP address + Pod Name
    -  Allocate more than 1000 ids.. See if New chunk is allocated

TODO:
    - Release Id test
    -  IP address allocation
    -  MongoDB instance restart
    -  Rst counter to be appended to identify pod.PodId should be = K8s Pod Id + Rst Count. 
       This makes sure that restarted pod even if it comes with same name then we treat it differently
    -  min REST APIs to trigger { allocate, claim }

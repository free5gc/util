#!/bin/bash
# SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0


set -xe

while getopts t:i: flag
do
    case "${flag}" in
        t) INTTEST=${OPTARG};;
        i) IPPOOL=${OPTARG};;
    esac
done

echo "Test name : $INTTEST";
echo "IP POOL : $IPPOOL";

POD_N=`kubectl get pods -n omec | grep dbtest | cut -d' ' -f1`
POD_IP=`kubectl get pod $POD_N -n omec --template '{{.status.podIP}}'`
echo $POD_IP

if [ -z "$POD_IP" ]
then
echo "POD IP empty"
return
fi

if [ ! -z $INTTEST ]
then
curl -X POST $POD_IP:8000/app/v1/integer-resource/resourceids?num=1
else
echo "Integer resource request not set "
fi

if [ ! -z $IPPOOL ]
then
curl -X POST $POD_IP:8000/app/v1/ipv4-resource/$IPPOOL
else
echo "IP resource request not set "
fi

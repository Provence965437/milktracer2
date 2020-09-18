#!/bin/bash
GETOPT_ARGS=`getopt -o n:p:u:w:a:s:d:v: -al ip:,port:,user:,pwd:,path:,script:,debug:,version: -- "$@"`
show_usage="args:[-n]链码名字"
chaincodename=""
version=""
sequence=""
while [ -n "$1" ]
do
    case "$1" in
        -n) chaincodename=$2; shift 2;;
	-v) version=$2; shift 2;;
        -s) sequence=$2; shift 2;;
        --) break ;;
        *) echo $show_usage; break ;;
    esac
done

eval set -- "$GETOPT_ARGS"
#peer1的环境变量
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/

#peer安装链码
peer lifecycle chaincode install $chaincodename.tar.gz
#peer2环境变量
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

#peer安装链码
peer lifecycle chaincode install $chaincodename.tar.gz



#!/bin/bash
set -x
GETOPT_ARGS=`getopt -o n:p:u:w:a:s:d:v:c: -al ip:,port:,user:,pwd:,path:,script:,debug:,version: -- "$@"`
show_usage="args:[-n]链码名字\
		[-a]参数列表格式“”，“”，“”\
                [-c]调用函数"
chaincodename=""
args=""
call=""
while [ -n "$1" ]
do
    case "$1" in
        -n) chaincodename=$2; shift 2;;
        -a) args=$2; shift 2;;
        -c) call=$2; shift 2;;
        --) break ;;
        *) echo $show_usage; break ;;
    esac
done

eval set -- "$GETOPT_ARGS"

export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
#peer2环境变量
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# calljson="'{\"function\":"\"$call\"",\"Args\":["$args"]}'"
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n $chaincodename --peerAddresses localhost:7051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"'$call'","Args":['$args']}'
echo $args

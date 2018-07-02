echo "installing bankcc"
peer chaincode install -n bankcc -v 1.0 -p github.com/chaincode/Bank/
echo "instantiating bankcc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n bankcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing businesscc"
peer chaincode install -n businesscc -v 1.0 -p github.com/chaincode/Business/
echo "instantiating businesscc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n businesscc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing walletcc"
peer chaincode install -n walletcc -v 1.0 -p github.com/chaincode/Wallet/
echo "instantiating walletcc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n walletcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing disbursementcc"
peer chaincode install -n disbursementcc -v 1.0 -p github.com/chaincode/Transactions/Disbursement/
echo "instantiating disbursementcc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n disbursementcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing txnbalcc"
peer chaincode install -n txnbalcc -v 1.0 -p github.com/chaincode/TxnBalance/
echo "instantiating txnbalcc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n txnbalcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing txncc"
peer chaincode install -n txncc -v 1.0 -p github.com/chaincode/Transactions/
echo "instantiating txncc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n txncc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing loancc"
peer chaincode install -n loancc -v 1.0 -p github.com/chaincode/Loan/
echo "instantiating loancc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n loancc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"

echo "installing loanbalcc"
peer chaincode install -n loanbalcc -v 1.0 -p github.com/chaincode/LoanBalance/
echo "instantiating loanbalcc"
peer chaincode instantiate -o orderer.example.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C myc -n loanbalcc -v 1.0 -c '{"Args":[]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
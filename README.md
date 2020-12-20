# WOOD-SUPPLYCHAIN-POC

# Steps to setup the hyperledger fabric network

Step-1: Clone the Repository:

    git clone https://github.com/Phaniramvarma/wood-supply-chain-blockchain.git 

Step-2: 

    cd blockchain-network/wood-supplychain/
  
  Remove old cryptomaterial and genesis.block, channel transaction files
  
    rm -rf crypto-config channel-artifacts/genesis.block channel-artifacts/genesis.json channel-artifacts/channel-transactions
  

Step-3: Generate the crypto-material

    cd scripts/Channel-Artifacts
  
  Generate certificates for organizations peer nodes and orderers
  
    ./generateCerts.sh

  Generate Genesis block
  
    ./generateGensisBlock.sh

Step-4: Using defined docker-compose yaml files for orderers, start the orderer for our blockchain network

    cd ../development
    docker-compose -f docker-compose-orderer0.yaml up -d
  

Step-5: using defined docker-compose yaml files for peer organizations, start the peers,couchdb,ca,cli for our four organizations(forest,cutter,manufacturer,logistics)

    docker-compose -f docker-compose-forest.yaml up -d
    docker-compose -f docker-compose-cutter.yaml up -d
    docker-compose -f docker-compose-logistics.yaml up -d
    docker-compose -f docker-compose-manufacture.yaml up -d
  

Step-6: Create channel transaction files using developed scripts

    ./channel-tx-gen.sh supplychain
  
Step-7: Login to forest cli container for creating supplychain channel

    docker exec -it cli.peer0-forest.com /bin/bash
    cd channel-artifacts/channel-transactions/supplychain/
    peer channel create -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem -f ./supplychain_channel.tx -c supplychain
  
  
  If the channel creation is success then you will get channels genesis block/channel 0 block
  
  Using the 0 block(supplychain.block) join the peers into the channel.
  
 
Step-8: Join forest peer into the channel

    peer channel join -b supplychain.block
  
  After joining the peer to channel check the peer joined the channel or not using below command
  
    peer channel list
  

Step-9: Now add other organizations peers to channel

  Adding Cutter organization peer to supplychain channel
  
    docker exec -it cli.peer0-cutter.com /bin/bash
    cd channel-artifacts/channel-transactions/supplychain/
    peer channel join -b supplychain.block
  
  After joining the peer to channel check the peer joined the channel or not using below command
  
    peer channel list
  
  
  Adding Logistics organization peer to supplychain channel
  
    docker exec -it cli.peer0-logistics.com /bin/bash
    cd channel-artifacts/channel-transactions/supplychain/
    peer channel join -b supplychain.block
  
  After joining the peer to channel check the peer joined the channel or not using below command
  
    peer channel list
  
  
  Adding Manufacture organization peer to supplychain channel
  
    docker exec -it cli.peer0-manufacture.com /bin/bash
    cd channel-artifacts/channel-transactions/supplychain/
    peer channel join -b supplychain.block
  
  After joining the peer to channel check the peer joined the channel or not using below command
  
    peer channel list
  
  

Step-10: Deploying Chaincode onto peers

    step-i: creating chaincode package
  
    docker exec -it cli.peer0-forest.com /bin/bash
    peer lifecycle chaincode package supplychain_1.0.tar.gz --path github.com/chaincodes/woodsupplychain/ --label supplychain_1.0
    Above command outputs *supplychain_1.0.tar.gz*
    
    step-ii: Install chaincode on all organizations:
  
    docker exec -it cli.peer0-forest.com /bin/bash
    peer lifecycle chaincode install supplychain_1.0.tar.gz
    exit
    
    docker exec -it cli.peer0-cutter.com /bin/bash
    peer lifecycle chaincode install supplychain_1.0.tar.gz
    exit
    
    docker exec -it cli.peer0-logistics.com /bin/bash
    peer lifecycle chaincode install supplychain_1.0.tar.gz
    exit
    
    docker exec -it cli.peer0-manufacture.com /bin/bash
    peer lifecycle chaincode install supplychain_1.0.tar.gz
    exit
  
    step-iii: Approve and check commitreadyness on all organizations:
  
    Execute below command on logging to all organization peers
    
    docker exec -it cli.peer0-forest.com /bin/bash
    peer lifecycle chaincode approveformyorg -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem --channelID supplychain --name supplychain --version 1.0 --package-id supplychain_1.0:2250a1459517ae8e1e8e5cb3ec5055941e07011f78e4e2fe3b127d2c8d01c367 --sequence 1
  
    docker exec -it cli.peer0-cutter.com /bin/bash
    peer lifecycle chaincode approveformyorg -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem --channelID supplychain --name supplychain --version 1.0 --package-id supplychain_1.0:2250a1459517ae8e1e8e5cb3ec5055941e07011f78e4e2fe3b127d2c8d01c367 --sequence 1
    
    docker exec -it cli.peer0-manufacture.com /bin/bash
    peer lifecycle chaincode approveformyorg -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem --channelID supplychain --name supplychain --version 1.0 --package-id supplychain_1.0:2250a1459517ae8e1e8e5cb3ec5055941e07011f78e4e2fe3b127d2c8d01c367 --sequence 1
    
    docker exec -it cli.peer0-logistics.com /bin/bash
    peer lifecycle chaincode approveformyorg -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem --channelID supplychain --name supplychain --version 1.0 --package-id supplychain_1.0:2250a1459517ae8e1e8e5cb3ec5055941e07011f78e4e2fe3b127d2c8d01c367 --sequence 1
  
    step-iv: Approve the chancode on channel
  
    docker exec -it cli.peer0-forest.com /bin/bash
    peer lifecycle chaincode commit -o orderer0.supplychain.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/supplychain.com/orderers/orderer0.supplychain.com/msp/tlscacerts/tlsca.supplychain.com-cert.pem --channelID supplychain --name supplychain --version 1.0 --sequence 1
    

Step-10: After Installing,Approving,Committing the chaincode on channel. start the API for our network for all organizations

  Forest Api:
  
    cd ../chaincodes/forest-api
    pm2 start app.js
  
  
  Manufacture Api:
  
    cd ../chaincodes/manufacture-api
    pm2 start app.js
    
  Logistics Api:
  
    cd ../chaincodes/logistics-api
    pm2 start app.js
  
  Cutter Api:
  
    cd ../chaincodes/cutter-api
    pm2 start app.js

#Next we need to generate genesis block for orderer. It can be done by using configtxgen script
# generate orderer.block on config directory

export FABRIC_CFG_PATH=/home/ubuntu/blockchain-network/wood-supplychain/fabric-config/
../../bin/configtxgen -profile MultiOrgsOrdererGenesis -outputBlock ../../channel-artifacts/genesis.block -channelID system-channel

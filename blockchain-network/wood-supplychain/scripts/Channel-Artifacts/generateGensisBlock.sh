#Next we need to generate genesis block for orderer. It can be done by using configtxgen script
# generate orderer.block on config directory

export FABRIC_CFG_PATH=../../fabric-config/
../../bin/configtxgen -profile MultiOrgsOrdererGenesis -outputBlock ../../channel-artifacts/genesis.block -channelID system-channel

if ! [ "$#" -ge 1 ] ; then
        echo "$1: exactly 1 arguments expected [channel-tx-gen.sh <channel-name> ]"
        exit 0
fi
# generate channel config transaction on config/channel.tx

export FABRIC_CFG_PATH=/home/ubuntu/blockchain-network/wood-supplychain/fabric-config/

mkdir -p ../../channel-artifacts/channel-transactions/$1

../../bin/configtxgen -profile SupplyChainOrgsChannel -outputCreateChannelTx ../../channel-artifacts/channel-transactions/$1/$1_channel.tx -channelID $1


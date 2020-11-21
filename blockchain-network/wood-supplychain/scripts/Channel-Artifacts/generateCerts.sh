##################################################################################################################
#
# Now we need to generate certificates/keys for our orderer, peers and ca.
#  We are using cryptogen script to generate them. This script will generates crypto materials in crypto directory.
#
####################################################################################################################

# generates certificates and keys on crypto-config directory
../../bin/cryptogen generate --config=../../fabric-config/crypto-config.yaml --output=../../crypto-config

####################################################################################################################33333333333333333333
#
#   In order to work with hlf we need to have crypto certificates, genesis block and channel configurations.
#  To generate them we are using scripts on bin directory and configurations(crypto-config.yaml and configtx.yaml) on config directory.
#        So we need to add bin directory $PATH and set the FABRIC_CFG_PATH with config directory (where configxtx.yaml exists).
#
##########################################################################################################################################

# add bin directory to PATH
export PATH=${PWD}/bin:${PWD}:$PATH

# define fabric config directory
export FABRIC_CFG_PATH=${PWD}/fabric-config/

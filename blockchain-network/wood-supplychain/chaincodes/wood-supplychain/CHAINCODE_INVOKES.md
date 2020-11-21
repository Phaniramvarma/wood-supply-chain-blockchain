SET_REWARD/UPDATE_REWARD:(addReward)

	peer chaincode invoke -C reward -n reward -o orderer.uec.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/uec.com/msp/tlscacerts/tlsca.uec.com-cert.pem -c '{"Args":["addReward","{\"id\":\"9000580405\",\"createdAt\":1604301987722,\"updatedAt\":1604301989684,\"payoutType\":\"DOWNLINE_MATCHING\",\"packageType\":\"AFFILIATE\",\"packageTypeId\":2,\"saleVolume\":500,\"saleVolumeInXrp\":2070.810033,\"distributeVolume\":5,\"distributeVolumeInXrp\":20.7081,\"paymentStatus\":\"pending\",\"active\": 1,\"type\":\"SuperMatch\",\"memberPackageId\":49467,\"isMigrated\":0,\"profiteerUserId\":517,\"profiteerMemberId\":8536,\"userId\":619,\"member\":12358,\"comment\":\"null\"}"]}'


	Response:
		Chaincode invoke successful. result: status:200 payload:"{\"RewardId\":\"9000580405\",\"data\":{\"asset_type\":\"reward\",\"id\":\"9000580405\",\"createdAt\":1604301987722,\"updatedAt\":1604301989684,\"payoutType\":\"DOWNLINE_MATCHING\",\"packageType\":\"AFFILIATE\",\"packageTypeId\":2,\"saleVolume\":500,\"saleVolumeInXrp\":2070.810033,\"distributeVolume\":5,\"distributeVolumeInXrp\":20.7081,\"paymentStatus\":\"pending\",\"active\":1,\"type\":\"SuperMatch\",\"memberPackageId\":49467,\"isMigrated\":0,\"profiteerUserId\":517,\"profiteerMemberId\":8536,\"userId\":619,\"member\":12358,\"comment\":\"null\"},\"message\":\"Reward added succesfull\",\"trxnID\":\"3327209699d1798ae7828a94068d3e16843b83c672741d24a1e0907e234dfbce\"}"

	

GET_REWARD_DETAILS:(getRewardById):
	peer chaincode query -C reward -n reward -o orderer.uec.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/uec.com/msp/tlscacerts/tlsca.uec.com-cert.pem -c  '{"Args":["getRewardById","9000580405"]}'

	Response:
		{"data":{"asset_type":"reward","id":"9000580405","createdAt":1604301987722,"updatedAt":1604301989684,"payoutType":"DOWNLINE_MATCHING","packageType":"AFFILIATE","packageTypeId":2,"saleVolume":500,"saleVolumeInXrp":2070.810033,"distributeVolume":5,"distributeVolumeInXrp":20.7081,"paymentStatus":"pending","active":1,"type":"SuperMatch","memberPackageId":49467,"isMigrated":0,"profiteerUserId":517,"profiteerMemberId":8536,"userId":619,"member":12358,"comment":"null"},"status":"true"}










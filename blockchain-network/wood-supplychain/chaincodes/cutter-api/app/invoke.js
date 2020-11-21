

'use strict';                                                                                                                                                                                                                                                                 


const { Gateway, Wallets } = require('fabric-network');                                
const path = require('path');
//const fs = require('fs');
var helper = require('./helper.js');                                                                                                                   
const { stringify } = require('querystring');                                                                          
const { getHashes } = require('crypto');                                        

async function invokeTransaction(fcn, args) {                                                                                                                                                                                                                                                                                                                  
    try {                     
        var result;                                                            
        var response;            
        var errResponse;              
        let orgName=args[0];                                             
        let userName=args[1];                                                                                            
        const ccp = await helper.getCCP(orgName); 

        // Create a new file system based wallet for managing identities.                                             
        const wallet = await helper.getWallet();
        //console.log(`Wallet path: ${walletPath}`);
                                                                                         
        const identity = await wallet.get(userName);                                     
        if (!identity) {
            console.log('An identity for the user "appUser" does not exist in the wallet');
            console.log('Run the registerUser.js application before retrying');
            return;
        }
        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();

	console.info(ccp)

        await gateway.connect(ccp, { wallet, identity: userName, discovery: { enabled: true, asLocalhost: true } });

	
        //console.info("Get the network (channel) our contract is deployed to.")
        // Get the network (channel) our contract is deployed to.
    
        const network = await gateway.getNetwork('supplychain');

        // Get the contract from the network.
        const contract = network.getContract('supplychain');

        switch(fcn){

                 case "addTreeDetailsByForestOfficer":
			result= await contract.submitTransaction(fcn,`${args[2]}`);
			response="Tree Details Added Successfully";
			break;
		case "updateTreeDetailsByForestOfficer":
			result= await contract.submitTransaction(fcn,`${args[2]}`,`${args[3]}`);
			response="Tree Details Updated Successfully";
			break;
		case "logUpdateByCutter":
			result= await contract.submitTransaction(fcn,`${args[2]}`);
			response="logUpdateByCutter Updated Successfully";
			break;
		case "logStatusUpdateByCutter":
			result= await contract.submitTransaction(fcn,`${args[2]}`,`${args[3]}`);
			response="logStatusUpdateByCutter Updated Successfully";
			break;
		case "loadingUpdateByLogistics":
			result= await contract.submitTransaction(fcn,`${args[2]}`);
			response="loadingUpdateByLogistics Updated Successfully";
			break;
		case "loadingStatusUpdateByLogistics":
			result= await contract.submitTransaction(fcn,`${args[2]}`,`${args[3]}`);
			response="loadingStatusUpdateByLogistics Updated Successfully";
			break;
		case "treeDetailsUpdateByManufacturer":
			result= await contract.submitTransaction(fcn,`${args[2]}`);
			response="treeDetailsUpdateByManufacturer Updated Successfully";
			break;
                case "getTreeDetailsById":
			result = await contract.evaluateTransaction(fcn,`${args[2]}`);
			result= result.toString();
			result = result.replace(/\\/g, '');
			return {
			       "statusCode":200,
				"name":"Success",
				"code":"Tree Details",
				"stack":"",
				"message":{"TreeDetails":JSON.parse(result)}    
                        }  ;                               
                case "queryWoodTrack":
			result = await contract.evaluateTransaction(fcn,`${args[2]}`);
			result= result.toString();
			result = result.replace(/\\/g, '');
			return {
			       "statusCode":200,
				"name":"Success",
				"code":"Tree Details",
				"stack":"",
				"message":{"TreeDetails":JSON.parse(result)}    
                        }  ;                               
        }
	    await gateway.disconnect();
	    return {success: true, message: response};
	}catch (error) {
		return {success: false, message: `Operation failed  : ${error}`};    
	}
}

exports.invokeTransaction = invokeTransaction;

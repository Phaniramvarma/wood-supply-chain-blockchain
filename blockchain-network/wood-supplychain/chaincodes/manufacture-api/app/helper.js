/**
 * Copyright 2017 IBM All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the 'License');
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an 'AS IS' BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */
'use strict'
const { Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const path = require('path');
const { get } = require('http');
const { stringify } = require('querystring');  
var _ = require('lodash');
var invokeTxn = require('./invoke.js');



async function getcaInfo(orgName){
        const pcn=orgName+'.com'
        const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', pcn, 'connection-'+orgName+'.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities['ca.'+pcn];
	console.info(caInfo)
        return caInfo;
}

async function getCCP(orgName){
    const pcn=orgName+'.com'
    const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', pcn, 'connection-'+orgName+'.json');
    const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8')); 
    return ccp;
};

async function getWallet(){
	const walletPath = path.join(process.cwd(), 'wallet');
	const wallet = await Wallets.newFileSystemWallet(walletPath);
	console.log(`Wallet path: ${walletPath}`);
	return wallet;;
};

async function getMSP(orgName){
	const MSP=orgName+'MSP';
	const OrgMSP=MSP;
	//const OrgMSP=MSP.charAt(0).toUpperCase() + MSP.slice(1);
        //console.info(MSP)
        //console.info(MSP.toUpperCase())
	//return MSP.toUpperCase();
	return OrgMSP;
};

async function registerUser(orgName, userName) {
    try {
        const caInfo=await getcaInfo(orgName);
        console.log(caInfo);
        const caURL = caInfo.url;
  
        const ca = new FabricCAServices(caURL);
          const OrgMSP= await getMSP(orgName);
        // Create a new file system based wallet for managing identities.
        const wallet=await getWallet();
        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(userName);
        if (userIdentity) {
            return response = 'An identity for the user '+userName+' already exists in the wallet';

        }

        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, 'admin');

        // Register the user, enroll the user, and import the new identity into the wallet.
        const secret = await ca.register({
            affiliation: orgName+'.department1',
            enrollmentID: userName,
            role: 'client'
        }, adminUser);
        const enrollment = await ca.enroll({
            enrollmentID:  userName,
            enrollmentSecret: secret
        });
        const x509Identity = {
            credentials: {
                certificate: enrollment.certificate,
                privateKey: enrollment.key.toBytes(),

            },
            mspId: OrgMSP,
            type: 'X.509',
        };
     await wallet.put(userName, x509Identity);
       
        var response = {
            success: true,
            secret: secret,
            x509Identity,  
            message: 'Successfully registered and enrolled admin user '+userName+' and imported it into the wallet successfully enrolled  user "admin" and imported it into the wallet',
        };
        return response;

    } catch (error) {
        return response = `Failed to enroll admin user "${userName}": ${error}`;
    }
};

async function revokeUser(orgName, userName) {

    try {
        var response;

        const caInfo=await getcaInfo(orgName);
        const caURL = caInfo.url;

        const ca = new FabricCAServices(caURL);
          const OrgMSP= await getMSP(orgName);
        // Create a new file system based wallet for managing identities.
        const wallet=await getWallet();
        // Check to see if we've already enrolled the user.
        const userIdentity = await wallet.get(userName);
       
       
         
        // Check to see if we've already enrolled the admin user.
        const adminIdentity = await wallet.get('admin');
        if (!adminIdentity) {
            console.log('An identity for the admin user "admin" does not exist in the wallet');
            console.log('Run the enrollAdmin.js application before retrying');
            return;
        }

        // build a user object for authenticating with the CA
        const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
        const adminUser = await provider.getUserContext(adminIdentity, 'admin');

        // Revoke
     //   response = await ca.revoke({
       //     affiliation: orgName+'.department1',
         //   enrollmentID: userName,
           // role: 'client'
        //}, adminUser);
        await wallet.remove(userName);
        var fcn = "RevokeUser";
             let userN='appUser';   
        var args = [
         orgName,
         userN,
         userName
                  ];

       let  result =  await invokeTxn.invokeTransaction(fcn, args);  

      // response = await ca.generateCRL({
        //revokedAfter: '2020-01-13T16:39:57-08:00',
        //revokedbefore: '2020-09-20T16:39:57-08:00'      
    //}, adminUser);


    //let crlResponse=await ca.generateCRL(ts,adminUser);
       
       /* var response = {
            success: true,
            secret: secret,
            message: 'Successfully registered and enrolled admin user '+userName+' and imported it into the wallet successfully enrolled  user "admin" and imported it into the wallet',
        };*/
        return true;
    } catch (error) {
        return false;
        }
};

async function enrolladmin(orgName, userName, passWord) {
        try {                         
                    const caInfo=await getcaInfo(orgName);
                   const caURL=caInfo.url
                   const caTLSCACerts = caInfo.tlsCACerts.pem;

            const ca = new FabricCAServices(caURL, { trustedRoots: caTLSCACerts, verify: false },caInfo.caName);
            console.info(ca)
            const OrgMSP= await getMSP(orgName);

            // Create a new file system based wallet for managing identities.
            //const walletPath = path.join(process.cwd(), orgName+'-wallet');
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet=await getWallet();

    
            // Check to see if we've already enrolled the admin user.
            const identity = await wallet.get(userName);
            if (identity) {
                
                  return response = 'An identity for the admin user "admin" already exists in the wallet';
            }
    
            // Enroll the admin user, and import the new identity into the wallet.
            const enrollment = await ca.enroll({ enrollmentID: userName, enrollmentSecret: passWord });
            const x509Identity = {
                credentials: {
                    certificate: enrollment.certificate,
                    privateKey: enrollment.key.toBytes(),
                },
                mspId:OrgMSP ,
                type: 'X.509',
            };
            await wallet.put(userName, x509Identity);
    
            var response = {
                success: true,
                message: 'Successfully enrolled admin user "admin" and imported it into the wallet',
            };
            return response;
    
        } catch (error) {
            return response = `Failed to enroll admin user "${userName}": ${error}`;
        }
    }
    ;

    async function checkUser(orgName, userName) {
        try {
            const caInfo=await getcaInfo(orgName);
            console.log(caInfo);
            const caURL = caInfo.url;
              var response;
              let userN='appUser';
            const ca = new FabricCAServices(caURL);
              const OrgMSP= await getMSP(orgName);
            // Create a new file system based wallet for managing identities.
            const wallet=await getWallet();
            // Check to see if we've already enrolled the user.
            const userIdentity = await wallet.get(userName);
            if (userIdentity) {
                return response = 'Identity Exists';
    
            }
            else
            {
                var fcn = "queryUserById";
                
            	var args = [
	         	orgName,
                 userN,
                 userName
                      	];

                let user = await invokeTxn.invokeTransaction(fcn, args);  
              
                response=  user.message.userDetails.toString();
                response=JSON.parse(response);
               // let isActive=response.userDetails.isActive;
                if(response.isActive=='1')
                return 'Identity not issued';
                if(response.isActive=='0')
                return 'Identity revoked';
                

               //  if(user.isActive=='0')
                 //response = user.message.userDetails.isActive;


            }
            return response;
    
        } catch (error) {
            return ` ${error}`;
        }
    };
     

    async function checkIdentity(orgName, userName,x509Identity) {
        try {
            const caInfo=await getcaInfo(orgName);
            console.log(caInfo);
            const caURL = caInfo.url;
    var response;
            const ca = new FabricCAServices(caURL);
              const OrgMSP= await getMSP(orgName);
            // Create a new file system based wallet for managing identities.
            const wallet=await getWallet();
            // Check to see if we've already enrolled the user.
            var userIdentity = await wallet.get(userName);
          response= _.isEqual(x509Identity,userIdentity );
             
            
         
            return response;
    
        } catch (error) {
            return response = ` ${error}`;
        }
    };
     

    
    
    exports.enrolladmin = enrolladmin;
exports.registerUser = registerUser;
exports.revokeUser = revokeUser;

exports.getCCP = getCCP;
exports.getMSP = getMSP;
exports.caInfo = getcaInfo;
exports.getWallet=getWallet;
exports.checkIdentity=checkIdentity;

exports.checkUser=checkUser

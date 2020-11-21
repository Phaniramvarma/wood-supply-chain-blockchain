/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');
const fs = require('fs');
const path = require('path');
const { get } = require('http');

async function getCA(orgName){
    const pcn=orgName+'.com'
        const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', pcn, 'connection-'+orgName+'.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new CA client for interacting with the CA.
        const caURL = ccp.certificateAuthorities['ca.'+pcn].url;
        const ca = new FabricCAServices(caURL);
        return ca;
}

async function getMSP(orgName)
{

    const MSP=orgName+'MSP';
   const OrgMSP=MSP;
//const OrgMSP=MSP.charAt(0).toUpperCase() + MSP.slice(1);
return OrgMSP;
}

async function registerUser(orgName, userName, passWord) {
    try {
        
        const caCLient=  await getCA(orgName);
        console.log(caCLient);
          const OrgMSP= await getMSP(orgName);
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

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
        const secret = await caCLient.register({
            affiliation: 'org1.department1',
            enrollmentID: userName,
            role: 'client'
        }, adminUser);
        const enrollment = await caCLient.enroll({
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
            message: 'Successfully registered and enrolled admin user '+userName+' and imported it into the walletuccessfully enrolled  user "admin" and imported it into the wallet',
        };
        return response;

    } catch (error) {
        return response = `Failed to enroll admin user "${userName}": ${error}`;
    }
}

exports.registerUser = registerUser;

/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const FabricCAServices = require('fabric-ca-client');
const { Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

async function enrolladmin(orgName, userName, passWord) {
    try {
        // load the network configuration
        const MSP=orgName+'MSP';
	const OrgMSP=MSP;

        //const OrgMSP=MSP.charAt(0).toUpperCase() + MSP.slice(1);

        const pcn=orgName+'.com'
        const ccpPath = path.resolve(__dirname, '..', 'organizations', 'peerOrganizations', pcn, 'connection-'+orgName+'.json');
        const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));
        console.log(`ccp: ${ccpPath}`);

        // Create a new CA client for interacting with the CA.
        const caInfo = ccp.certificateAuthorities['ca.'+pcn];
        const caTLSCACerts = caInfo.tlsCACerts.pem;
        const ca = new FabricCAServices(caInfo.url, { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

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
            mspId: OrgMSP,
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


exports.enrolladmin = enrolladmin;

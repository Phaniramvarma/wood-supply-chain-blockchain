

'use strict';

var log4js = require('log4js');

var logger = log4js.getLogger();
var express = require('express');
var axios = require('axios');
const cookieSession = require('cookie-session');

var bodyParser = require('body-parser');
var http = require('http');

var app = express();
var helper = require('./app/helper.js');
var invokeTxn = require('./app/invoke.js');

var host = "0.0.0.0";
var port = 3000;


var cors = require('cors');
//const { ADDRCONFIG } = require('dns');
//const { stderr } = require('process');


///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// SET CONFIGURATONS ////////////////////////////
///////////////////////////////////////////////////////////////////////////////
app.options('*', cors());
app.use(cors());
//support parsing of application/json type post data
app.use(bodyParser.json());
//support parsing of application/x-www-form-urlencoded post data
app.use(bodyParser.urlencoded({
        extended: false
}));


// cookieSession config
//app.use(cookieSession({
  //  maxAge: 24 * 60 * 60 * 1000, // One day in milliseconds
    //keys: ['randomstringhere']
//}));


//enroll the admin into dlt 
app.post('/api/enrollAdmin', async function (req, res){

	console.info("STARTED ENROLLADMIN")
    var orgName = req.body.orgName;
    var userName = req.body.userName;
    var passWord = req.body.passWord;

    if (!orgName) {
                res.json(getErrorMessage('\'orgname\''));
                return;
        }

    if (!userName) {
                res.json(getErrorMessage('\'username\''));
                return;
        }

    if (!passWord) {
                res.json(getErrorMessage('\'password\''));
                return;
        }

    let response = await helper.enrolladmin(orgName, userName, passWord);
        if (response && typeof response !== 'string') {
                res.json(response);
        } else {
                res.json({success: false, message: response});
        }

});


//register the user into dlt
app.post('/api/registerUser', async function (req, res){

        var orgName = req.body.orgName;
        var userId = req.body.userId;
        var userName=userId;


        if (!orgName) {
                res.json(getErrorMessage('\'orgname\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }
        if (!userId) {
                res.json(getErrorMessage('\'userId\''));
                return;
        }


        let response = await helper.registerUser(orgName, userName);
        if (response && typeof response !== 'string') {
                res.json(response);
        } else {
                res.json({success: false, message: response});
        }

});



//addTreeDetailsByForestOfficer into dlt
app.post('/api/addTreeDetailsByForestOfficer', async function (req, res){

	//parsing the body coming from request 
	var orgName  = req.body.orgName;
	var userName = req.body.userName;
	var Tree   = req.body.Tree
	console.log(Tree)


	//validating all input parameters are there or not
	if (!orgName) {
		res.json(getErrorMessage('\'orgName\''));
		return;
	}

	if (!userName) {
		res.json(getErrorMessage('\'userName\''));
		return;
	}	

	if (!Tree.hasOwnProperty('tree_id')) {
		res.json(getErrorMessage('\'tree_id\''));
		return;
	}

	if (Tree.tree_id.replace(/ /g, '').length===0){
		res.json(getErrorMessage('\'tree_id\''));
		return;
	}

	if (!Tree.hasOwnProperty('age')) {
		res.json(getErrorMessage('\'age\''));
		return;
	}

	if (Tree.age.replace(/ /g, '').length===0) {
		res.json(getErrorMessage('\'age\''));
		return;
	}

	if (!Tree.hasOwnProperty('quality')) {
		res.json(getErrorMessage('\'quality\''));
		return;
	}

	if (Tree.quality.replace(/ /g, '').length===0) {
		res.json(getErrorMessage('\'quality\''));
		return;
	}

	if (!Tree.hasOwnProperty('tree_dimensions')) {
		res.json(getErrorMessage('\'tree_dimensions\''));
		return;
	}

	if (Tree.tree_dimensions.replace(/ /g, '').length===0) {
		res.json(getErrorMessage('\'tree_dimensions\''));
		return;
	}

	if (!Tree.hasOwnProperty('rfid_created_time')) {
		res.json(getErrorMessage('\'rfid_created_time\''));
		return;
	}

	if (Tree.rfid_created_time.replace(/ /g, '').length===0) {
		res.json(getErrorMessage('\'rfid_created_time\''));
		return;
	}
	if (!Tree.hasOwnProperty('location')){
		res.json(getErrorMessage('\'location\''));
		return;
	}

	if (Tree.location.replace(/ /g, '').length===0) {
		res.json(getErrorMessage('\'location\''));
		return;
	}

	let json_string={"tree":Tree}
	var fcn = "addTreeDetailsByForestOfficer"
	var args = [
		orgName,
		userName,
		JSON.stringify(json_string),
		]
	logger.debug('fcn  : ' + fcn);
	logger.debug('args  : ' + args);
	let message = await invokeTxn.invokeTransaction(fcn, args);
	res.send(message);

});


//updateTreeDetailsByForestOfficer into dlt
app.post('/api/updateTreeDetailsByForestOfficer', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var TreeId  = req.body.TreeId;
	var TreeSts = req.body.TreeSts;


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!TreeId) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

	 if (TreeId.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

        if (!TreeSts) {
                res.json(getErrorMessage('\'TreeSts\''));
                return;
        }

	 if (TreeSts.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'TreeSts\''));
                return;
        }
	 var fcn = "updateTreeDetailsByForestOfficer"
        var args = [
                orgName,
                userName,
		TreeId,
		TreeSts,
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});



//logUpdateByCutter into dlt
app.post('/api/logUpdateByCutter', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var CutterDetails   = req.body.CutterDetails
        console.log(CutterDetails)


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!CutterDetails.hasOwnProperty('tree_id')) {
                res.json(getErrorMessage('\'tree_id\''));
                return;
        }

        if (CutterDetails.tree_id.replace(/ /g, '').length===0){
                res.json(getErrorMessage('\'tree-id\''));
                return;
        }

        if (!CutterDetails.hasOwnProperty('log_dimensions')) {
                res.json(getErrorMessage('\'log_dimensions\''));
                return;
        }

        if (CutterDetails.log_dimensions.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'log_dimensions\''));
                return;
        }

        if (!CutterDetails.hasOwnProperty('log_status')) {
                res.json(getErrorMessage('\'log_status\''));
                return;
        }

        if (CutterDetails.log_status.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'log_status\''));
                return;
        }

        if (!CutterDetails.hasOwnProperty('log_time')) {
                res.json(getErrorMessage('\'log_time\''));
                return;
        }

        if (CutterDetails.log_time.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'log_time\''));
                return;
        }

        if (!CutterDetails.hasOwnProperty('log_location')) {
                res.json(getErrorMessage('\'log_location\''));
                return;
        }

        if (CutterDetails.log_location.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'log_location\''));
                return;
        }

        var fcn = "logUpdateByCutter"
        var args = [
                orgName,
                userName,
                JSON.stringify(CutterDetails),
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});


//logStatusUpdateByCutter into dlt
app.post('/api/logStatusUpdateByCutter', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var TreeId  = req.body.TreeId;
        var LogSts = req.body.LogSts;


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!TreeId) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

         if (TreeId.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

        if (!LogSts) {
                res.json(getErrorMessage('\'LogSts\''));
                return;
        }

         if (LogSts.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'LogSts\''));
                return;
        }
         var fcn = "logStatusUpdateByCutter"
        var args = [
                orgName,
                userName,
                TreeId,
                LogSts,
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});


//loadingUpdateByLogistics into dlt
app.post('/api/loadingUpdateByLogistics', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var LogisticsDetails   = req.body.LogisticsDetails
        console.log(LogisticsDetails)


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!LogisticsDetails.hasOwnProperty('tree_id')) {
                res.json(getErrorMessage('\'tree_id\''));
                return;
        }

        if (LogisticsDetails.tree_id.replace(/ /g, '').length===0){
                res.json(getErrorMessage('\'tree_id\''));
                return;
        }

        if (!LogisticsDetails.hasOwnProperty('log_dimensions')) {
                res.json(getErrorMessage('\'log_dimensions\''));
                return;
        }

        if (LogisticsDetails.log_dimensions.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'log_dimensions\''));
                return;
        }

        if (!LogisticsDetails.hasOwnProperty('loading_status')) {
                res.json(getErrorMessage('\'loading_status\''));
                return;
        }

        if (LogisticsDetails.loading_status.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'loading_status\''));
                return;
        }

        if (!LogisticsDetails.hasOwnProperty('loading_time')) {
                res.json(getErrorMessage('\'loading_time\''));
                return;
        }

        if (LogisticsDetails.loading_time.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'loading_time\''));
                return;
        }

        var fcn = "loadingUpdateByLogistics"
        var args = [
                orgName,
                userName,
                JSON.stringify(LogisticsDetails),
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});

//loadingStatusUpdateByLogistics into dlt
app.post('/api/loadingStatusUpdateByLogistics', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var TreeId  = req.body.TreeId;
        var LoadingSts = req.body.LoadingSts;


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!TreeId) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

         if (TreeId.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

        if (!LoadingSts) {
                res.json(getErrorMessage('\'LoadingSts\''));
                return;
        }

         if (LoadingSts.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'LoadingSts\''));
                return;
        }
         var fcn = "loadingStatusUpdateByLogistics"
        var args = [
                orgName,
                userName,
                TreeId,
                LoadingSts,
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});


//treeDetailsUpdateByManufacturer into dlt
app.post('/api/treeDetailsUpdateByManufacturer', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var ManufactureDetails   = req.body.ManufactureDetails
        console.log(ManufactureDetails)


        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }

        if (!ManufactureDetails.hasOwnProperty('tree_id')) {
                res.json(getErrorMessage('\'tree_id\''));
                return;
        }

        if (ManufactureDetails.tree_id.replace(/ /g, '').length===0){
                res.json(getErrorMessage('\'tree_id\''));
                return;
        }

        if (!ManufactureDetails.hasOwnProperty('product_dimensions')) {
                res.json(getErrorMessage('\'product_dimensions\''));
                return;
        }

        if (ManufactureDetails.product_dimensions.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'product_dimensions\''));
                return;
        }

        if (!ManufactureDetails.hasOwnProperty('status')) {
                res.json(getErrorMessage('\'status\''));
                return;
        }

        if (ManufactureDetails.status.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'status\''));
                return;
        }

        if (!ManufactureDetails.hasOwnProperty('qr_code')) {
                res.json(getErrorMessage('\'qr_code\''));
                return;
        }

        if (ManufactureDetails.qr_code.replace(/ /g, '').length===0) {
                res.json(getErrorMessage('\'qr_code\''));
                return;
        }

        var fcn = "treeDetailsUpdateByManufacturer"
        var args = [
                orgName,
                userName,
                JSON.stringify(ManufactureDetails),
                ]
        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        res.send(message);

});



//get the tree details by tree id from dlt on supply channel
app.post('/api/getTreeDetailsById', async function (req, res){

	//parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
        var TreeId = req.body.TreeId;

	//validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }


        if (!TreeId) {
                res.json(getErrorMessage('\'TreeId\''));
                return;
        }

        var fcn = "getTreeDetailsById"
        var args = [
                        orgName,
                        userName,
                        TreeId,
                ]

        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        //res.send(JSON.parse(message));
        res.send(message);
});


//get the tree details by selector parameters from dlt on supply channel
app.post('/api/queryWoodTrack', async function (req, res){

        //parsing the body coming from request
        var orgName  = req.body.orgName;
        var userName = req.body.userName;
	var Query= req.body.Query;

        //validating all input parameters are there or not
        if (!orgName) {
                res.json(getErrorMessage('\'orgName\''));
                return;
        }

        if (!userName) {
                res.json(getErrorMessage('\'userName\''));
                return;
        }


        if (!Query) {
                res.json(getErrorMessage('\'Query\''));
                return;
        }

        var fcn = "queryWoodTrack"
        var args = [
                        orgName,
                        userName,
			JSON.stringify(Query),
                ]

        logger.debug('fcn  : ' + fcn);
        logger.debug('args  : ' + args);
        let message = await invokeTxn.invokeTransaction(fcn, args);
        //res.send(JSON.parse(message));
        res.send(message);
});



///////////////////////////////////////////////////////////////////////////////
//////////////////////////////// START SERVER /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////
var server = http.createServer(app).listen(port, function() {});
logger.info('****************** SERVER STARTED ************************');
logger.info('***************  http://%s:%s  ******************',host,port);
server.timeout = 240000;


function getErrorMessage(field) {
        var response = {
                success: false,
                message: field + ' field is missing or Invalid in the request'
        };
        return response;
}



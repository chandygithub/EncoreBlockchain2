var express = require('express')
var router = express.Router()
var url = require('url');


// define the home page route
router.get('/', function (req, res) {
	'use strict';
	'use strict';
	/*
	* Copyright IBM Corp All Rights Reserved
	*
	* SPDX-License-Identifier: Apache-2.0
	*/
	/*
	 * Chaincode query
	 */
	var url_parts = url.parse(req.url, true);
	var query = url_parts.query;
//	var data = req.body;

	//console.log("data"+JSON.stringify(data));
	//console.log("data"+JSON.parse(data));
	//var myArgs = JSON.parse(query.arguments);
	var myArgs = JSON.parse(query.arguments);
	console.log("myArgs "+myArgs);
	var fcn_args = myArgs.slice(2);
	
	console.log("Arguments " + myArgs[0] + myArgs[1]);
	console.log("fcn_args "+fcn_args);

	var Fabric_Client = require('fabric-client');
	var path = require('path');
	var util = require('util');
	var os = require('os');

	//
	var fabric_client = new Fabric_Client();

	// setup the fabric network
	var channel = fabric_client.newChannel('myc');
	var peer = fabric_client.newPeer('grpc://localhost:7051');
	channel.addPeer(peer);

	//
	var member_user = null;
	var store_path = path.join(__dirname, '/../hfc-key-store');
	console.log('Store path:' + store_path);
	var tx_id = null;

	// The below piece  till line 53 is for the purpose of creating user profile for peers as network admin and allows the fabric to check if the user has rights
	// create the key value store as defined in the fabric-client/config/default.json 'key-value-store' setting
	Fabric_Client.newDefaultKeyValueStore({
		path: store_path
	}).then((state_store) => {
		// assign the store to the fabric client
		fabric_client.setStateStore(state_store);
		var crypto_suite = Fabric_Client.newCryptoSuite();
		// use the same location for the state store (where the users' certificate are kept)
		// and the crypto store (where the users' keys are kept)
		var crypto_store = Fabric_Client.newCryptoKeyStore({ path: store_path });
		crypto_suite.setCryptoKeyStore(crypto_store);
		fabric_client.setCryptoSuite(crypto_suite);

		// get the enrolled user from persistence, this user will sign all requests
		return fabric_client.getUserContext('user1', true);
	}).then((user_from_store) => {
		if (user_from_store && user_from_store.isEnrolled()) {
			console.log('Successfully loaded user1 from persistence');
			member_user = user_from_store;
		} else {
			throw new Error('Failed to get user1.... run registerUser.js');
		}

		// queryCar chaincode function - requires 1 argument, ex: args: ['CAR4'],
		// queryAllCars chaincode function - requires no arguments , ex: args: [''],
		console.log("myArgs[0] "+myArgs[0]);
		console.log("myArgs[1] "+myArgs[1]);
		console.log("fcn_args "+fcn_args);
		const request = {
			//targets : --- letting this default to the peers assigned to the channel
			chaincodeId: myArgs[0],
			fcn: myArgs[1],
			args: fcn_args
			//args: ['']
		};

		// send the query proposal to the peer
		return channel.queryByChaincode(request);
	})
	.then((query_responses) => {
		console.log("Query has completed, checking results");
		console.log('query resp'+JSON.parse(query_responses));
		// query_responses could have more than one  results if there multiple peers were used as targets
		if (query_responses && query_responses.length == 1) {
			//console.log(""+query_responses.);
			if (query_responses[0] instanceof Error) {
				console.error("Error from query = ", query_responses[0]);
			} else {
				//console.log("Response is resolve json parse", (JSON.parse(query_responses[0].toString('utf8'))));
			for(let i = 0; i < query_responses.length; i++) {
				console.log(util.format('Query result from peer [%s]: %s', i, JSON.parse(query_responses[i].toString('utf8'))));
				console.log(util.format('Query result from peer [%s]: %s', i, query_responses[i].toString()));
				console.log('Query result from peer [%s]: %s', query_responses[i].toString('utf8'));
				//console.log('Query result from peer ', JSON.parse(query_responses[i].toString('utf8')));
			}
			//console.log(util.format('Query result from peer  %s',  JSON.parse(query_responses[0].toString('ascii'))));
			//query_responses[0].
			}
		} else {
			console.log("No payloads were returned from query");
		} 

	/**.then((response_payloads) => {
		console.log("length is "+response_payloads[0].length)
        for(let i = 0; i < response_payloads.length; i++) {
            console.log(util.format(' Query result from peer  %s ',  response_payloads[i].toString()));
        }    */
	})	
	.catch((err) => {
		console.error('Failed to query successfully :: ' + err+'Error message from go');
	});
	

});
module.exports = router;

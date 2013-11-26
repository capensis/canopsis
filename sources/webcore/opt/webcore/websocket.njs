#!/usr/bin/env node
/*
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------
*/

//####################################################
//# Default Configurations
//####################################################

var default_config = {
	amqp: {},
	nowjs: { port: 8085, debug: false, socketio_loglevel: 0, heartbeat: 60},
	mongodb: {}
};
var config = default_config;

//####################################################
//#  Logging
//####################################################

var log = {
	info:		function(message, author){
		this.print(this.date()+' INFO '+author+' '+message)
	},
	debug:		function(message, author){
		if (config.nowjs.debug)
			this.print(this.date()+' DEBUG '+author+' '+message)
	},
	warning:	function(message, author){
		this.print(this.date()+' WARN '+author+' '+message)
	},
	error:		function(message, author){
		this.print(this.date()+' ERR '+author+' '+message)
	},
	dump: 		function(message){
		if (config.nowjs.debug)
			console.log(message)
	},
	print: 		function(message){
		console.log(message)
	},
	date:		function(){
		var date = new Date()
		return date
	}
};

//####################################################
//#  Build Event and rk
//####################################################

var build_event = function(event){

	if (! event.component){
		//event.component = 
		log.error("Missing 'component' field", "build_event")
		return undefined
	}
	
	if (! event.resource)
		event.source_type = 'component'
	else
		event.source_type = 'resource'
	
	if (! event.state)
		event.state = 0

	if (! event.event_type)
		event.event_type = 'log'
		
	if (! event.output)
		event.output = ''
		
	if (! event.connector_name)	
		event.connector_name = 'canopsis'
		
	event.connector = 'websocket'
	event.timestamp = parseInt(new Date().getTime() / 1000)
	
	return event
};

var build_rk = function(event){
	//<connector>.<connector_name>.<event_type>.<source_type>.<component>[.<resource>]
	if (event.source_type == 'resource')
		return event.connector +"."+ event.connector_name +"."+ event.event_type +"."+ event.source_type +"."+ event.component +"."+ event.resource + "." + new Date().getTime()
	else
		return event.connector +"."+ event.connector_name +"."+ event.event_type +"."+ event.source_type +"."+ event.component + "." + new Date().getTime()
}

//####################################################
//#  Extend object (http://onemoredigit.com/post/1527191998/extending-objects-in-node-js)
//####################################################
Object.defineProperty(Object.prototype, "extend", {
    enumerable: false,
    value: function(from) {
        var props = Object.getOwnPropertyNames(from);
        var dest = this;
        props.forEach(function(name) {
            if (name in dest) {
                var destination = Object.getOwnPropertyDescriptor(from, name);
                Object.defineProperty(dest, name, destination);
            }else{
				// Hack
				dest[name] = Object.getOwnPropertyDescriptor(from, name).value;
			}
        });
        return this;
    }
});

//####################################################
//#  Load Module
//####################################################

log.info("Load modules ...", "main")
try {
	//var fs = require('fs');
	var mongodb = require('mongodb');
	var http = require('http');
	var nowjs = require("now");
	var amqp   = require('amqp');
	var util   = require('util');
	var iniparser = require('iniparser');
	
} catch (err) {
	log.error("Impossible to load modules", "main")
	log.dump(err)
	process.exit(1);
}
log.info(" + Ok", "main")


//####################################################
//#  Load Configurations
//####################################################

//GLOBAL
var config = {};
	
var read_config = function(callback){
	
	log.info("Read configuration's file ...", "config")
	
	var read_config_ini = function(file, field, section, callback){
		log.info(" + Read "+file+"...", "config")
		
		iniparser.parse(file, function(err, data){
			if (err) {
				log.error(err, "config");
				process.exit(1);
			} else {
				if (data[section])
					config[field] = default_config[field].extend(data[section])
				else
					config[field] = default_config[field]
					
				log.info("   + Ok", "config")
				callback()
			}
		});	
	}
	
	// MongoDB
	read_config_ini(process.env.HOME+'/etc/cstorage.conf', "mongodb", "master", function(){
		// AMQP
		read_config_ini(process.env.HOME+'/etc/amqp.conf', "amqp", "master", function(){
			// Now
			read_config_ini(process.env.HOME+'/etc/websocket.conf', "nowjs", "master", function(){
				//Main Callback
				log.info(" + Ok", "config")
				callback(config)
			});
		});
	});
}

//####################################################
//#  Connect to MongoDB
//####################################################

//GLOBAL
var mongodb_server = undefined
var mongodb_client = undefined
var mongodb_collections = {}

var init_mongo = function(callback){
	log.info("Connect to MongoDB ...", "mongodb")
	mongodb_server = new mongodb.Server(config.mongodb.host, parseInt(config.mongodb.port), {})
	mongodb_client = new mongodb.Db(config.mongodb.db, mongodb_server, {safe:false});

	mongodb_client.open(function(err, p_client) {
		if (err) {
			log.error(err, "mongodb");
		} else {
			log.info(" + Ok", "mongodb")
			callback()
		}
	});
}

var mongodb_getCollection = function(name){
	if (mongodb_collections[name])
		return mongodb_collections[name]
		
	mongodb_collections[name] = new mongodb.Collection(mongodb_client, name);
	return mongodb_collections[name]
}

// Ex: mongodb_find('object', {'crecord_type': 'account'}, { 'limit': 1 }, ['id'], console.log)
var mongodb_find = function(collection_name, filter, options, callback, callback_err){
	if (! options)
		options = {}
	
	if (mongodb_client){
		mongodb_getCollection(collection_name).find(filter, {'media_bin': 0}, options).toArray(function(err, records){
			if (err){
				log.error("Find "+collection_name+" "+filter+":", "mongodb");
				log.error(err, "mongodb");
				if (callback_err)
					callback_err(err)
			}else{
				if (callback)
					callback(records)
			}
		});
	}else{
		log.error("MongoDB Client is not ready", "mongodb");
		if (callback_err)
			callback_err()
	}		
}

var mongodb_findOne = function(collection_name, filter, options, callback, callback_err){
	if (!options)
		options = {}
		
	if (mongodb_client){
		mongodb_getCollection(collection_name).findOne(filter, options, function(err, record){
			if (err){
				log.error("FindOne "+collection_name+" "+filter+":", "mongodb");
				log.error(err, "mongodb");
				if (callback_err)
					callback_err(err)
			}else{
				if (callback)
					callback(record)
			}
		});
	}else{
		log.error("MongoDB Client is not ready", "mongodb");
		if (callback_err)
			callback_err()
	}
}

var mongodb_count = function(collection_name, filter, callback, callback_err){		
	if (mongodb_client)
		mongodb_getCollection(collection_name).count(filter, function(err, count){
			if (err){
				log.error("count "+collection_name+" "+filter+":", "mongodb");
				log.error(err, "mongodb");
				if (callback_err)
					return callback_err(err)
			}else
				return callback(count)
		});
			
	return 0
}

//####################################################
//#  Connect to AMQP Broker
//####################################################

//GLOBAL
var amqp_connection = undefined

var init_amqp = function(callback){
	log.info("Connect to AMQP Broker ...", "amqp")
	
	amqp_connection = amqp.createConnection({
		host: config.amqp.host,
		port: config.amqp.port,
		vhost: config.amqp.virtual_host
	});

	amqp_connection.addListener('ready', function(){
		log.info(" + Connected", "amqp");
		if (callback)
			callback();
	});

	amqp_connection.addListener('error', function(exception){
		log.error(" + Disconnected", "amqp");
	});	
}


//####################################################
//#  Bind AMQP Queue on Now group
//####################################################

//GLOBAL
var amqp_queues = {};
var amqp_exchanges = {};

var amqp_subscribe_queue = function(queue_name){
	var short_name = queue_name;
	var queueId = "amqp-" + queue_name;
	var queue_name = 'websocket_'+queueId
	
	log.info("Create Queue '"+queue_name+"'", "amqp")
	if (! amqp_queues[queue_name]){
		var queue = amqp_connection.queue(queue_name, {durable: false, exclusive: true}, function(){
			log.debug(" + Ok", "amqp")
				
			log.debug("Subscribe Queue '"+queue_name+"'", "amqp")
			this.subscribe( {ack:true}, function(message, headers, deliveryInfo){
				if (message['media_bin'])
					delete ['media_bin']
					
				nowjs.getGroup(queueId).now[queueId](message, deliveryInfo.routingKey)
				queue.shift()
			});
			
			log.debug("Bind '#' on '"+queue_name+"'", "amqp")
			this.bind("canopsis."+short_name, "#");
			this.on('queueBindOk', function() { log.debug(" + Ok", "amqp") });
			
			this.short_name = short_name

			amqp_queues[queue_name] = this;	
		});
	}else{
		log.info(" + Already exist", "amqp")
	}
}

var amqp_unsubscribe_queue = function(queue_name){
	queue_name = 'websocket_'+queue_name
	var queue = amqp_queues[queue_name]
	if (queue){
		log.info("Close AMQP queue '" + queue_name + "'", "amqp")
		queue.destroy()
		delete amqp_queues[queue_name];
	}
}

var amqp_publish = function(exchange, rk, message){
	if(amqp_connection){
		if (! amqp_exchanges[exchange]){
			log.info("Open exchange '"+exchange+"'", "amqp")
			amqp_exchanges[exchange] = amqp_connection.exchange(exchange, {type: "topic", durable: true, auto_delete: false});
		}
		
		log.info("Publish message to '"+rk+"@"+exchange+"'", "amqp")
		amqp_exchanges[exchange].publish(rk, message, {contentType: 'application/json', contentEncoding: 'utf-8'});
	}
}


//####################################################
//#  Start Now Server
//####################################################

var sessions = {
	sessions: {},
	clientIds: {}, 
	
	create: function(id, authId){
		if (this.check(id))
			return
		
		if (authId == undefined){
			log.error("You must specify authId !", "session")
			return
		}
			
		log.debug("Create session "+id+" ("+authId+")", "session")
		this.sessions[id] = authId
		
		if (this.clientIds[authId])
			this.clientIds[authId].push(id)
		else
			this.clientIds[authId] = [ id ]
	},
	
	drop: function(id){
		if (this.sessions[id]){
			var authId = this.sessions[id]
			log.debug("Drop session "+id+" ("+authId+")", "session")
			delete this.sessions[id]
			this.clientIds[authId].splice(this.clientIds[authId].indexOf(id), 1)
		}else{
			log.warning("Unknown session "+id, "session")
		}
	},
	
	check: function(id){
		return this.sessions[id]
	},
	
	getclientIds: function(authId){
		return this.clientIds[authId]
	}
}

var everyone = undefined
var init_now = function(callback){
	var server = http.createServer(function(req, res){});
	server.listen(parseInt(config.nowjs.port));

	everyone = nowjs.initialize(server, {socketio: {'log level': config.nowjs.socketio_loglevel}});
	
	////////////////// Utils
	var check_session = function(event){
		log.debug("Check session for "+event.now.authId+" ("+event.user.clientId+")", "nowjs");
		if (! sessions.check(event.user.clientId)){
			log.debug(" + You must auth !", "nowjs");
			return false
		}
		log.debug(" + Ok", "nowjs");
		return true
	}
	
	var check_authToken = function (clientId, authId, authToken, callback){
		if (mongodb_client) {
			mongodb_findOne('object', {'_id': authId}, {'fields': ['authkey']}, function(record){			
				if (record.authkey == authToken){
					log.info(" + Auth Ok", "nowjs")
					sessions.create(clientId, authId)
				} else {
					log.info(" + "+clientId + ": Invalid auth (authId: '"+authId+"')", "nowjs");
				}	
				callback();
			});
		}else{
			log.warning("MongoDB not ready.", "nowjs");
		}
	};
	
	////////////////// RPC	
	everyone.now.auth = function(callback){
		var clientId = this.user.clientId
		log.info("Auth " + this.now.authId + " ..." , "nowjs");
		check_authToken(clientId, this.now.authId, this.now.authToken, callback)
	}
	
	everyone.now.subscribe = function(type, queue_name){
		if (check_session(this)){
			var queueId = type+"-"+queue_name;
			
			log.info(this.now.authId + " subscribe to "+queueId, "nowjs");
			
			if (type == 'amqp')
				amqp_subscribe_queue(queue_name)
			
			var group = nowjs.getGroup(queueId)
			group.addUser(this.user.clientId);
		}
	}

	everyone.now.unsubscribe = function(type, queue_name){
		if (check_session(this)){
			var queueId = type+"-"+queue_name;
			
			log.info(this.now.authId + " unsubscribe from "+queueId, "nowjs");
				
			nowjs.getGroup(queueId).removeUser(this.user.clientId);
		}
	}
	
	everyone.now.publish = function(type, queue_name, message){
		if (check_session(this)){
			var queueId = type+"-"+queue_name;
			
			if (type == 'amqp'){
				var event = build_event(message)
				if (event){
					event.clientId = this.user.clientId
					event.authorId = sessions.check(this.user.clientId)
					var rk = build_rk(event)
					amqp_publish('canopsis.events', rk, event)
				}else{
					log.error('Invalid event.', 'nowjs')
				}
			}else{
				var group = nowjs.getGroup(queueId)
				if (group)
					group.now.on_message(message)
			}
		}
	}
	
	everyone.now.direct = function(authId, message){
		if (check_session(this)){
			var from_authId = sessions.check(this.user.clientId)
			var to_clientIds = sessions.getclientIds(authId)
			
			log.info(this.now.authId + " send direct message to "+authId, "nowjs");
			log.debug(" + from_authId:   "+ from_authId , "nowjs");
			log.debug(" + from_clientId: "+ this.user.clientId , "nowjs");
			
			for (var i in to_clientIds){
				var to_clientId = to_clientIds[i]
				
				log.debug(" + to_clientId: "+ to_clientId , "nowjs");
			
				if (to_clientId)
					nowjs.getClient(to_clientId, function(){
						if (this.now && this.now['on_direct'])
							this.now.on_direct(message)
							log.debug("   + Sended" , "nowjs");
					});
			}
		}
	}
	
	////////////////// Binding events
	nowjs.on("connect", function(){
		var clientId = this.user.clientId
		var authId = this.now.authId 
		if (authId == undefined)
			authId = 'Unknown'
		log.info(authId + " connected ("+clientId+")", "nowjs");
	});

	nowjs.on("disconnect", function(){
		var clientId = this.user.clientId
		log.info(this.now.authId + " disconnected ("+clientId+")", "nowjs");
		sessions.drop(clientId)
	});
	
	callback()
}

var heartbeat = function(){

	// Check nowjs and amqp Queue (Close amqp queue if group is empty)
	nowjs.getGroups(function(groups){
		for (var i in groups){
			var group = groups[i]
			if (group != 'everyone'){
				nowjs.getGroup(group).count(function(count){
					log.debug("Group '"+group+"': "+count+" Client(s)", "heartbeat")
					if (count == 0){
						log.debug(" + Remove group", "heartbeat")
						amqp_unsubscribe_queue(group)
						nowjs.removeGroup(group)
					}
				})
			}
		}
	});

	// Check amqp connection, reconnect if possible
	if (! amqp_connection.readable){
		log.info(" + Try to reconnect to AMQP", "heartbeat")
		init_amqp(function(){
			log.info(" + Re-create AMQP Queues", "amqp")
			for (var i in amqp_queues){
				var queue = amqp_queues[i]
				var short_name = queue.short_name
				delete amqp_queues[i]
				amqp_subscribe_queue(short_name)
			}
		})
	}
}

//####################################################
//#  Stream helper
//####################################################
// TODO: split this file in external files .....

var stream_getComments = function(referer, limit, callback){
	log.debug("getComments for '"+referer+"'", "widget-stream")
	mongodb_find('events_log', { "$and": [{"referer": referer }, {"event_type": "comment"} ]}, { 'limit': limit, 'sort': {"timestamp": -1} }, callback)
}

var stream_countComments = function(referer, callback){
	log.debug("countComments for '"+referer+"'", "widget-stream")
	mongodb_count('events_log', { "$and": [{"referer": referer }, {"event_type": "comment"} ]}, callback)
}
var stream_getHistory= function(limit, tags, tags_op, from, to, callback){
	log.debug("getHistory (tags: "+tags+" (op: "+tags_op+")) "+from+" -> "+to, "widget-stream")
	
	var mfilter = { "$and": [{"state_type": 1 }, {"event_type": {"$ne": "comment"}}]}
	
	if (tags)
		if (tags_op)
			mfilter["$and"].push({"tags": {"$all": tags}})
		else
			mfilter["$and"].push({"tags": {"$in": tags}})
		
	if (from && to)
		mfilter["$and"].push({"timestamp": { "$gte": from, "$lte": to } })

	mongodb_find('events_log', mfilter, { 'limit': limit, 'sort': {"timestamp": -1} }, callback)
}

//####################################################
//#  Stop daemon
//####################################################

var stop = function(){
	log.info("Stop daemon", "main")

	log.info(" + Stop AMQP", "main")
	amqp_connection.end()

	log.info("Bye !", "main")
	process.exit(0)
}

//####################################################
//#  Main Program
//####################################################

process.on('SIGINT', function () {
	stop();
});

process.on('SIGTERM', function () {
	stop();
});

read_config(function(){
	log.debug("Configurations:", "main")
	config.nowjs.debug = (config.nowjs.debug === 'true') || (config.nowjs.debug === 'True')
	
	// Force debug
	//config.nowjs.debug = true

	log.dump(config)
	
	init_mongo(function(){
		init_amqp(function(){
			init_now(function(){
				log.info(" + heartbeat interval: " + config.nowjs.heartbeat + " sec", "main")
				log.info("Initialization completed, Ready for action !", "main")
				setInterval(heartbeat, config.nowjs.heartbeat * 1000);
				
				everyone.now.stream_getComments = stream_getComments;
				everyone.now.stream_getHistory = stream_getHistory;
				everyone.now.stream_countComments = stream_countComments;
			});
		});
	});
});

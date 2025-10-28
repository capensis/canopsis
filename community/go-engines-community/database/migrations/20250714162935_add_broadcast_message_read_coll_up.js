db.broadcast_message_read.createIndex({message: 1, user: 1}, {name: "message_1_user_1"});
db.broadcast_message.createIndex({start: 1, end: 1}, {name: "start_1_end_1"});

db.pbehavior_type.updateMany({type: "active"}, {$set: {color: "#2FAB63"}});
db.pbehavior_type.updateMany({type: "inactive"}, {$set: {color: "#979797"}});
db.pbehavior_type.updateMany({type: "maintenance"}, {$set: {color: "#BF360C"}});
db.pbehavior_type.updateMany({type: "pause"}, {$set: {color: "#5A6D80"}});

db.permission.deleteMany({_id: "api_entitycomment"});
db.role.updateMany({}, {$unset: {"permissions.api_entitycomment": ""}});

db.periodical_alarm.find({
    "v.resolved": null,
    tags: {$ne: null},
    etags: null,
    itags: null,
}).forEach(function (doc) {
    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            etags: doc.tags,
        }
    });
});

db.alarm_tag.createIndex({value: 1}, {name: "value_1", unique: true});
db.alarm_tag.find({type: null}).forEach(function (doc) {
    db.alarm_tag.updateOne({_id: doc._id}, {
        $set: {
            type: 0,
            updated: doc.created,
        }
    });
});

db.periodical_alarm.createIndex({itags: 1, itags_upd: 1}, {name: "itags_1_itags_upd_1"});

if (!db.permission.findOne({_id: "api_alarm_tag"})) {
    db.permission.insertOne({
        _id: "api_alarm_tag",
        crecord_name: "api_alarm_tag",
        crecord_type: "action",
        description: "Alarm tag",
        type: "CRUD"
    });
    db.role.updateMany({"permissions.api_alarm_read": 1}, {
        $set: {
            "permissions.api_alarm_tag": 4
        }
    });
    db.role.updateOne({name: "admin"}, {
        $set: {
            "permissions.api_alarm_tag": 15
        }
    });
}

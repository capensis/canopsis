db.periodical_alarm.find({
    "v.resolved": null,
    $or: [
        {etags: {$ne: null}},
        {itags: {$ne: null}}
    ],
}).forEach(function (doc) {
    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            tags: doc.etags,
        },
        $unset: {
            etags: "",
            itags: "",
            itags_upt: "",
        }
    });
});

db.alarm_tag.dropIndex("value_1");
db.alarm_tag.deleteMany({type: 1});
db.alarm_tag.updateMany({type: 0}, {
    $unset: {
        type: "",
        updated: ""
    },
});

db.periodical_alarm.dropIndex("itags_1_itags_upd_1");

db.permission.deleteOne({_id: "api_alarm_tag"});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_alarm_tag": ""
    }
});

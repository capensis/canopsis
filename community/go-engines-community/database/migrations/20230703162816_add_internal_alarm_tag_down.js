db.periodical_alarm.find({
    "v.resolved": null,
    $or: [
        {external_tags: {$ne: null}},
        {internal_tags: {$ne: null}}
    ],
}).forEach(function (doc) {
    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            tags: doc.external_tags,
        },
        $unset: {
            external_tags: "",
            internal_tags: "",
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

db.periodical_alarm.dropIndex("internal_tags_1_internal_tags_updated_1");

db.permission.deleteOne({_id: "api_alarm_tag"});
db.role.updateMany({}, {
    $unset: {
        "permissions.api_alarm_tag": ""
    }
});

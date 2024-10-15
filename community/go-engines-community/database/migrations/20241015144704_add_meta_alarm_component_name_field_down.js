db.meta_alarm_states.aggregate([
    {$match: {meta_alarm_component_name: {$ne: null}}},
    {
        $lookup: {
            from: "meta_alarm_rules",
            localField: "_id",
            foreignField: "_id",
            as: "rule"
        }
    },
    {$unwind: "$rule"},
    {
        $match: {
            "rule.type": {
                $in: [
                    "attribute",
                    "complex",
                    "timebased",
                    "valuegroup",
                ]
            }
        }
    }
]).forEach(function (doc) {
    db.meta_alarm_states.updateOne({_id: doc._id}, {
        $unset: {
            meta_alarm_component_name: "",
        }
    });
});

db.meta_alarm_states.aggregate([
    {$match: {meta_alarm_component_name: {$ne: null}}},
    {
        $lookup: {
            from: "meta_alarm_rules",
            localField: "_id",
            foreignField: "_id",
            as: "rule"
        }
    },
    {$unwind: "$rule"},
    {
        $match: {
            "rule.type": "corel"
        }
    },
]).forEach(function (doc) {
    db.meta_alarm_states.updateOne({_id: doc._id}, {
        $set: {
            meta_alarm_name: doc.meta_alarm_name + "/" + doc.meta_alarm_component_name,
        },
        $unset: {
            meta_alarm_component_name: ""
        },
    });
});

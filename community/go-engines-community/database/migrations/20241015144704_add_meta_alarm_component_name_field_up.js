db.meta_alarm_states.aggregate([
    {$match: {meta_alarm_component_name: null}},
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
        $set: {
            meta_alarm_component_name: "metaalarm",
        }
    });
});

db.meta_alarm_states.aggregate([
    {$match: {meta_alarm_component_name: null}},
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
    {
        $lookup: {
            from: "default_entities",
            localField: "meta_alarm_name",
            foreignField: "_id",
            as: "entity"
        }
    },
    {$unwind: "$entity"},
    {
        $match: {
            "entity.type": "resource"
        }
    },
]).forEach(function (doc) {
    db.meta_alarm_states.updateOne({_id: doc._id}, {
        $set: {
            meta_alarm_component_name: doc.entity.component,
            meta_alarm_name: doc.entity.name,
        }
    });
});

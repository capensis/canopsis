const ruleTypesToUpdate = [
    "attribute",
    "complex",
    "timebased",
    "valuegroup",
];
const ruleTypeCorel = "corel"
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
    {$match: {"rule.type": {$in: ruleTypesToUpdate}}}
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
    {$match: {"rule.type": ruleTypeCorel}},
    {
        $lookup: {
            from: "default_entities",
            localField: "meta_alarm_name",
            foreignField: "_id",
            as: "entity"
        }
    },
    {$unwind: "$entity"},
    {$match: {"entity.type": "resource"}},
]).forEach(function (doc) {
    db.meta_alarm_states.updateOne({_id: doc._id}, {
        $set: {
            meta_alarm_component_name: doc.entity.component,
            meta_alarm_name: doc.entity.name,
        }
    });
});

(function () {
    let ruleTypes = {};
    db.meta_alarm_rules.find().forEach(function (rule) {
        ruleTypes[rule._id] = rule.type;
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
        {$unwind: {path: "$rule", preserveNullAndEmptyArrays: true}},
        {$match: {rule: null}}
    ]).forEach(function (doc) {
        let ruleType = "";
        for (const ruleID in ruleTypes) {
            if (doc._id.startsWith(ruleID)) {
                ruleType = ruleTypes[ruleID];
                break;
            }
        }

        if (!ruleTypesToUpdate.includes(ruleType)) {
            return;
        }

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
        {$unwind: {path: "$rule", preserveNullAndEmptyArrays: true}},
        {$match: {rule: null}},
        {
            $lookup: {
                from: "default_entities",
                localField: "meta_alarm_name",
                foreignField: "_id",
                as: "entity",
                pipeline: [
                    {$match: {type: "resource"}}
                ]
            }
        },
        {$unwind: "$entity"}
    ]).forEach(function (doc) {
        let ruleType = "";
        for (const ruleID in ruleTypes) {
            if (doc._id.startsWith(ruleID)) {
                ruleType = ruleTypes[ruleID];
                break;
            }
        }

        if (ruleType !== ruleTypeCorel) {
            return;
        }

        db.meta_alarm_states.updateOne({_id: doc._id}, {
            $set: {
                meta_alarm_component_name: doc.entity.component,
                meta_alarm_name: doc.entity.name,
            }
        });
    });
})();

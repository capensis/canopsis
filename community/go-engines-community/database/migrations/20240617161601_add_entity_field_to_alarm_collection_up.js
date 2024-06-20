db.periodical_alarm.aggregate([
    {$match: {e: null}},
    {
        $lookup: {
            from: "default_entities",
            localField: "d",
            foreignField: "_id",
            as: "entity",
        }
    },
    {$unwind: "$entity"},
    {
        $project: {
            entity: 1,
        }
    },
]).forEach(function (doc) {
    db.periodical_alarm.updateOne({_id: doc._id}, {
        $set: {
            e: doc.entity,
        }
    });
});

db.periodical_alarm.createIndex({"e._id": 1}, {name: "e_id_1"});
db.periodical_alarm.createIndex({"e.enabled": 1}, {name: "e_enabled_1"});
db.periodical_alarm.createIndex({"e.type": 1}, {name: "e_type_1"});
db.periodical_alarm.createIndex({"e.connector": 1}, {name: "e_connector_1"});
db.periodical_alarm.createIndex({"e.component": 1}, {name: "e_component_1"});
db.periodical_alarm.createIndex({"e.services": 1}, {name: "e_services_1"});
db.periodical_alarm.createIndex({"e.type": 1}, {
    name: "e_type_service_1",
    partialFilterExpression: {type: {$eq: "service"}}
});

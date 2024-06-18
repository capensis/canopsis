db.periodical_alarm.aggregate([
    {$match: {entity: null}},
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
            entity: doc.entity,
        }
    });
});

db.periodical_alarm.createIndex({"entity._id": 1}, {name: "entity_id_1"});
db.periodical_alarm.createIndex({"entity.enabled": 1}, {name: "entity_enabled_1"});
db.periodical_alarm.createIndex({"entity.type": 1}, {name: "entity_type_1"});
db.periodical_alarm.createIndex({"entity.connector": 1}, {name: "entity_connector_1"});
db.periodical_alarm.createIndex({"entity.component": 1}, {name: "entity_component_1"});
db.periodical_alarm.createIndex({"entity.services": 1}, {name: "entity_services_1"});
db.periodical_alarm.createIndex({"entity.type": 1}, {
    name: "entity_type_service_1",
    partialFilterExpression: {type: {$eq: "service"}}
});

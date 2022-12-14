db.createCollection("declare_ticket_rule");

if (!db.default_rights.findOne({_id: "api_declare_ticket_rule"})) {
    db.default_rights.insertOne({
        _id: "api_declare_ticket_rule",
        crecord_name: "api_declare_ticket_rule",
        crecord_type: "action",
        desc: "Declare ticket rule",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_declare_ticket_rule": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

db.action_scenario.find({"actions.type": "webhook"}).forEach(function (doc) {
    if (!doc.actions) {
        return;
    }

    var set = {};
    var unset = {};
    doc.actions.forEach(function (action, index) {
        if (action.type !== "webhook") {
            return;
        }

        if (action.parameters.retry_count) {
            set["actions." + index + ".parameters.request.retry_count"] = action.parameters.retry_count;
            unset["actions." + index + ".parameters.retry_count"] = "";
        }
        if (action.parameters.retry_delay) {
            set["actions." + index + ".parameters.request.retry_delay"] = action.parameters.retry_delay;
            unset["actions." + index + ".parameters.retry_delay"] = "";
        }
    });

    if (Object.keys(set).length > 0) {
        var update = {};
        update["$set"] = set;
        update["$unset"] = unset;
        db.action_scenario.updateOne({_id: doc._id}, update);
    }
});

db.eventfilter.find({type: "enrichment", external_data: {$ne: null}}).forEach(function (doc) {
    if (!doc.external_data) {
        return;
    }

    var set = {};
    var unset = {};
    Object.keys(doc.external_data).forEach(function (key) {
        if (doc.external_data[key].type !== "api") {
            return;
        }

        if (doc.external_data[key].retry_count) {
            set["external_data." + key + ".request.retry_count"] = doc.external_data[key].retry_count;
            unset["external_data." + key + ".retry_count"] = "";
        }
        if (doc.external_data[key].retry_delay) {
            set["external_data." + key + ".request.retry_delay"] = doc.external_data[key].retry_delay;
            unset["external_data." + key + ".retry_delay"] = "";
        }
    });

    if (Object.keys(set).length > 0) {
        var update = {};
        update["$set"] = set;
        update["$unset"] = unset;
        db.eventfilter.updateOne({_id: doc._id}, update);
    }
});

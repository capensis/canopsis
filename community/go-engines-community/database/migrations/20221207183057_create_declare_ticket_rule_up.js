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
    doc.actions.forEach(function (index, action) {
        if (action.type !== "webhook") {
            return;
        }

        if (action.retry_count) {
            set["action." + index + ".request.retry_count"] = action.retry_count;
            unset["action." + index + ".retry_count"] = "";
        }
        if (action.retry_delay) {
            set["action." + index + ".request.retry_delay"] = action.retry_delay;
            unset["action." + index + ".retry_delay"] = "";
        }
    });

    if (Object.keys(set) > 0) {
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

    if (Object.keys(set) > 0) {
        var update = {};
        update["$set"] = set;
        update["$unset"] = unset;
        db.eventfilter.updateOne({_id: doc._id}, update);
    }
});

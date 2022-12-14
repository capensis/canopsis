db.declare_ticket_rule.drop();

db.default_rights.deleteMany({
    _id: {
        $in: [
            "api_declare_ticket_rule",
        ]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_declare_ticket_rule": "",
    }
});

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

        if (action.parameters && action.parameters.request) {
            if (action.parameters.request.retry_count) {
                set["actions." + index + ".parameters.retry_count"] = action.parameters.request.retry_count;
                unset["actions." + index + ".parameters.request.retry_count"] = "";
            }
            if (action.parameters.request.retry_delay) {
                set["actions." + index + ".parameters.retry_delay"] = action.parameters.request.retry_delay;
                unset["actions." + index + ".parameters.request.retry_delay"] = "";
            }
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

        if (doc.external_data[key].request) {
            if (doc.external_data[key].request.retry_count) {
                set["external_data." + key + ".retry_count"] = doc.external_data[key].request.retry_count;
                unset["external_data." + key + ".request.retry_count"] = "";
            }
            if (doc.external_data[key].request.retry_delay) {
                set["external_data." + key + ".retry_delay"] = doc.external_data[key].request.retry_delay;
                unset["external_data." + key + ".request.retry_delay"] = "";
            }
        }
    });

    if (Object.keys(set).length > 0) {
        var update = {};
        update["$set"] = set;
        update["$unset"] = unset;
        db.eventfilter.updateOne({_id: doc._id}, update);
    }
});

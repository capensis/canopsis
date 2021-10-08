(function () {
    function genID() {
        return UUID().toString().split('"')[1];
    }

    var now = Math.ceil((new Date()).getTime() / 1000);
    var priority = -1;
    db.idle_rule.find().sort({priority: -1}).limit(1).forEach(function (doc) {
        priority = doc.priority
    });

    db.heartbeat.find().forEach(function (doc) {
        var entityPattern = {};
        var oldPattern = doc.pattern;

        if (oldPattern.source_type) {
            entityPattern["type"] = oldPattern.source_type;
            delete oldPattern.source_type;
        }

        if (oldPattern.resource) {
            if (oldPattern.component) {
                entityPattern["_id"] = oldPattern.resource + "/" + oldPattern.component;
            } else {
                entityPattern["name"] = oldPattern.resource;
            }
        } else if (oldPattern.component) {
            entityPattern["_id"] = oldPattern.component;
        } else if (oldPattern.connector_name) {
            if (oldPattern.connector) {
                entityPattern["_id"] = oldPattern.connector + "/" + oldPattern.connector_name;
            } else {
                entityPattern["name"] = oldPattern.connector_name;
            }
        } else if (oldPattern.connector) {
            entityPattern["_id"] = {regex_match: oldPattern.connector + "/*"};
        }

        delete oldPattern.resource;
        delete oldPattern.component;
        delete oldPattern.connector;
        delete oldPattern.connector_name;

        if (oldPattern.event_type === "check") {
            delete oldPattern.event_type;
        }

        if (Object.keys(oldPattern).length > 0) {
            return;
        }

        priority++;
        var val = parseInt(doc.expected_interval);
        var unit = doc.expected_interval.replace(val.toString(), "");
        var seconds = val;
        switch (unit) {
            case "s":
                seconds = val;
                break
            case "m":
                seconds = val * 60;
                break
            case "h":
                seconds = val * 60 * 60;
                break
        }

        db.idle_rule.insertOne({
            _id: genID(),
            type: "entity",
            enabled: true,
            name: doc.name,
            description: doc.description,
            author: doc.author,
            priority: priority,
            created: doc.created ? doc.created : now,
            updated: now,
            entity_patterns: [entityPattern],
            duration: {
                seconds: seconds,
                unit: unit,
            }
        });
    });
})();

function genID() {
    var hex;
    try {
        hex = UUID().hex(); // mongo
    } catch (e) {
        hex = UUID().toString('hex'); // mongosh
    }

    return hex.match(/^(.{8})(.{4})(.{4})(.{4})(.{12})$/).slice(1,6).join('-')
}

db.eventfilter.find({
    "actions.type": {$in:["set_field_from_template", "set_field"]},
    "actions.name": {$in:["Resource", "Component", "Connector", "ConnectorName"]}
}).forEach(function (doc) {
    var now = Math.ceil((new Date()).getTime() / 1000);

    var config = {};

    doc.actions.forEach(function (action) {
        if (action.type === "set_field" || action.type === "set_field_from_template") {
            if (action.name === "Resource") {
                config.resource = action.value
            }

            if (action.name === "Component") {
                config.component = action.value
            }

            if (action.name === "Connector") {
                config.connector = action.value
            }

            if (action.name === "ConnectorName") {
                config.connector_name = action.value
            }
        }
    })

    db.eventfilter.insertOne({
        _id: genID(),
        description: doc.description,
        type: "change_entity",
        patterns: doc.patterns,
        priority: doc.priority,
        enabled: doc.enabled,
        config: config,
        external_data: doc.external_data || null,
        created: now,
        updated: now,
        author: doc.author || null
    });
});

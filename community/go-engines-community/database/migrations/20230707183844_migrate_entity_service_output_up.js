if (db.entity_service_counters.countDocuments() === 0) {
    db.default_entities.find({
        "type": "service"
    }).forEach(function (doc) {
        var outputTemplate = doc.output_template
        if (outputTemplate === undefined) {
            outputTemplate = ""
        }
        db.entity_service_counters.insertOne({
            _id: doc._id,
            output_template: outputTemplate
        });
    })
}

db.default_entities.find({
    "type": "service"
}).forEach(function (doc) {
    outputTemplate = doc.output_template
    if (outputTemplate === undefined) {
        outputTemplate = ""
    }
    db.entity_service_counters.insert({
        _id: doc._id,
        output_template: outputTemplate
    });
})

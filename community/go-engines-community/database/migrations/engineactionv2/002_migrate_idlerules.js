(function () {
    db.idle_rule.find().forEach(function (doc) {
        var updated = false;
        var parameters = doc.operation.parameters;
        if (parameters.message) {
            parameters.output = parameters.message;
            delete parameters.message;
            updated = true;
        }
        if (doc.operation.type === "snooze") {
            parameters.duration = {
                seconds: parameters.duration,
                unit: "s"
            };
            updated = true;
        }

        if (updated) {
            db.idle_rule.updateOne({_id: doc._id}, {
                "$set": {"operation.parameters": parameters},
            });
        }
    });
})();

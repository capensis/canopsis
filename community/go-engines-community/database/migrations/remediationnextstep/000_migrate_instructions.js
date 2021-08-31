(function () {
    db.instruction.find().forEach(function (doc) {
        var set = {};
        var unset = {};

        if (!doc.type) {
            set["type"] = 0;
        }
        if (!doc.timeout_after_execution) {
            set["timeout_after_execution"] = {
                seconds: 60,
                unit: "m"
            };
        }

        if (doc.last_executed_by) {
            unset["last_executed_by"] = "";
        }

        var update = {};
        if (Object.keys(set).length > 0) {
            update["$set"] = set;
        }
        if (Object.keys(unset).length > 0) {
            update["$unset"] = unset;
        }

        if (Object.keys(update).length > 0) {
            db.instruction.updateOne({_id: doc._id}, update);
        }
    });
})();

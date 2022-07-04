(function () {
    var filter = {"v.steps._t": "comment", "v.last_comment": null};
    var updatePipeline = [
        {
            $set: {
                "v.last_comment": {
                    "$arrayElemAt": [
                        {
                            "$filter": {
                                "input": "$v.steps",
                                "cond": {
                                    "$eq": ["$$this._t", "comment"],
                                },
                            }
                        },
                        -1,
                    ],
                }
            }
        },
    ];

    db.periodical_alarm.updateMany(filter, updatePipeline);
    db.resolved_alarms.updateMany(filter, updatePipeline);
})();

(function () {
    db.widgets.find().forEach(function (doc) {
        if (!doc.parameters) {
            return;
        }

        var update = {};
        var oldField = "v.lastComment";
        var newField = "v.last_comment";
        var columnFields = ["widgetColumns", "widgetExportColumns", "widgetGroupColumns"];

        for (var columnField of columnFields) {
            if (doc.parameters[columnField]) {
                for (var i in doc.parameters[columnField]) {
                    var value = doc.parameters[columnField][i].value;
                    if (value.startsWith(oldField)) {
                        update["parameters." + columnField + "." + i + ".value"] = value.replace(oldField, newField);
                    }
                }
            }
        }

        if (doc.parameters.sort && doc.parameters.sort.column && doc.parameters.sort.column.startsWith(oldField)) {
            update["parameters.sort.column"] = doc.parameters.sort.column.replace(oldField, newField);
        }

        if (Object.keys(update).length > 0) {
            db.widgets.updateOne({_id: doc._id}, {$set: update});
        }
    });
})();

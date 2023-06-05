db.widgets.find({"parameters.serviceDependenciesColumns": {$ne: null}}).forEach(function (doc) {
    if (!doc.parameters || !doc.parameters.serviceDependenciesColumns || doc.parameters.serviceDependenciesColumns.length === 0) {
        return;
    }

    var columns = [];
    var updated = false;
    for (var column of doc.parameters.serviceDependenciesColumns) {
        if (!column.value.startsWith("alarm.")) {
            column.value = "entity." + column.value;
            updated = true;
        }
        columns.push(column);
    }

    if (updated) {
        db.widgets.updateOne({_id: doc._id}, {
            $set: {
                "parameters.serviceDependenciesColumns": columns,
            }
        });
    }
});

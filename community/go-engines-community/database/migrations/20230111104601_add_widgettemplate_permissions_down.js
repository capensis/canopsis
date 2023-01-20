db.widget_templates.drop();

db.default_rights.deleteMany({
    _id: {
        $in: ["api_widgettemplate", "models_widgetTemplate"]
    }
});
db.default_rights.updateMany({crecord_type: "role"}, {
    $unset: {
        "rights.api_widgettemplate": "",
        "rights.models_widgetTemplate": "",
    }
});

function migrateEntityColumns(columns) {
    if (!columns) {
        return false;
    }

    for (var column of columns) {
        column.value = "entity." + column.value;
    }

    return true;
}

db.widgets.find().forEach(function (doc) {
    if (!doc.parameters) {
        return;
    }

    var set = {};
    var unset = {};

    switch (doc.type) {
        case "AlarmsList":
        case "Context":
        case "ServiceWeather":
            if (migrateEntityColumns(doc.parameters.serviceDependenciesColumns)) {
                set["parameters.serviceDependenciesColumns"] = doc.parameters.serviceDependenciesColumns;
            }
            break;
        case "Map":
            if (doc.parameters.alarmsColumns) {
                set["parameters.alarms_columns"] = doc.parameters.alarmsColumns;
                unset["parameters.alarmsColumns"] = "";
            }
            if (doc.parameters.entitiesColumns) {
                set["parameters.entities_columns"] = doc.parameters.entitiesColumns;
                unset["parameters.entitiesColumns"] = "";
            }
            break;
    }

    var update = {};
    if (Object.keys(set).length > 0) {
        update["$set"] = set;
    }
    if (Object.keys(unset).length > 0) {
        update["$unset"] = unset;
    }
    if (Object.keys(update).length > 0) {
        db.widgets.updateOne({_id: doc._id}, update);
    }
});

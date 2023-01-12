if (!db.default_rights.findOne({_id: "api_widgettemplate"})) {
    db.default_rights.insertOne({
        _id: "api_widgettemplate",
        crecord_name: "api_widgettemplate",
        crecord_type: "action",
        desc: "Widget templates",
        type: "CRUD"
    });
    db.default_rights.updateOne({_id: "admin"}, {
        $set: {
            "rights.api_widgettemplate": {
                checksum: 15,
                crecord_type: "right"
            }
        }
    });
}

function migrateAlarmColumns(columns) {
    if (!columns) {
        return false;
    }

    var updated = false;
    for (var column of columns) {
        if (column.value === "priority") {
            column.value = "impact_state";
            updated = true;
        } else if (column.value.startsWith("infos")) {
            column.value = column.value.replace("infos", "v.infos");
            updated = true;
        }
    }

    return updated;
}

function migrateEntityColumns(columns) {
    if (!columns) {
        return false;
    }

    var updated = false;
    for (var column of columns) {
        if (column.value.startsWith("entity.")) {
            column.value = column.value.replace("entity.", "");
            updated = true;
        }
    }

    return updated;
}

db.widgets.find().forEach(function (doc) {
    if (!doc.parameters) {
        return;
    }

    var set = {};
    var unset = {};

    switch (doc.type) {
        case "AlarmsList":
            if (migrateAlarmColumns(doc.parameters.widgetColumns)) {
                set["parameters.widgetColumns"] = doc.parameters.widgetColumns;
            }
            if (migrateAlarmColumns(doc.parameters.widgetGroupColumns)) {
                set["parameters.widgetGroupColumns"] = doc.parameters.widgetGroupColumns;
            }
            if (migrateEntityColumns(doc.parameters.serviceDependenciesColumns)) {
                set["parameters.serviceDependenciesColumns"] = doc.parameters.serviceDependenciesColumns;
            }
            break;
        case "Context":
            if (migrateAlarmColumns(doc.parameters.activeAlarmsColumns)) {
                set["parameters.activeAlarmsColumns"] = doc.parameters.activeAlarmsColumns;
            }
            if (migrateAlarmColumns(doc.parameters.resolvedAlarmsColumns)) {
                set["parameters.resolvedAlarmsColumns"] = doc.parameters.resolvedAlarmsColumns;
            }
            if (migrateEntityColumns(doc.parameters.serviceDependenciesColumns)) {
                set["parameters.serviceDependenciesColumns"] = doc.parameters.serviceDependenciesColumns;
            }
            break;
        case "ServiceWeather":
            if (doc.parameters.alarmsList && migrateAlarmColumns(doc.parameters.alarmsList.widgetColumns)) {
                set["parameters.alarmsList.widgetColumns"] = doc.parameters.alarmsList.widgetColumns;
            }
            if (migrateEntityColumns(doc.parameters.serviceDependenciesColumns)) {
                set["parameters.serviceDependenciesColumns"] = doc.parameters.serviceDependenciesColumns;
            }
            break;
        case "Counter":
        case "StatsCalendar":
            if (doc.parameters.alarmsList && migrateAlarmColumns(doc.parameters.alarmsList.widgetColumns)) {
                set["parameters.alarmsList.widgetColumns"] = doc.parameters.alarmsList.widgetColumns;
            }
            break;
        case "Map":
            if (doc.parameters.alarms_columns) {
                migrateAlarmColumns(doc.parameters.alarms_columns)
                set["parameters.alarmsColumns"] = doc.parameters.alarms_columns;
                unset["parameters.alarms_columns"] = "";
            }
            if (doc.parameters.entities_columns) {
                set["parameters.entitiesColumns"] = doc.parameters.entities_columns;
                unset["parameters.entities_columns"] = "";
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

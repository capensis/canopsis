function genID() {
    return UUID().toString().split('"')[1];
}

function isInt(value) {
    return typeof value === "number" || value instanceof NumberLong;
}

function migrateOldMongoQueryForAlarmList(oldMongoQuery) {
    if (!oldMongoQuery || oldMongoQuery === "") {
        return null;
    }
    var query = JSON.parse(oldMongoQuery);
    var newAlarmPattern = [];
    var newPbehaviorPattern = [];
    var newEntityPattern = [];

    if (typeof query === "object" && query && Object.keys(query).length === 1) {
        var highLevelAnd = query["$and"];
        var highLevelOr = query["$or"];

        if (highLevelAnd && Array.isArray(highLevelAnd)) {
            var newGroups = migrateOldGroupForAlarmList(highLevelAnd);
            if (newGroups === null) {
                return;
            }

            if (newGroups[0]) {
                newAlarmPattern.push(newGroups[0])
            }
            if (newGroups[1]) {
                newPbehaviorPattern.push(newGroups[1])
            }
            if (newGroups[2]) {
                newEntityPattern.push(newGroups[2])
            }
        } else if (highLevelOr && Array.isArray(highLevelOr)) {
            for (var oldGroup of highLevelOr) {
                if (typeof oldGroup !== "object" || !oldGroup) {
                    return;
                }

                var and = oldGroup["$and"];
                var newGroups = null;
                if (Object.keys(oldGroup).length === 1 && and && Array.isArray(and)) {
                    newGroups = migrateOldGroupForAlarmList(and);
                } else {
                    newGroups = migrateOldGroupForAlarmList([oldGroup]);
                }

                if (newGroups === null) {
                    return;
                }

                if (newGroups[0]) {
                    newAlarmPattern.push(newGroups[0])
                }
                if (newGroups[1]) {
                    newPbehaviorPattern.push(newGroups[1])
                }
                if (newGroups[2]) {
                    newEntityPattern.push(newGroups[2])
                }
            }
        } else {
            var newGroups = migrateOldGroupForAlarmList([query]);
            if (newGroups === null) {
                return;
            }

            if (newGroups[0]) {
                newAlarmPattern.push(newGroups[0])
            }
            if (newGroups[1]) {
                newPbehaviorPattern.push(newGroups[1])
            }
            if (newGroups[2]) {
                newEntityPattern.push(newGroups[2])
            }
        }
    }

    if (newAlarmPattern.length > 1 && (newPbehaviorPattern.length > 0 || newEntityPattern.length > 0)) {
        return null;
    }
    if (newPbehaviorPattern.length > 1 && (newAlarmPattern.length > 0 || newEntityPattern.length > 0)) {
        return null;
    }
    if (newEntityPattern.length > 1 && (newPbehaviorPattern.length > 0 || newAlarmPattern.length > 0)) {
        return null;
    }

    if (newAlarmPattern.length === 0) {
        newAlarmPattern = null;
    }
    if (newPbehaviorPattern.length === 0) {
        newPbehaviorPattern = null;
    }
    if (newEntityPattern.length === 0) {
        newEntityPattern = null;
    }

    if (!newAlarmPattern && !newPbehaviorPattern && !newEntityPattern) {
        return null;
    }

    return [newAlarmPattern, newPbehaviorPattern, newEntityPattern];
}

function migrateOldMongoQueryForEntityList(oldMongoQuery) {
    if (!oldMongoQuery || oldMongoQuery === "") {
        return null;
    }
    var query = JSON.parse(oldMongoQuery);
    var newEntityPattern = [];

    if (typeof query === "object" && query && Object.keys(query).length === 1) {
        var highLevelAnd = query["$and"];
        var highLevelOr = query["$or"];

        if (highLevelAnd && Array.isArray(highLevelAnd)) {
            var newGroup = migrateOldGroupForEntityList(highLevelAnd);
            if (newGroup === null) {
                return;
            }
            newEntityPattern.push(newGroup)
        } else if (highLevelOr && Array.isArray(highLevelOr)) {
            for (var oldGroup of highLevelOr) {
                if (typeof oldGroup !== "object" || !oldGroup) {
                    return;
                }

                var and = oldGroup["$and"];
                var newGroup = null;
                if (Object.keys(oldGroup).length === 1 && and && Array.isArray(and)) {
                    newGroup = migrateOldGroupForEntityList(and);
                } else {
                    newGroup = migrateOldGroupForEntityList([oldGroup]);
                }

                if (newGroup === null) {
                    return;
                }
                newEntityPattern.push(newGroup)
            }
        } else {
            var newGroup = migrateOldGroupForEntityList([query]);
            if (newGroup === null) {
                return;
            }
            newEntityPattern.push(newGroup)
        }
    }

    if (newEntityPattern.length === 0) {
        return null;
    }

    return newEntityPattern;
}

function migrateOldMongoQueryForWeather(oldMongoQuery) {
    if (!oldMongoQuery || oldMongoQuery === "") {
        return null;
    }
    var query = JSON.parse(oldMongoQuery);
    var newWeatherPattern = [];
    var newEntityPattern = [];

    if (typeof query === "object" && query && Object.keys(query).length === 1) {
        var highLevelAnd = query["$and"];
        var highLevelOr = query["$or"];

        if (highLevelAnd && Array.isArray(highLevelAnd)) {
            var newGroups = migrateOldGroupForWeather(highLevelAnd);
            if (newGroups === null) {
                return;
            }

            if (newGroups[0]) {
                newEntityPattern.push(newGroups[0])
            }
            if (newGroups[1]) {
                newWeatherPattern.push(newGroups[1])
            }
        } else if (highLevelOr && Array.isArray(highLevelOr)) {
            for (var oldGroup of highLevelOr) {
                if (typeof oldGroup !== "object" || !oldGroup) {
                    return;
                }

                var and = oldGroup["$and"];
                var newGroups = null;
                if (Object.keys(oldGroup).length === 1 && and && Array.isArray(and)) {
                    newGroups = migrateOldGroupForWeather(and);
                } else {
                    newGroups = migrateOldGroupForWeather([oldGroup]);
                }

                if (newGroups === null) {
                    return;
                }

                if (newGroups[0]) {
                    newEntityPattern.push(newGroups[0])
                }
                if (newGroups[1]) {
                    newWeatherPattern.push(newGroups[1])
                }
            }
        } else {
            var newGroups = migrateOldGroupForWeather([query]);
            if (newGroups === null) {
                return;
            }

            if (newGroups[0]) {
                newEntityPattern.push(newGroups[0])
            }
            if (newGroups[1]) {
                newWeatherPattern.push(newGroups[1])
            }
        }
    }

    if (newEntityPattern.length > 1 && newWeatherPattern.length > 0) {
        return null;
    }
    if (newWeatherPattern.length > 1 && newEntityPattern.length > 0) {
        return null;
    }

    if (newEntityPattern.length === 0) {
        newEntityPattern = null;
    }
    if (newWeatherPattern.length === 0) {
        newWeatherPattern = null;
    }

    if (!newEntityPattern && !newWeatherPattern) {
        return null;
    }

    return [newEntityPattern, newWeatherPattern];
}

function migrateOldGroupForAlarmList(oldGroup) {
    var newAlarmGroup = [];
    var newPbehaviorGroup = [];
    var newEntityGroup = [];

    for (var oldCond of oldGroup) {
        if (typeof oldCond !== "object" || !oldCond) {
            return null;
        }

        for (var field of Object.keys(oldCond)) {
            var value = oldCond[field];
            var strCond = null;
            if (typeof value === "string") {
                strCond = {
                    type: "eq",
                    value: value,
                };
            } else if (typeof value === "object" && value && Object.keys(value).length === 1) {
                if (value["$regex"] && typeof value["$regex"] === "string") {
                    strCond = {
                        type: "regexp",
                        value: value["$regex"],
                    };
                } else if (value["$in"] && Array.isArray(value["$in"]) && value["$in"].length > 0 && typeof value["$in"][0] === "string") {
                    strCond = {
                        type: "is_one_of",
                        value: value["$in"],
                    };
                } else if (value["$ne"] && typeof value["$ne"] === "string") {
                    strCond = {
                        type: "neq",
                        value: value["$ne"],
                    };
                }
            }

            var intCond = null;
            if (isInt(value)) {
                intCond = {
                    type: "eq",
                    value: value,
                };
            } else if (typeof value === "object" && value && Object.keys(value).length === 1) {
                if (value["$ne"] && isInt(value["$ne"])) {
                    intCond = {
                        type: "neq",
                        value: value["$ne"],
                    };
                } else if (value["$gt"] && isInt(value["$gt"])) {
                    intCond = {
                        type: "gt",
                        value: value["$gt"],
                    };
                } else if (value["$lt"] && isInt(value["$lt"])) {
                    intCond = {
                        type: "lt",
                        value: value["$lt"],
                    };
                } else if (value["$gte"] && isInt(value["$gte"])) {
                    intCond = {
                        type: "gt",
                        value: value["$gte"] - 1,
                    };
                } else if (value["$lte"] && isInt(value["$lte"])) {
                    intCond = {
                        type: "lt",
                        value: value["$lte"] + 1,
                    };
                }
            }

            switch (field) {
                case "v.ack":
                case "v.ticket":
                case "v.canceled":
                case "v.snooze":
                case "v.activation_date":
                    if (typeof value === "object" && value && typeof value["$exists"] === "boolean") {
                        newAlarmGroup.push({
                            field: field,
                            cond: {
                                type: "exist",
                                value: value["$exists"],
                            },
                        });
                    } else {
                        return null;
                    }
                    break;
                case "v.connector":
                case "v.connector_name":
                case "v.resource":
                case "v.component":
                case "v.display_name":
                case "v.output":
                case "v.long_output":
                case "v.initial_output":
                case "v.initial_long_output":
                case "v.ack.a":
                case "v.ack.m":
                case "v.ack.initiator":
                    if (strCond === null) {
                        return null;
                    }
                    newAlarmGroup.push({
                        field: field,
                        cond: strCond,
                    });
                    break;
                case "connector":
                case "connector_name":
                case "resource":
                case "component":
                    if (strCond === null) {
                        return null;
                    }
                    newAlarmGroup.push({
                        field: "v." + field,
                        cond: strCond,
                    });
                    break;
                case "v.total_state_changes":
                    if (intCond === null) {
                        return null;
                    }
                    newAlarmGroup.push({
                        field: field,
                        cond: intCond,
                    });
                    break;
                case "v.state.val":
                case "v.status.val":
                    if (isInt(value)) {
                        newAlarmGroup.push({
                            field: field,
                            cond: {
                                type: "eq",
                                value: value,
                            },
                        });
                    } else if (typeof value === "object" && value && isInt(value["$ne"])) {
                        newAlarmGroup.push({
                            field: field,
                            cond: {
                                type: "neq",
                                value: value["$ne"],
                            },
                        });
                    } else {
                        return null;
                    }
                    break;
                case "v.pbehavior_info":
                    if (typeof value === "object" && value && typeof value["$exists"] === "boolean") {
                        if (value["$exists"]) {
                            newPbehaviorGroup.push({
                                field: "pbehavior_info.canonical_type",
                                cond: {
                                    type: "neq",
                                    value: "active",
                                },
                            });
                        } else {
                            newPbehaviorGroup.push({
                                field: "pbehavior_info.canonical_type",
                                cond: {
                                    type: "eq",
                                    value: "active"
                                },
                            });
                        }
                    } else {
                        return null;
                    }
                    break;
                case "v.pbehavior_info.name":
                    if (typeof value === "string") {
                        var pbehavior = db.pbehavior.findOne({name: value});
                        if (pbehavior) {
                            newPbehaviorGroup.push({
                                field: "pbehavior_info._id",
                                cond: {
                                    type: "eq",
                                    value: pbehavior._id,
                                },
                            });
                        }
                    } else {
                        return null;
                    }
                    break;
                default:
                    if (field.startsWith("entity.")) {
                        var entityField = field.replace("entity.", "");

                        switch (entityField) {
                            case "_id":
                            case "name":
                            case "type":
                            case "component":
                                if (strCond === null) {
                                    return null;
                                }
                                newEntityGroup.push({
                                    field: entityField,
                                    cond: strCond,
                                });
                                break;
                            default:
                                if (entityField.startsWith("infos.") && entityField.endsWith(".value")) {
                                    var info = entityField.replace(".value", "");
                                    if (strCond !== null) {
                                        newEntityGroup.push({
                                            field: info,
                                            field_type: "string",
                                            cond: strCond,
                                        });
                                    } else if (value === null) {
                                        newEntityGroup.push({
                                            field: info,
                                            cond: {
                                                type: "exist",
                                                value: false
                                            },
                                        });
                                    } else {
                                        return null;
                                    }
                                } else {
                                    return null;
                                }
                        }
                    } else if (field.startsWith("infos.") && field.endsWith(".value")) {
                        var info = field.replace(".value", "");
                        if (strCond !== null) {
                            newEntityGroup.push({
                                field: info,
                                field_type: "string",
                                cond: strCond,
                            });
                        } else if (value === null) {
                            newEntityGroup.push({
                                field: info,
                                cond: {
                                    type: "exist",
                                    value: false
                                },
                            });
                        } else {
                            return null;
                        }
                    } else if (field.startsWith("v.infos.*.") && strCond !== null) {
                        var info = field.replace("v.infos.*.", "");
                        newAlarmGroup.push({
                            field: "v.infos." + info,
                            field_type: "string",
                            cond: strCond,
                        });
                    } else {
                        return null;
                    }
            }
        }
    }

    if (newAlarmGroup.length === 0) {
        newAlarmGroup = null;
    }
    if (newPbehaviorGroup.length === 0) {
        newPbehaviorGroup = null;
    }
    if (newEntityGroup.length === 0) {
        newEntityGroup = null;
    }

    return [newAlarmGroup, newPbehaviorGroup, newEntityGroup];
}

function migrateOldGroupForEntityList(oldGroup) {
    var newEntityGroup = [];

    for (var oldCond of oldGroup) {
        if (typeof oldCond !== "object" || !oldCond) {
            return null;
        }

        for (var field of Object.keys(oldCond)) {
            var value = oldCond[field];
            var strCond = null;
            if (typeof value === "string") {
                strCond = {
                    type: "eq",
                    value: value,
                };
            } else if (typeof value === "object" && value) {
                if (value["$regex"] && typeof value["$regex"] === "string") {
                    strCond = {
                        type: "regexp",
                        value: value["$regex"],
                    };
                } else if (value["$in"] && Array.isArray(value["$in"]) && value["$in"].length > 0 && typeof value["$in"][0] === "string") {
                    strCond = {
                        type: "is_one_of",
                        value: value["$in"],
                    };
                } else if (value["$ne"] && typeof value["$ne"] === "string") {
                    strCond = {
                        type: "neq",
                        value: value["$ne"],
                    };
                }
            }

            switch (field) {
                case "_id":
                case "name":
                case "type":
                case "component":
                    if (strCond === null) {
                        return null;
                    }
                    newEntityGroup.push({
                        field: field,
                        cond: strCond,
                    });
                    break;
                default:
                    if (field.startsWith("infos.") && field.endsWith(".value")) {
                        var info = field.replace(".value", "");
                        if (strCond !== null) {
                            newEntityGroup.push({
                                field: info,
                                field_type: "string",
                                cond: strCond,
                            });
                        } else if (value === null) {
                            newEntityGroup.push({
                                field: info,
                                cond: {
                                    type: "exist",
                                    value: false
                                },
                            });
                        } else {
                            return null;
                        }
                    } else {
                        return null;
                    }
            }
        }
    }

    if (newEntityGroup.length === 0) {
        return null;
    }

    return newEntityGroup;
}

function migrateOldGroupForWeather(oldGroup) {
    var newEntityGroup = [];
    var newWeatherGroup = [];

    for (var oldCond of oldGroup) {
        if (typeof oldCond !== "object" || !oldCond) {
            return null;
        }

        for (var field of Object.keys(oldCond)) {
            var value = oldCond[field];
            var strCond = null;
            if (typeof value === "string") {
                strCond = {
                    type: "eq",
                    value: value,
                };
            } else if (typeof value === "object" && value) {
                if (value["$regex"] && typeof value["$regex"] === "string") {
                    strCond = {
                        type: "regexp",
                        value: value["$regex"],
                    };
                } else if (value["$in"] && Array.isArray(value["$in"]) && value["$in"].length > 0 && typeof value["$in"][0] === "string") {
                    strCond = {
                        type: "is_one_of",
                        value: value["$in"],
                    };
                } else if (value["$ne"] && typeof value["$ne"] === "string") {
                    strCond = {
                        type: "neq",
                        value: value["$ne"],
                    };
                }
            }

            switch (field) {
                case "_id":
                case "name":
                case "type":
                case "component":
                    if (strCond === null) {
                        return null;
                    }
                    newEntityGroup.push({
                        field: field,
                        cond: strCond,
                    });
                    break;
                case "icon":
                case "secondary_icon":
                    if (strCond === null) {
                        return null;
                    }
                    newEntityGroup.push({
                        field: field,
                        cond: strCond,
                    });
                    break;
                case "is_grey":
                    if (typeof value === "boolean") {
                        newEntityGroup.push({
                            field: field,
                            cond: {
                                type: "eq",
                                value: value,
                            }
                        });
                    } else {
                        return null;
                    }
                    break;
                default:
                    if (field.startsWith("infos.") && field.endsWith(".value")) {
                        var info = field.replace(".value", "");
                        if (strCond !== null) {
                            newEntityGroup.push({
                                field: info,
                                field_type: "string",
                                cond: strCond,
                            });
                        } else if (value === null) {
                            newEntityGroup.push({
                                field: info,
                                cond: {
                                    type: "exist",
                                    value: false
                                },
                            });
                        } else {
                            return null;
                        }
                    } else {
                        return null;
                    }
            }
        }
    }

    if (newEntityGroup.length === 0) {
        newEntityGroup = null;
    }
    if (newWeatherGroup.length === 0) {
        newWeatherGroup = null;
    }

    return [newEntityGroup, newWeatherGroup];
}

function migrateOldFilter(widget, oldFilter) {
    var newFilter = {
        _id: genID(),
        widget: widget._id,
        title: oldFilter.title,
        old_mongo_query: oldFilter.filter,
    };
    switch (widget.type) {
        case "AlarmsList":
        case "Counter":
        case "StatsCalendar":
            var patterns = migrateOldMongoQueryForAlarmList(oldFilter.filter);
            if (patterns) {
                if (patterns[0]) {
                    newFilter.alarm_pattern = patterns[0];
                }
                if (patterns[1]) {
                    newFilter.pbehavior_pattern = patterns[1];
                }
                if (patterns[2]) {
                    newFilter.entity_pattern = patterns[2];
                }
            }
            break;
        case "Context":
            var entityPattern = migrateOldMongoQueryForEntityList(oldFilter.filter);
            if (entityPattern) {
                newFilter.entity_pattern = entityPattern;
            }
            break;
        case "ServiceWeather":
            var patterns = migrateOldMongoQueryForWeather(oldFilter.filter);
            if (patterns) {
                if (patterns[0]) {
                    newFilter.entity_pattern = patterns[0];
                }
                if (patterns[1]) {
                    newFilter.weather_service_pattern = patterns[1];
                }
            }
            break;
    }

    return newFilter;
}

function migrateOldMainFilter(widget, oldMainFilter) {
    if (Array.isArray(oldMainFilter)) {
        if (oldMainFilter.length === 0) {
            return null;
        }
        var and = [];
        var mergedTitle = "";
        oldMainFilter.forEach(function (oldFilter, i) {
            if (i > 0) {
                mergedTitle += " and ";
            }
            mergedTitle += oldFilter.title;
            and.push(JSON.parse(oldFilter.filter));
        });
        var mergedOldMainFilter = {
            title: mergedTitle,
            filter: JSON.stringify({$and: and}),
        };
        return migrateOldFilter(widget, mergedOldMainFilter);
    }
    if (!oldMainFilter || !oldMainFilter.filter) {
        return null
    }
    return migrateOldFilter(widget, oldMainFilter);
}

db.widgets.find({
        $or: [
            {"parameters.mainFilter": {$ne: null}},
            {"parameters.viewFilters": {$ne: null}}
        ]
    }
).forEach(function (widget) {
    var now = Math.ceil((new Date()).getTime() / 1000);
    var author = widget.author;
    if (!author) {
        author = "root";
    }
    var created = widget.created;
    if (!created) {
        created = now;
    }
    var updated = widget.updated;
    if (!updated) {
        updated = now;
    }

    var mainFilter = widget.parameters.mainFilter;
    var newMainFilter = null;
    var newFilters = [];

    if (widget.parameters.viewFilters) {
        widget.parameters.viewFilters.forEach(function (filter, fi) {
            var newFilter = migrateOldFilter(widget, filter);
            newFilter.is_private = false;
            newFilter.position = fi;
            newFilter.author = author;
            newFilter.created = created;
            newFilter.updated = updated;
            newFilters.push(newFilter);

            if (mainFilter && mainFilter.title === newFilter.title && mainFilter.filter === newFilter.old_mongo_query) {
                newMainFilter = newFilter._id;
            }
        });
    }
    if (mainFilter && !newMainFilter) {
        var newFilter = migrateOldMainFilter(widget, mainFilter);
        if (newFilter) {
            newFilter.is_private = false;
            newFilter.position = newFilters.length;
            newFilter.author = author;
            newFilter.created = created;
            newFilter.updated = updated;
            newFilters.push(newFilter);
            newMainFilter = newFilter._id;
        }
    }

    db.widgets.updateOne({_id: widget._id}, {
        $set: {"parameters.mainFilter": newMainFilter},
        $unset: {"parameters.viewFilters": ""},
    });
    if (newFilters.length > 0) {
        db.widget_filters.insertMany(newFilters);
    }
});

db.userpreferences.aggregate([
    {$match: {widget: {$ne: null}}},
    {$sort: {updated: -1, _id: 1}},
    {
        $group: {
            _id: {
                widget: "$widget",
                user: "$user",
            },
            userPref: {"$first": "$$ROOT"},
            extraIds: {"$push": "$_id"}
        }
    },
    {
        $addFields: {
            extraIds: {
                $filter: {
                    input: "$extraIds",
                    cond: {$ne: ["$$this", "$userPref._id"]}
                }
            }
        }
    },
    {
        $lookup: {
            from: "widgets",
            localField: "userPref.widget",
            foreignField: "_id",
            as: "widget",
        }
    },
    {$unwind: {path: "$widget", preserveNullAndEmptyArrays: true}},
    {
        $lookup: {
            from: "widget_filters",
            localField: "widget._id",
            foreignField: "widget",
            as: "widgetFilters",
        }
    },
]).forEach(function (doc) {
    if (doc.extraIds && doc.extraIds.length > 0) {
        db.userpreferences.deleteMany({_id: {$in: doc.extraIds}});
    }

    var userPref = doc.userPref;
    var widget = doc.widget;
    var widgetFilters = doc.widgetFilters;
    if (!widget) {
        db.userpreferences.deleteOne({_id: userPref._id});
        return;
    }
    if (!userPref.content) {
        return;
    }

    var updated = userPref.updated;
    if (!updated) {
        updated = Math.ceil((new Date()).getTime() / 1000);
    }

    var mainFilter = userPref.content.mainFilter;
    var viewFilters = userPref.content.viewFilters;
    if (!mainFilter && !viewFilters) {
        return;
    }

    var newFilters = [];
    var newMainFilter = null;

    if (viewFilters) {
        viewFilters.forEach(function (filter, fi) {
            var newFilter = migrateOldFilter(widget, filter);
            newFilter.is_private = true;
            newFilter.author = userPref.user;
            newFilter.position = fi;
            newFilter.created = updated;
            newFilter.updated = updated;
            newFilters.push(newFilter);

            if (mainFilter && mainFilter.title === newFilter.title && mainFilter.filter === newFilter.old_mongo_query) {
                newMainFilter = newFilter._id;
            }
        });
    }

    if (mainFilter && !newMainFilter && widgetFilters) {
        widgetFilters.forEach(function (widgetFilter) {
            if (mainFilter.title === widgetFilter.title && mainFilter.filter === widgetFilter.old_mongo_query) {
                newMainFilter = widgetFilter._id;
            }
        });
    }
    if (mainFilter && !newMainFilter) {
        var newFilter = migrateOldMainFilter(widget, mainFilter);
        if (newFilter) {
            newFilter.is_private = true;
            newFilter.author = userPref.user;
            newFilter.position = newFilters.length;
            newFilter.created = updated;
            newFilter.updated = updated;
            newFilters.push(newFilter);
            newMainFilter = newFilter._id;
        }
    }

    db.userpreferences.updateOne({_id: userPref._id}, {
        $set: {"content.mainFilter": newMainFilter},
        $unset: {"content.viewFilters": ""}
    });

    if (newFilters.length > 0) {
        db.widget_filters.insertMany(newFilters);
    }
});

db.widget_filters.createIndex({widget: 1}, {name: "widget_1"});

function migrateOldEntityPatterns(oldEntityPatterns, forbiddenFields) {
    if (!oldEntityPatterns) {
        return null;
    }

    var newPattern = [];
    for (var group of oldEntityPatterns) {
        var newGroup = [];

        for (var field of Object.keys(group)) {
            if (forbiddenFields && forbiddenFields[field]) {
                return null;
            }

            var value = group[field];
            switch (field) {
                case "_id":
                case "name":
                case "type":
                case "category":
                case "component":
                    var cond = migrateOldStringPattern(value);
                    if (!cond) {
                        return null;
                    }
                    newGroup.push({
                        field: field,
                        cond: cond,
                    });
                    break;
                case "last_event_date":
                    return null;
                case "infos":
                    for (var infoKey of Object.keys(value)) {
                        var info = value[infoKey];

                        if (info === null) {
                            newGroup.push({
                                field: field + "." + infoKey,
                                cond: {
                                    type: "exist",
                                    value: false,
                                },
                            });

                            continue
                        }

                        if (info.description || info.name) {
                            return null;
                        }

                        if (info.value) {
                            var cond = migrateOldStringPattern(info.value);
                            if (!cond) {
                                return null;
                            }
                            newGroup.push({
                                field: field + "." + infoKey,
                                field_type: "string",
                                cond: cond,
                            });
                        }
                    }
                    break;
                case "component_infos":
                    for (var componentInfoKey of Object.keys(value)) {
                        var componentInfo = value[componentInfoKey];

                        if (componentInfo === null) {
                            newGroup.push({
                                field: field + "." + componentInfoKey,
                                cond: {
                                    type: "exist",
                                    value: false,
                                },
                            });

                            continue
                        }

                        if (componentInfo.description || componentInfo.name) {
                            return null;
                        }

                        if (componentInfo.value) {
                            var cond = migrateOldStringPattern(componentInfo.value);
                            if (!cond) {
                                return null;
                            }
                            newGroup.push({
                                field: field + "." + componentInfoKey,
                                field_type: "string",
                                cond: cond,
                            });
                        }
                    }
                    break;
                default:
                    return null;
            }
        }

        if (newGroup.length > 0) {
            newPattern.push(newGroup)
        }
    }

    if (newPattern.length > 0) {
        return newPattern;
    }

    return null;
}

function migrateOldAlarmPatterns(oldAlarmPatterns) {
    if (!oldAlarmPatterns) {
        return null;
    }

    var newPattern = [];
    for (var group of oldAlarmPatterns) {
        var newGroup = [];

        for (var field of Object.keys(group)) {
            var value = group[field];
            switch (field) {
                case "t":
                    var creationDateCond = migrateOldTimePattern(value);
                    if (!creationDateCond) {
                        return null;
                    }
                    newGroup.push({
                        field: "v.creation_date",
                        cond: creationDateCond,
                    });
                    break;
                case "v":
                    for (var vField of Object.keys(value)) {
                        var vValue = value[vField];
                        var newField = field + "." + vField;

                        switch (vField) {
                            case "ack":
                                var ackCond = migrateOldAlarmStepPattern(vValue, newField, {t: true, a: true, m: true, initiator: true});
                                if (!ackCond) {
                                    return null;
                                }

                                newGroup.push(ackCond);
                                break;
                            case "canceled":
                            case "ticket":
                            case "snooze":
                                var stepCond = migrateOldAlarmStepPattern(vValue, newField, {});
                                if (!stepCond) {
                                    return null;
                                }

                                newGroup.push(stepCond);
                                break;
                            case "state":
                            case "status":
                                var sCond = migrateOldStateAndStatusAlarmStepPattern(vValue, newField)
                                if (!sCond) {
                                    return null;
                                }

                                newGroup.push(sCond);
                                break;
                            case "creation_date":
                                var timeCond = migrateOldTimePattern(vValue)
                                if (!timeCond) {
                                    return null;
                                }

                                newGroup.push({
                                    field: newField,
                                    cond: timeCond,
                                });
                                break;
                            case "activation_date":
                                if (vValue === null) {
                                    newGroup.push({
                                        field: newField,
                                        cond: {
                                            type: "exist",
                                            value: false,
                                        },
                                    });
                                } else {
                                    var timeCond = migrateOldTimePattern(vValue)
                                    if (!timeCond) {
                                        return null;
                                    }
                                    newGroup.push({
                                        field: newField,
                                        cond: timeCond,
                                    });
                                }
                                break;
                            case "last_update_date":
                            case "last_event_date":
                            case "resolved":
                                return null;
                            case "resource":
                            case "component":
                            case "connector":
                            case "connector_name":
                            case "display_name":
                            case "output":
                            case "long_output":
                            case "initial_output":
                            case "initial_long_output":
                                var cond = migrateOldStringPattern(vValue);
                                if (!cond) {
                                    return null;
                                }
                                newGroup.push({
                                    field: newField,
                                    cond: cond,
                                });
                                break;
                            case "total_state_changes":
                                var cond = migrateOldIntPattern(vValue);
                                if (!cond) {
                                    return null;
                                }
                                newGroup.push({
                                    field: newField,
                                    cond: cond,
                                });
                                break;
                            default:
                                return null;
                        }
                    }
                    break;
                default:
                    return null;
            }
        }

        if (newGroup.length > 0) {
            newPattern.push(newGroup)
        }
    }

    if (newPattern.length > 0) {
        return newPattern;
    }

    return null;
}

function migrateOldEventPatterns(oldEventPatterns) {
    if (!oldEventPatterns) {
        return null;
    }

    var newPattern = [];
    var newEntityPattern = null;
    for (var group of oldEventPatterns) {
        var newGroup = [];

        for (var field of Object.keys(group)) {
            var value = group[field];
            switch (field) {
                case "connector":
                case "connector_name":
                case "component":
                case "resource":
                case "output":
                case "long_output":
                case "event_type":
                case "source_type":
                    var cond = migrateOldStringPattern(value);
                    if (!cond) {
                        return null;
                    }
                    newGroup.push({
                        field: field,
                        cond: cond,
                    });
                    break;
                case "current_entity":
                    if (oldEventPatterns.length > 1) {
                        return null;
                    }

                    newEntityPattern = migrateOldEntityPatterns([value])
                    if (!newEntityPattern) {
                        return null;
                    }

                    break;
                case "state":
                    var cond = migrateOldIntPattern(value);
                    if (!cond) {
                        return null;
                    }
                    newGroup.push({
                        field: field,
                        cond: cond,
                    });
                    break;
                case "_id":
                case "perf_data":
                case "status":
                case "timestamp":
                case "author":
                case "routing_key":
                case "ack_resources":
                case "duration":
                case "ticket":
                case "stat_name":
                case "debug":
                    return null;
                default:
                    var newField = "extra." + field;
                    if (value === null) {
                        newGroup.push({
                            field: newField,
                            cond: {
                                type: "exist",
                                value: false,
                            },
                        });
                    } else if (isInt(value)) {
                        newGroup.push({
                            field: newField,
                            field_type: "int",
                            cond: {
                                type: "eq",
                                value: value,
                            },
                        });
                    } else if (typeof value === "string") {
                        newGroup.push({
                            field: newField,
                            field_type: "string",
                            cond: {
                                type: "eq",
                                value: value,
                            },
                        });
                    } else if (typeof value === "object") {
                        var intCond = migrateOldIntPattern(value);
                        var strCond = migrateOldStringPattern(value);
                        if (intCond) {
                            newGroup.push({
                                field: newField,
                                field_type: "int",
                                cond: intCond,
                            });
                        } else if (strCond) {
                            newGroup.push({
                                field: newField,
                                field_type: "string",
                                cond: strCond,
                            });
                        } else if (typeof value === "object" && (value.has_every || value.has_one_of || value.has_not || value.is_empty)) {
                            if (value.has_every) {
                                if (isStringArray(value.has_every)) {
                                    newGroup.push({
                                        field: newField,
                                        field_type: "string_array",
                                        cond: {
                                            type: "has_every",
                                            value: value.has_every,
                                        },
                                    });
                                } else {
                                    return null;
                                }
                            }
                            if (value.has_one_of) {
                                if (isStringArray(value.has_one_of)) {
                                    newGroup.push({
                                        field: newField,
                                        field_type: "string_array",
                                        cond: {
                                            type: "has_one_of",
                                            value: value.has_one_of,
                                        },
                                    });
                                } else {
                                    return null;
                                }
                            }
                            if (value.has_not) {
                                if (isStringArray(value.has_not)) {
                                    newGroup.push({
                                        field: newField,
                                        field_type: "string_array",
                                        cond: {
                                            type: "has_not",
                                            value: value.has_not,
                                        },
                                    });
                                } else {
                                    return null;
                                }
                            }
                            if (value.is_empty) {
                                if (typeof value.is_empty === "boolean") {
                                    newGroup.push({
                                        field: newField,
                                        field_type: "bool",
                                        cond: {
                                            type: "is_empty",
                                            value: value.is_empty,
                                        },
                                    });
                                } else {
                                    return null;
                                }
                            }
                        } else {
                            return null;
                        }
                    } else {
                        return null;
                    }
            }
        }

        if (newGroup.length > 0) {
            newPattern.push(newGroup)
        }
    }

    if (newPattern.length > 0) {
        return [newPattern, newEntityPattern];
    }

    return [null, newEntityPattern];
}

function migrateOldStringPattern(oldStringPattern) {
    if (oldStringPattern.regex_match) {
        return {
            type: "regexp",
            value: oldStringPattern.regex_match,
        };
    }

    if (typeof oldStringPattern === "string") {
        return {
            type: "eq",
            value: oldStringPattern,
        };
    }

    return null;
}

function migrateOldTimePattern(oldTimePattern) {
    var from = null;
    var to = null
    if (oldTimePattern[">="]) {
        from = oldTimePattern[">="];
    }
    if (oldTimePattern["<="]) {
        to = oldTimePattern["<="];
    }
    if (oldTimePattern[">"]) {
        from = oldTimePattern[">"];
    }
    if (oldTimePattern["<"]) {
        to = oldTimePattern["<"];
    }

    if (from && to) {
        return {
            type: "absolute_time",
            value: {
                from: from,
                to: to,
            },
        };
    }

    return null;
}

function migrateOldAlarmStepPattern(oldAlarmStepPattern, stepField, allowedFields) {
    if (oldAlarmStepPattern === null) {
        return {
            field: stepField,
            cond: {
                type: "exist",
                value: false
            }
        };
    }

    var res = null;
    for (var field of Object.keys(oldAlarmStepPattern)) {
        var value = oldAlarmStepPattern[field];

        switch (field) {
            case "t":
                if (!allowedFields || !allowedFields[field]) {
                    return null;
                }

                var cond = migrateOldTimePattern(value);
                if (!cond) {
                    return null;
                }
                res = {
                    field: stepField + "." + field,
                    cond: cond,
                };
                break;
            case "a":
            case "m":
            case "initiator":
                if (!allowedFields || !allowedFields[field]) {
                    return null;
                }

                var cond = migrateOldStringPattern(value);
                if (!cond) {
                    return null;
                }
                res = {
                    field: stepField + "." + field,
                    cond: cond,
                };
                break;
            default:
                return null;
        }
    }

    return res;
}

function migrateOldStateAndStatusAlarmStepPattern(oldAlarmStepPattern, stepField) {
    if (oldAlarmStepPattern === null) {
        return null;
    }

    var res = null;
    for (var field of Object.keys(oldAlarmStepPattern)) {
        switch (field) {
            case "val":
                var cond = migrateOldIntPattern(oldAlarmStepPattern.val);
                if (!cond) {
                    return null;
                }
                res = {
                    field: stepField + "." + field,
                    cond: cond,
                };
                break;
            default:
                return null;
        }
    }

    return res;
}

function migrateOldIntPattern(oldIntPattern) {
    if (oldIntPattern["<="]) {
        return {
            type: "lt",
            value: oldIntPattern["<="] + 1,
        };
    }
    if (oldIntPattern[">="]) {
        return {
            type: "gt",
            value: oldIntPattern[">="] - 1,
        };
    }
    if (oldIntPattern["<"]) {
        return {
            type: "lt",
            value: oldIntPattern["<"],
        };
    }
    if (oldIntPattern[">"]) {
        return {
            type: "gt",
            value: oldIntPattern[">"],
        };
    }

    if (isInt(oldIntPattern)) {
        return {
            type: "eq",
            value: oldIntPattern,
        };
    }

    return null;
}

function isStringArray(value) {
    if (Array.isArray(value)) {
        if (value.length === 0) {
            return false;
        }
        for (var item of value) {
            if (typeof item !== "string") {
                return false;
            }
        }
    }

    return true;
}

db.default_entities.updateMany({type: "service"}, {
    $rename: {
        entity_patterns: "old_entity_patterns",
    }
});
db.default_entities.find({type: "service"}).forEach(function (doc) {
    var forbiddenFields = {
        "connector": true,
        "component_infos": true,
    };

    if (doc.old_entity_patterns) {
        var newPattern = migrateOldEntityPatterns(doc.old_entity_patterns, forbiddenFields);
        if (newPattern) {
            db.default_entities.updateOne({_id: doc._id}, {
                $set: {entity_pattern: newPattern},
            });
        }
    }
});

db.kpi_filter.updateMany({}, {
    $rename: {
        entity_patterns: "old_entity_patterns",
    }
});
db.kpi_filter.find().forEach(function (doc) {
    if (doc.old_entity_patterns) {
        var newPattern = migrateOldEntityPatterns(doc.old_entity_patterns);
        if (newPattern) {
            db.kpi_filter.updateOne({_id: doc._id}, {
                $set: {entity_pattern: newPattern},
            });
        }
    }
});

for (var collectionName of ["idle_rule", "dynamic_infos", "instruction", "resolve_rule", "flapping_rule"]) {
    var collection = db.getCollection(collectionName);
    collection.updateMany({}, {
        $rename: {
            entity_patterns: "old_entity_patterns",
            alarm_patterns: "old_alarm_patterns",
        }
    });
    collection.find().forEach(function (doc) {
        var set = {};
        if (doc.old_entity_patterns) {
            var newEntityPattern = migrateOldEntityPatterns(doc.old_entity_patterns);
            if (newEntityPattern) {
                set.entity_pattern = newEntityPattern;
            }
        }
        if (doc.old_alarm_patterns) {
            var newAlarmPattern = migrateOldAlarmPatterns(doc.old_alarm_patterns);
            if (newAlarmPattern) {
                set.alarm_pattern = newAlarmPattern;
            }
        }

        if (Object.keys(set).length > 0) {
            collection.updateOne({_id: doc._id}, {
                $set: set,
            });
        }
    });
}

db.action_scenario.find().forEach(function (doc) {
    var newActions = [];

    doc.actions.forEach(function (action) {
        var newAction = action;

        if (newAction.entity_patterns) {
            newAction.old_entity_patterns = newAction.entity_patterns;
            var newEntityPattern = migrateOldEntityPatterns(newAction.entity_patterns);
            if (newEntityPattern) {
                newAction.entity_pattern = newEntityPattern;
            }
        }
        if (newAction.alarm_patterns) {
            newAction.old_alarm_patterns = newAction.alarm_patterns;
            var newAlarmPattern = migrateOldAlarmPatterns(newAction.alarm_patterns);
            if (newAlarmPattern) {
                newAction.alarm_pattern = newAlarmPattern;
            }
        }

        delete newAction.entity_patterns;
        delete newAction.alarm_patterns;

        newActions.push(newAction);
    });

    db.action_scenario.updateOne({_id: doc._id}, {$set: {actions: newActions}});
});

db.eventfilter.updateMany({}, {
    $rename: {
        patterns: "old_patterns",
    }
});
db.eventfilter.find().forEach(function (doc) {
    if (doc.old_patterns) {
        var newPatterns = migrateOldEventPatterns(doc.old_patterns);
        if (newPatterns) {
            var set = {};
            if (newPatterns[0]) {
                set["event_pattern"] = newPatterns[0];
            }
            if (newPatterns[1]) {
                set["entity_pattern"] = newPatterns[1];
            }

            if (Object.keys(set).length > 0) {
                db.eventfilter.updateOne({_id: doc._id}, {
                    $set: set,
                });
            }
        }
    }
});

db.meta_alarm_rules.updateMany({}, {
    $rename: {
        "config.alarm_patterns": "old_alarm_patterns",
        "config.entity_patterns": "old_entity_patterns",
        "config.total_entity_patterns": "old_total_entity_patterns",
        "config.event_patterns": "old_event_patterns",
    }
});
db.meta_alarm_rules.find().forEach(function (doc) {
    var set = {};
    if (doc.old_alarm_patterns) {
        var newAlarmPattern = migrateOldAlarmPatterns(doc.old_alarm_patterns);
        if (newAlarmPattern) {
            set.alarm_pattern = newAlarmPattern;
        }
    }
    if (doc.old_entity_patterns) {
        var newEntityPattern = migrateOldEntityPatterns(doc.old_entity_patterns);
        if (newEntityPattern) {
            set.entity_pattern = newEntityPattern;
        }
    }
    if (doc.old_total_entity_patterns) {
        var newEntityPattern = migrateOldEntityPatterns(doc.old_total_entity_patterns);
        if (newEntityPattern) {
            set.total_entity_pattern = newEntityPattern;
        }
    }
    if (doc.old_event_patterns && !set.entity_pattern) {
        var newPatterns = migrateOldEventPatterns(doc.old_event_patterns);
        if (newPatterns && newPatterns[0] === null && newPatterns[1]) {
            set.entity_pattern = newPatterns[1];
        }
    }

    if (Object.keys(set).length > 0) {
        db.meta_alarm_rules.updateOne({_id: doc._id}, {
            $set: set,
        });
    }
});

db.pbehavior.updateMany({}, {
    $rename: {
        filter: "old_mongo_query",
    }
});

function migrateOldGroup(oldGroup) {
    var newGroup = [];

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
                    newGroup.push({
                        field: field,
                        cond: strCond,
                    });
                    break;
                default:
                    if (field.startsWith("infos.") && field.endsWith(".value")) {
                        var info = field.replace(".value", "");
                        if (strCond !== null) {
                            newGroup.push({
                                field: info,
                                field_type: "string",
                                cond: strCond,
                            });
                        } else if (value === null) {
                            newGroup.push({
                                field: info,
                                cond: {
                                    type: "exist",
                                    value: false
                                },
                            });
                        } else {
                            return null;
                        }
                        break;
                    }

                    return null;
            }
        }
    }

    if (newGroup.length > 0) {
        return newGroup;
    }

    return null;
}

db.pbehavior.find().forEach(function (doc) {
    if (!doc.old_mongo_query || doc.old_mongo_query === "") {
        return;
    }

    var query = JSON.parse(doc.old_mongo_query);
    var newPattern = [];
    if (typeof query === "object" && Object.keys(query).length === 1) {
        var highLevelAnd = query["$and"];
        var highLevelOr = query["$or"];

        if (highLevelAnd && Array.isArray(highLevelAnd)) {
            var newGroup = migrateOldGroup(highLevelAnd);
            if (newGroup === null) {
                return;
            }

            newPattern.push(newGroup)
        } else if (highLevelOr && Array.isArray(highLevelOr)) {
            for (var oldGroup of highLevelOr) {
                if (typeof oldGroup !== "object" || !oldGroup) {
                    return;
                }

                var and = oldGroup["$and"];
                if (Object.keys(oldGroup).length === 1 && and && Array.isArray(and)) {
                    var newGroup = migrateOldGroup(and);
                    if (newGroup === null) {
                        return;
                    }

                    newPattern.push(newGroup)
                } else {
                    var newGroup = migrateOldGroup([oldGroup]);
                    if (newGroup === null) {
                        return;
                    }

                    newPattern.push(newGroup)
                }
            }
        } else {
            var newGroup = migrateOldGroup([query]);
            if (newGroup === null) {
                return;
            }

            newPattern.push(newGroup)
        }
    }

    if (newPattern.length > 0) {
        db.pbehavior.updateOne({_id: doc._id}, {$set: {entity_pattern: newPattern}});
    }
});

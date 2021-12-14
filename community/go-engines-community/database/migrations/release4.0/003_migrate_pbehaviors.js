(function () {
    function genID() {
        return UUID().toString().split('"')[1]
    }

    // Find default active type
    var activeType = db.pbehavior_type.findOne({type: "active"});
    var activeTypeID
    if (activeType) {
        activeTypeID = activeType._id;
    } else {
        return;
    }

    // Migrate pbehaviors
    db.default_pbehavior.find().forEach(function (doc) {
        var now = Math.ceil((new Date()).getTime() / 1000);
        // Find or create reason
        var reason = db.pbehavior_reason.findOne({name: doc.reason});
        var reasonID;
        if (reason) {
            reasonID = reason._id;
        } else {
            var reasonInsert = db.pbehavior_reason.insertOne({
                _id: genID(),
                name: doc.reason,
                description: doc.reason,
                created: now,
            });

            reasonID = reasonInsert.insertedId;
            if (!reasonID) {
                return;
            }
        }

        // Find type
        var canonicalType;

        switch (doc.type_.toLowerCase()) {
            case "maintenance":
                canonicalType = "maintenance";
                break;
            case "pause":
                canonicalType = "pause";
                break;
            case "hors plage horaire de surveillance":
                canonicalType = "inactive";
                break;
            default:
                return;
        }

        var type = db.pbehavior_type.findOne({type: canonicalType});
        var typeID;
        if (type) {
            typeID = type._id;
        } else {
            return;
        }

        // Update exdates
        var diff = doc.tstop - doc.tstart;
        var exdates = [];
        if (doc.exdate && doc.exdate.length > 0) {
            doc.exdate.forEach(function (start) {
                exdates.push({
                    begin: start,
                    end: start + diff,
                    type: activeTypeID
                });
            });
        }

        if (doc.comments === undefined) {
            doc.comments = [];
        }
        // Insert pbehavior
        var r = db.pbehavior.insertOne({
            _id: genID(),
            author: doc.author,
            enabled: doc.enabled,
            filter: doc.filter,
            name: doc.name,
            reason: reasonID,
            rrule: doc.rrule,
            tstart: doc.tstart,
            tstop: doc.tstop,
            type_: typeID,
            comments: doc.comments,
            exdates: exdates,
            created: now,
        });

        var pbehaviorID = r.insertedId
        if (!pbehaviorID) {
            return;
        }

        // Update alarms
        if (doc.eids && doc.eids.length > 0) {
            var aggregate = db.default_entities.aggregate([
                { $match: { _id: { $in: doc.eids } } },
                { $group: { _id: "", ids: { $push: "$_id" } } }
            ]);
            var res = (aggregate.hasNext()?aggregate.next():null);
            if (res && res.ids) {
                db.periodical_alarm.update(
                    {
                        $and: [
                            { d: { $in: res.ids } },
                            { "v.pbehavior_info": { "$exists": false } },
                            {
                                $or: [
                                    { "v.resolved": null },
                                    { "v.resolved": { "$exists": false } },
                                ]
                            },
                        ]
                    },
                    {
                        $set: {
                            "v.pbehavior_info": {
                                id: pbehaviorID,
                                name: doc.name,
                                reason: doc.reason,
                                type: type._id,
                                type_name: type.name,
                                canonical_type: type.type,
                            }
                        }
                    }
                );
            }
        }
    });
})();

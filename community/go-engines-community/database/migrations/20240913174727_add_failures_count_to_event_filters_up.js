db.eventfilter.find().forEach(function (doc) {
    var failuresCount = db.eventfilter_failure.countDocuments({rule: doc._id});
    var unreadFailuresCount = db.eventfilter_failure.countDocuments({rule: doc._id, unread: true});
    var set = {};
    if (failuresCount > 0) {
        set["failures_count"] = failuresCount;
    }

    if (unreadFailuresCount > 0) {
        set["unread_failures_count"] = unreadFailuresCount;
    }

    if (Object.keys(set).length > 0) {
        db.eventfilter.updateOne({_id: doc._id}, {$set: set});
    }
});

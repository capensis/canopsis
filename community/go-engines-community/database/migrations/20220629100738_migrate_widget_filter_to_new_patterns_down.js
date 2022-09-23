// down script can be applied only if no widget filter has been updated by API
db.widget_filters.aggregate([
    {$match: {is_private: false}},
    {
        $group: {
            _id: "$widget",
            filters: {$push: "$$ROOT"},
        }
    },
    {
        $lookup: {
            from: "widgets",
            localField: "_id",
            foreignField: "_id",
            as: "widget",
        }
    },
    {$unwind: "$widget"},
]).forEach(function (doc) {
    var viewFilters = [];
    var oldMainFilter = null;
    for (var filter of doc.filters) {
        viewFilters.push({
            title: filter.title,
            filter: filter.old_mongo_query,
        });

        if (doc.widget.parameters && doc.widget.parameters.mainFilter === filter._id) {
            oldMainFilter = {
                title: filter.title,
                filter: filter.old_mongo_query,
            };
        }
    }

    db.widgets.updateOne({_id: doc.widget._id}, {
        $set: {
            "parameters.mainFilter": oldMainFilter,
            "parameters.viewFilters": viewFilters,
        },
    });
});

db.widget_filters.aggregate([
    {$match: {is_private: true}},
    {
        $group: {
            _id: {
                widget: "$widget",
                user: "$author",
            },
            widget: {$first: "$widget"},
            user: {$first: "$author"},
            filters: {$push: "$$ROOT"},
        }
    },
    {
        $lookup: {
            from: "userpreferences",
            let: {widget: "$widget", user: "$user"},
            pipeline: [
                {
                    $match: {
                        $and: [
                            {$expr: {$eq: ["$user", "$$user"]}},
                            {$expr: {$eq: ["$widget", "$$widget"]}}
                        ]
                    }
                }
            ],
            as: "userpreferences",
        }
    },
    {$unwind: "$userpreferences"},
]).forEach(function (doc) {
    var viewFilters = [];
    var oldMainFilter = null;
    for (var filter of doc.filters) {
        viewFilters.push({
            title: filter.title,
            filter: filter.old_mongo_query,
        });

        if (doc.userpreferences.content && doc.userpreferences.content.mainFilter === filter._id) {
            oldMainFilter = {
                title: filter.title,
                filter: filter.old_mongo_query,
            };
        }
    }

    db.userpreferences.updateOne({_id: doc.userpreferences._id}, {
        $set: {
            "content.mainFilter": oldMainFilter,
            "content.viewFilters": viewFilters,
        },
    });
});

db.widget_filters.drop();

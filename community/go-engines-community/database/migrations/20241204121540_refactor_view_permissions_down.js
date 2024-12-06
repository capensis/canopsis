db.permission.aggregate([
    {$match: {view: {$ne: null}}},
    {
        $lookup: {
            from: "views",
            localField: "view",
            foreignField: "_id",
            as: "view",
        }
    },
    {$unwind: "$view"},
]).forEach(function (doc) {
    switch (doc.name) {
        case "view_general":
            db.permission.updateOne({_id: doc._id}, {
                $set: {
                    name: doc.view._id,
                    groups: [
                        "commonviews",
                        doc.view.group_id,
                    ],
                },
                $unset: {
                    view: ""
                }
            });

            return;
        case "view_actions":
            db.permission.deleteOne({_id: doc._id});
            db.role.updateMany({["permissions." + doc._id]: {$ne: null}}, {$unset: {["permissions." + doc._id]: ""}});

            return;
    }
});

db.permission.aggregate([
    {$match: {title: {$ne: null}}},
    {
        $lookup: {
            from: "view_playlist",
            localField: "_id",
            foreignField: "_id",
            as: "playlist",
        }
    },
    {$unwind: "$playlist"},
]).forEach(function (doc) {
    db.permission.updateOne({_id: doc._id}, {
        $set: {
            name: doc.playlist._id,
        },
        $unset: {
            title: ""
        }
    });
});

db.permission_group.updateMany({name: {$ne: null}}, [
    {
        $unset: {
            name: ""
        }
    }
]);

db.permission.aggregate([
    {$match: {view: null}},
    {
        $lookup: {
            from: "views",
            localField: "_id",
            foreignField: "_id",
            as: "view",
        }
    },
    {$unwind: "$view"},
]).forEach(function (doc) {
    db.permission.updateOne({_id: doc._id}, {
        $set: {
            name: "view_general",
            view: doc.view._id,
            groups: [
                "commonviews",
                doc.view.group_id,
                doc.view._id
            ],
        }
    });
    const newPermID = genID();
    db.permission.insertOne({
        _id: newPermID,
        name: "view_actions",
        view: doc.view._id,
        description: doc.view.title,
        groups: [
            "commonviews",
            doc.view.group_id,
            doc.view._id
        ],
    });
    db.role.updateOne({name: "admin"}, {$set: {["permissions." + newPermID]: 1}});
    db.role.updateMany({["permissions." + doc._id]: {$bitsAllSet: [3]}}, {$set: {["permissions." + newPermID]: 1}});
});

db.permission.aggregate([
    {$match: {title: null}},
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
            title: doc.playlist.name,
        },
        $unset: {
            name: ""
        }
    });
});

db.permission_group.updateMany({name: null}, [
    {
        $set: {
            name: "$_id"
        }
    }
]);

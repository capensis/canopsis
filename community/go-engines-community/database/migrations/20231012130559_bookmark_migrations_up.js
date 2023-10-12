if (!db.permission.findOne({_id: "listalarm_addBookmark"})) {
    db.permission.insertOne({
        _id: "listalarm_addBookmark",
        name: "listalarm_addBookmark",
        description: "List alarm add bookmark"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.listalarm_addBookmark": 1}});
}

if (!db.permission.findOne({_id: "listalarm_removeBookmark"})) {
    db.permission.insertOne({
        _id: "listalarm_removeBookmark",
        name: "listalarm_removeBookmark",
        description: "List alarm remove bookmark"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.listalarm_removeBookmark": 1}});
}

if (!db.permission.findOne({_id: "listalarm_filterByBookmark"})) {
    db.permission.insertOne({
        _id: "listalarm_filterByBookmark",
        name: "listalarm_filterByBookmark",
        description: "List alarm filter by bookmark"
    });
    db.role.updateMany({_id: "admin"}, {$set: {"permissions.listalarm_filterByBookmark": 1}});
}


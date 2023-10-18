db.permission.deleteMany({_id: {$in: ["listalarm_addBookmark", "listalarm_removeBookmark", "listalarm_filterByBookmark"]}});
db.role.updateMany(
    {},
    {
        $unset: {
            "permissions.listalarm_addBookmark": "",
            "permissions.listalarm_removeBookmark": "",
            "permissions.listalarm_filterByBookmark": "",
        }
    },
);

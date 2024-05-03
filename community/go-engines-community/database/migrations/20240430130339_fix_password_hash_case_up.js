db.user.updateMany(
    {password: {$not: {$regex: /^\$2a/}}},
    [
        {
            $set: {
                password: {$toLower: "$password"}
            }
        }
    ]
);

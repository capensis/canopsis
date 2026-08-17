db.eventfilter.updateMany(
    {name: {$exists: false}},
    [
        {
            $set: {
                name: {
                    $cond: [
                        {$eq: ["$description", ""]},
                        "$_id",
                        "$description"
                    ]
                }
            }
        }
    ]
);

db.pbehavior.updateMany(
    {
        alarm_count: {
            $exists: false
        }
    },
    {
        $set: {
            alarm_count: 0
        }
    }
)

db.pbehavior.updateMany(
    {},
    {
        $rename: {
            type_: "type",
        }
    }
)

db.webhooks.update(
    {
        disable_if_active_pbehavior: true
    },
    {
        $set: {
            disable_on_periods: ["inactive", "maintenance", "pause"]
        }
    }
);
db.webhooks.update(
    {},
    {
        $unset: {
            disable_if_active_pbehavior: ""
        }
    }
);
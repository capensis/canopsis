db.meta_alarm_rules.createIndex(
    {soft_deleted: 1},
    {name: "soft_deleted_1", partialFilterExpression: {soft_deleted: {$exists: true}}}
);

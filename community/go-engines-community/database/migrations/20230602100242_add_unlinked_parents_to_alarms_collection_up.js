db.periodical_alarm.updateMany(
    {"v.unlinked_parents": {$exists: false}},
    {$set: {"v.unlinked_parents": []}},
);

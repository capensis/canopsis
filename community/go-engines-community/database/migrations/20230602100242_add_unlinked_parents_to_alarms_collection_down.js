db.periodical_alarm.updateMany(
    {"v.unlinked_parents": {$exists: true}},
    {$unset: {"v.unlinked_parents": ""}},
);

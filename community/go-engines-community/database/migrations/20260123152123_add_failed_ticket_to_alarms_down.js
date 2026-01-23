db.periodical_alarm.updateMany(
    {"v.failed_ticket": {$ne: null}},
    {$unset: {"v.failed_ticket": null}},
);

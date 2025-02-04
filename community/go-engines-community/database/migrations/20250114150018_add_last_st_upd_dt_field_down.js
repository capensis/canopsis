db.periodical_alarm.updateMany({}, {
    $unset: {
        "v.last_st_upd_dt": ""
    }
})

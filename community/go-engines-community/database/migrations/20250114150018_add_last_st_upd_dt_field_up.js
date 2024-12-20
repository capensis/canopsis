db.periodical_alarm.updateMany({"v.last_st_upd_dt": {$exists: false}}, [
    {
        $set: {
            "v.last_st_upd_dt": "$v.last_update_date"
        }
    }
])

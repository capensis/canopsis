db.periodical_alarm.updateMany({"v.state._t": "changestate"}, [{$set: {"v.change_state": "$v.state"}}])

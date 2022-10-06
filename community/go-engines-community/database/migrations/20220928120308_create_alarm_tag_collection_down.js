db.alarm_tag.dropIndex("value_1");

db.configuration.deleteOne({_id: "alarm_tag_color"});

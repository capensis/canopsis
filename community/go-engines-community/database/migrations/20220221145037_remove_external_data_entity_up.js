db.eventfilter.deleteMany({"external_data.entity": {$exists: true}});

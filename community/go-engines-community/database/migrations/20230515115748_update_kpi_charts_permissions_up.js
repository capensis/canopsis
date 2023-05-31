db.default_rights.updateOne({_id: "barchart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_interval": {$exists: true}}, {$set: {"rights.barchart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_sampling": {$exists: true}}, {$set: {"rights.barchart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_listFilters": {$exists: true}}, {$set: {"rights.barchart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_editFilter": {$exists: true}}, {$set: {"rights.barchart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_addFilter": {$exists: true}}, {$set: {"rights.barchart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_userFilter": {$exists: true}}, {$set: {"rights.barchart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_interval": {$exists: true}}, {$set: {"rights.linechart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_sampling": {$exists: true}}, {$set: {"rights.linechart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_listFilters": {$exists: true}}, {$set: {"rights.linechart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_editFilter": {$exists: true}}, {$set: {"rights.linechart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_addFilter": {$exists: true}}, {$set: {"rights.linechart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_userFilter": {$exists: true}}, {$set: {"rights.linechart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_interval": {$exists: true}}, {$set: {"rights.piechart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_sampling": {$exists: true}}, {$set: {"rights.piechart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_listFilters": {$exists: true}}, {$set: {"rights.piechart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_editFilter": {$exists: true}}, {$set: {"rights.piechart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_addFilter": {$exists: true}}, {$set: {"rights.piechart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_userFilter": {$exists: true}}, {$set: {"rights.piechart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_interval"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_interval": {$exists: true}}, {$set: {"rights.numbers_interval.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_sampling": {$exists: true}}, {$set: {"rights.numbers_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_listFilters": {$exists: true}}, {$set: {"rights.numbers_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_editFilter": {$exists: true}}, {$set: {"rights.numbers_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_addFilter": {$exists: true}}, {$set: {"rights.numbers_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_userFilter": {$exists: true}}, {$set: {"rights.numbers_userFilter.checksum": 1}});

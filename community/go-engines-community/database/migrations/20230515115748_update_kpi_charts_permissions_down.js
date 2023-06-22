db.default_rights.updateOne({_id: "barchart_interval"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_interval": {$exists: true}}, {$set: {"rights.barchart_interval.checksum": 15}});

db.default_rights.updateOne({_id: "barchart_sampling"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_sampling": {$exists: true}}, {$set: {"rights.barchart_sampling.checksum": 15}});

db.default_rights.updateOne({_id: "barchart_listFilters"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_listFilters": {$exists: true}}, {$set: {"rights.barchart_listFilters.checksum": 15}});

db.default_rights.updateOne({_id: "barchart_editFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_editFilter": {$exists: true}}, {$set: {"rights.barchart_editFilter.checksum": 15}});

db.default_rights.updateOne({_id: "barchart_addFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_addFilter": {$exists: true}}, {$set: {"rights.barchart_addFilter.checksum": 15}});

db.default_rights.updateOne({_id: "barchart_userFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.barchart_userFilter": {$exists: true}}, {$set: {"rights.barchart_userFilter.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_interval"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_interval": {$exists: true}}, {$set: {"rights.linechart_interval.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_sampling"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_sampling": {$exists: true}}, {$set: {"rights.linechart_sampling.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_listFilters"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_listFilters": {$exists: true}}, {$set: {"rights.linechart_listFilters.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_editFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_editFilter": {$exists: true}}, {$set: {"rights.linechart_editFilter.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_addFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_addFilter": {$exists: true}}, {$set: {"rights.linechart_addFilter.checksum": 15}});

db.default_rights.updateOne({_id: "linechart_userFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.linechart_userFilter": {$exists: true}}, {$set: {"rights.linechart_userFilter.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_interval"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_interval": {$exists: true}}, {$set: {"rights.piechart_interval.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_sampling"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_sampling": {$exists: true}}, {$set: {"rights.piechart_sampling.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_listFilters"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_listFilters": {$exists: true}}, {$set: {"rights.piechart_listFilters.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_editFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_editFilter": {$exists: true}}, {$set: {"rights.piechart_editFilter.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_addFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_addFilter": {$exists: true}}, {$set: {"rights.piechart_addFilter.checksum": 15}});

db.default_rights.updateOne({_id: "piechart_userFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.piechart_userFilter": {$exists: true}}, {$set: {"rights.piechart_userFilter.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_interval"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_interval": {$exists: true}}, {$set: {"rights.numbers_interval.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_sampling"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_sampling": {$exists: true}}, {$set: {"rights.numbers_sampling.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_listFilters"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_listFilters": {$exists: true}}, {$set: {"rights.numbers_listFilters.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_editFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_editFilter": {$exists: true}}, {$set: {"rights.numbers_editFilter.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_addFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_addFilter": {$exists: true}}, {$set: {"rights.numbers_addFilter.checksum": 15}});

db.default_rights.updateOne({_id: "numbers_userFilter"}, {$set: {"type": "RW"}});
db.default_rights.updateMany({crecord_type: "role", "rights.numbers_userFilter": {$exists: true}}, {$set: {"rights.numbers_userFilter.checksum": 15}});

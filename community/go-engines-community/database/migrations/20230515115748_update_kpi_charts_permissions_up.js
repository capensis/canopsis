db.default_rights.updateOne({_id: "barchart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "barchart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.barchart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "linechart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.linechart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_interval"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_interval.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "piechart_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.piechart_userFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_interval"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_interval.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_sampling"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_sampling.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_listFilters"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_listFilters.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_editFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_editFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_addFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_addFilter.checksum": 1}});

db.default_rights.updateOne({_id: "numbers_userFilter"}, {$unset: {"type": 1}});
db.default_rights.updateOne({_id: "admin"}, {$set: {"rights.numbers_userFilter.checksum": 1}});

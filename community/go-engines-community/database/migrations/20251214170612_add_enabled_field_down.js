db.meta_alarm_rules.updateMany({}, {$unset: {enabled: true}});
db.flapping_rule.updateMany({}, {$unset: {enabled: true}});
db.default_snmprules.updateMany({}, {$unset: {enabled: true}});
db.resolve_rule.updateMany({}, {$unset: {enabled: true}});

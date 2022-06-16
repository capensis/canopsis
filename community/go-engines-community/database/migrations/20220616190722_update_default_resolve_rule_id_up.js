var rule = db.resolve_rule.findOne({name: "Default rule"});

db.resolve_rule.deleteOne({_id: rule._id})

rule._id = "default-rule"

db.resolve_rule.insert(rule)

db.meta_alarm_rules.updateMany({
    output_template: ""
}, {$set: {output_template: "Règle : {{ .Rule.Name }}"}})
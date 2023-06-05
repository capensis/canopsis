db.alarm_tag.createIndex({value: 1}, {name: "value_1", unique: true});

if (!db.configuration.findOne({_id: "alarm_tag_color"})) {
    db.configuration.insertOne({
        _id: "alarm_tag_color",
        colors: [
            "#B71C1C",
            "#880E4F",
            "#4A148C",
            "#0D47A1",
            "#01579B",
            "#006064",
            "#004D40",
            "#33691E",
            "#9E9D24",
            "#F57F17",
            "#E65100",
            "#BF360C",
            "#8D6E63",
            "#607D8B",
            "#9E9E9E",
        ],
    });
}

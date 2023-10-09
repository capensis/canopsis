db.alarm_tag.aggregate([
    { $sort: { created: -1 } },
    {
        $group: {
            _id: "$label",
            color: {
                $first: "$color"
            }
        }
    },
    { $out: "alarm_tag_color" }
]);

db.alarm_tag_color.find({}).forEach(function (doc) {
    db.alarm_tag.updateMany({"label": doc._id}, {"$set": {"color": doc.color}})
});

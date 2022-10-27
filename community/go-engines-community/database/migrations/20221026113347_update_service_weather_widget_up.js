db.widgets.find({type: "ServiceWeather"}).forEach(function (doc) {
    var set = {};
    var unset = {};

    set["parameters.counters.pbehavior_enabled"] = false;
    set["parameters.counters.pbehavior_types"] = [];
    set["parameters.counters.state_enabled"] = false;
    set["parameters.counters.state_types"] = [];
    unset["parameters.counters.enabled"] = "";
    unset["parameters.counters.types"] = "";

    if (doc.parameters && doc.parameters.counters) {
        if (doc.parameters.counters.enabled) {
            set["parameters.counters.pbehavior_enabled"] = doc.parameters.counters.enabled;
        }
        if (doc.parameters.counters.types) {
            set["parameters.counters.pbehavior_types"] = doc.parameters.counters.types;
        }
    }

    db.widgets.updateOne({_id: doc._id}, {$set: set, $unset: unset});
});

(function () {
    db.configuration.find({"_id": "user_interface"}).forEach(function (doc) {
        var popupTimeout = doc.popup_timeout;

        if (!popupTimeout) {
            return;
        }

        Object.keys(popupTimeout).forEach(function (key) {
            var val = popupTimeout[key].interval;
            var unit = popupTimeout[key].unit;
            var seconds = val;
            switch (unit) {
                case "s":
                    seconds = val;
                    break
                case "m":
                    seconds = val * 60;
                    break
                case "h":
                    seconds = val * 60 * 60;
                    break
            }

            delete popupTimeout[key].interval;
            popupTimeout[key].seconds = seconds;
        });

        db.configuration.updateOne({_id: doc._id}, {
            $set: {
                "popup_timeout": popupTimeout,
            },
        });
    });
})();

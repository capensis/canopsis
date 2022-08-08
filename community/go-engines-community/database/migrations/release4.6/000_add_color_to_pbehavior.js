(function () {
    db.pbehavior.find({"color": {$in: ["", null]}}).forEach(function (doc) {
        db.pbehavior.updateOne({_id: doc._id}, {$set: {"color": getRandomColor()}});
    });

    function getRandomColor() {
        var letters = '0123456789ABCDEF';
        var color = '#';
        for (var i = 0; i < 6; i++) {
            color += letters[Math.floor(Math.random() * 16)];
        }
        return color;
    }
})();

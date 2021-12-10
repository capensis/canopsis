(function () {
    db.views.find().forEach(function (doc) {
        if (!doc.tabs) {
            return;
        }

        var tabs = [];
        var widgets = [];

        doc.tabs.forEach(function (tab) {
            if (tab.widgets) {
                tab.widgets.forEach(function (widget) {
                    widget.tab = tab._id;
                    widget.author = doc.author;
                    widget.created = doc.created;
                    widget.updated = doc.updated;
                    widgets.push(widget);
                });
            }

            delete tab.widgets;
            tab.view = doc._id;
            tab.author = doc.author;
            tab.created = doc.created;
            tab.updated = doc.updated;
            tabs.push(tab);
        });

        db.views.updateOne({_id: doc._id}, {$unset: {tabs: ""}});
        db.viewtabs.insertMany(tabs);
        if (widgets.length > 0) {
            db.widgets.insertMany(widgets);
        }
    });
})();

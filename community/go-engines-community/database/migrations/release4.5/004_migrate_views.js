(function () {
    db.views.find().forEach(function (doc) {
        if (!doc.tabs) {
            return;
        }

        var tabs = [];
        var widgets = [];
        var filters = [];

        var now = Math.ceil((new Date()).getTime() / 1000);
        var created = doc.created;
        if (!created) {
            created = now;
        }
        var updated = doc.updated;
        if (!updated) {
            updated = now;
        }
        var author = doc.author;
        if (!author) {
            author = "";
        }

        doc.tabs.forEach(function (tab) {
            if (tab.widgets) {
                tab.widgets.forEach(function (widget) {
                    if (widget.parameters.viewFilters) {
                        var mainFilterId = null;

                        var mainFilter = widget.parameters.mainFilter;
                        widget.parameters.viewFilters.forEach(function (filter, fi) {
                            var filterId = widget._id + "filter_" + (fi+1);
                            filters.push({
                                _id: filterId,
                                loader_id: filterId,
                                widget: widget._id,
                                title: filter.title,
                                query: filter.filter,
                                author: author,
                                created: created,
                                updated: updated,
                            });

                            if (mainFilter && mainFilter.title === filter.title) {
                                mainFilterId = filterId;
                            }
                        });

                        delete widget.parameters.viewFilters;
                        delete widget.parameters.mainFilter;

                        if (mainFilterId) {
                            widget.parameters.main_filter = mainFilterId;
                        }
                    }

                    widget.loader_id = widget._id;
                    widget.tab = tab._id;
                    widget.author = author;
                    widget.created = created;
                    widget.updated = updated;
                    widgets.push(widget);
                });
            }

            delete tab.widgets;
            tab.loader_id = tab._id;
            tab.view = doc._id;
            tab.author = author;
            tab.created = created;
            tab.updated = updated;
            tabs.push(tab);
        });

        db.views.updateOne({_id: doc._id}, {$unset: {tabs: ""}});
        db.viewtabs.insertMany(tabs);
        if (widgets.length > 0) {
            db.widgets.insertMany(widgets);
        }
        if (filters.length > 0) {
            db.widget_filters.insertMany(filters);
        }
    });

    db.userpreferences.find().forEach(function (doc) {
        var filters = [];

        var now = Math.ceil((new Date()).getTime() / 1000);
        var updated = doc.updated;
        if (!updated) {
            updated = now;
        }

        if (doc.content && doc.content.viewFilters) {
            var mainFilterId = null;

            var mainFilter = doc.content.mainFilter;
            doc.content.viewFilters.forEach(function (filter, fi) {
                var filterId = doc.widget + "_" + doc.user + "filter_" + (fi+1);
                filters.push({
                    _id: filterId,
                    loader_id: filterId,
                    widget: doc.widget,
                    user: doc.user,
                    title: filter.title,
                    query: filter.filter,
                    author: doc.user,
                    created: updated,
                    updated: updated,
                });

                if (mainFilter && mainFilter.title === filter.title) {
                    mainFilterId = filterId;
                }
            });

            if (mainFilter && !mainFilterId) {
                var filter = db.widget_filters.findOne({
                    widget: doc.widget,
                    user: null,
                    title: mainFilter.title,
                });
                if (filter) {
                    mainFilterId = filter._id;
                }
            }

            delete doc.content.viewFilters;
            delete doc.content.mainFilter;

            if (mainFilterId) {
                doc.content.main_filter = mainFilterId;
            }
        }

        if (filters.length > 0) {
            db.widget_filters.insertMany(filters);
        }
    });
})();

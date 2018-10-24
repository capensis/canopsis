# Upgrading to Canopsis 3.2.5

## Add rights on menu items

To add new rights for menu items, just rerun canopsinit (be carefull, it can destroy userviews !)
```bash
su - canopsis
/opt/canopsis/bin/canopsinit --authorize-reinit
```

Or just insert the corresponding values in `default_rights`

```bash
mongodb:PRIMARY> db.default_rights.insert([{
    "_id" : "showview_view_filters",
    "crecord_type" : "action",
    "loader_id" : "showview_view_filters",
    "crecord_name" : "showview_view_filters",
    "desc" : "Access to view.filters menu"
},{
    "_id" : "showview_view_context",
    "crecord_type" : "action",
    "loader_id" : "showview_view_context",
    "crecord_name" : "showview_view_context",
    "desc" : "Access to view.context menu"
},{
    "_id" : "showview_view_selectors",
    "crecord_type" : "action",
    "loader_id" : "showview_view_selectors",
    "crecord_name" : "showview_view_selectors",
    "desc" : "Access to view.selectors menu"
},{
    "_id" : "showview_view_series",
    "crecord_type" : "action",
    "loader_id" : "showview_view_series",
    "crecord_name" : "showview_view_series",
    "desc" : "Access to view.series menu"
},{
    "_id" : "showview_view_jobs",
    "crecord_type" : "action",
    "loader_id" : "showview_view_jobs",
    "crecord_name" : "showview_view_jobs",
    "desc" : "Access to view.jobs menu"
},{
    "_id" : "showview_view_notifications",
    "crecord_type" : "action",
    "loader_id" : "showview_view_notifications",
    "crecord_name" : "showview_view_notifications",
    "desc" : "Access to view.notifications menu"
}])
```

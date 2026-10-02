const defaultTemplateID = "62500000-0000-4000-8000-000000000001";
const defaultTemplateTitle = "Default ticket formats";
const defaultTemplateType = "ticket_formats";

if (db.widget_templates.countDocuments({_id: defaultTemplateID}) === 0 &&
    db.widget_templates.countDocuments({title: defaultTemplateTitle, type: defaultTemplateType}) === 0) {
    const now = Math.ceil((new Date()).getTime() / 1000);
    db.widget_templates.insertOne({
        _id: defaultTemplateID,
        title: defaultTemplateTitle,
        type: defaultTemplateType,
        formats: [],
        author: "root",
        created: now,
        updated: now
    });
}

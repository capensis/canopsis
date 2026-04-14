db.event_records.find({ c: { $gte: 0 }, e: { $exists: false } }).forEach(function(doc) {
  var recordingEndedAt = doc.t + 1;
  db.event_records.updateOne({ _id: doc._id }, { $set: { e: recordingEndedAt } });
});

db.event_records.createIndex({ t: -1 }, { name: "t_-1", partialFilterExpression: { t: { $exists: true } } });
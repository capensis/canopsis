if (!db.getCollectionNames().includes('message_rate_statistic_minute')) {
    db.createCollection("message_rate_statistic_minute", {capped: true, size: 100000, max: 1440});
}

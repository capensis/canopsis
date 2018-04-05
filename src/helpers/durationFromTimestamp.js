Ember.Application.initializer({
    name: 'durationFromTimestampHelpers',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);
        var datesUtils = container.lookupFactory('utility:dates');
        var __ = Ember.String.loc;
        Ember.Handlebars.helper('durationFromTimestamp', function(timestamp , record) {
			if(timestamp || record.timeStampState) {
                timestamp = record.timeStampState || timestamp;
				timestamp = timestamp - (timestamp % 60)
                var time = datesUtils.second2Duration(timestamp);
				time = time.replace(" 00s", "")
                return time;
            }
        });
    }
});

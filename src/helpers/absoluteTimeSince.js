Ember.Application.initializer({
    name: 'AbsoluteTimesinceHelper',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);
        var datesUtils = container.lookupFactory('utility:dates');
        var __ = Ember.String.loc;
        Ember.Handlebars.helper('absoluteTimeSince', function(timestamp , record) {
            if(timestamp || record.timeStampState) {
                timestamp = record.timeStampState || timestamp;
                var time = datesUtils.durationFromNow(timestamp);
                return time;
            } else {
                return '';
            }
        });
    }
});
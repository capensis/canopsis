App.IndexRoute = Ember.Route.extend({
    model: function() {
        return {
            'array': [
                1,
                2,
                3
            ],
            'boolean': true,
            'null': null,
            'number': 123,
            'object': {
                'a': 'b',
                'c': 'd',
                'e': 'f'
            },
            'string': 'Hello World'
        };
    }
});

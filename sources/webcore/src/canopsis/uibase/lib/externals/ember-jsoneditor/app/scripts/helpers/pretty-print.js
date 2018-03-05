Ember.Handlebars.helper('pretty-print', function(value, replacer, space) {
    console.log(arguments);
    return JSON.stringify(value, replacer, space);
});

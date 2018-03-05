define(["ember"], function(Ember) {
  Ember.Handlebars.registerHelper("date", function(type) {
    return (new Date()).getFullYear();
  });
});

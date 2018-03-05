QUnit.config.testTimeout = 2000;

require.config({
  "paths": {
    "jquery": "jquery-1.9.1",
    "ember": "ember-1.0.0-rc.6",
    "handlebars": "handlebars-1.0.0-rc.4",
    "ehbs": "../ehbs"
  },
  "shim": {
    "ember": {
      "deps": ["jquery", "handlebars"],
      "exports": "Ember"
    },
    "handlebars": {
      "exports": "Handlebars"
    }
  },
  "ehbs": {
    "casing": "camel"
  }
});

asyncTest("Application template", function() {
  require(["ember", "ehbs!application"], function() {
    ok(!!Ember.TEMPLATES.application, "Application template compiled.");
    start();
  });
});

asyncTest("Template with single date helper", function() {
  require(["ember", "ehbs!main"], function() {
    ok(!!Ember.Handlebars.helpers.date, "Date helper included.");
    ok(!!Ember.TEMPLATES.main, "Main template compiled.");
    start();
  });
});

asyncTest("Template with single view", function() {
  require(["ember", "ehbs!mainWithView"], function() {
    ok(!!App.testView, "testView included.");
    ok(!!Ember.TEMPLATES.mainWithView, "Main template compiled.");
    start();
  });
});

asyncTest("Template with single partial", function() {
  require(["ember", "ehbs!mainWithPartial"], function() {
    ok(!!Ember.TEMPLATES.testPartial, "Partial template compiled.");
    ok(!!Ember.TEMPLATES.mainWithPartial, "Main template compiled.");
    start();
  });
});

asyncTest("Template with single controller", function() {
  require(["ember", "ehbs!mainWithController"], function() {
    ok(!!App.testController, "Controller included.");
    ok(!!Ember.TEMPLATES.mainWithController, "Main template compiled.");
    start();
  });
});

asyncTest("Template with single render", function() {
  require(["ember", "ehbs!mainWithRender"], function() {
    ok(!!Ember.TEMPLATES.testRender, "Test render compiled.");
    ok(!!Ember.TEMPLATES.mainWithRender, "Main template compiled.");
    start();
  });
});


var canopsisEnv = {
  url: 'http://demo-devel.canopsis.net/',
  username: 'root',
  password: 'root',
  authkey: '523fd618-87dc-11e7-943e-6aaa6f7effb9'
};

var eventAlarm =  { 
  connector: "nightwatch", 
  connector_name: "nishtwatchdocker" ,
  event_type: "check",
  source_type: "resource",
  component: "nightwatchcomponent",
  resource: "nightwatchresource",
  output: "nightwatch fake output",
  state: 3
};

module.exports = {
  canopsisEnv: canopsisEnv, 
  eventAlarm: eventAlarm,
  beforeEach: function (browser, done) {
    require('nightwatch-video-recorder').start(browser, done)
  },
  afterEach: function (browser, done) {
    require('nightwatch-video-recorder').stop(browser, done)
  }
}


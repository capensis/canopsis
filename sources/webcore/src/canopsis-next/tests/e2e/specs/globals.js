'use strict'

var canopsisEnv = {
  url: 'http://demo-devel.canopsis.net',
  username: 'root',
  password: 'root',
};

module.exports = {
  canopsisEnv: canopsisEnv, 
  beforeEach: function (browser, done) {
    require('nightwatch-video-recorder').start(browser, done)
  },
  afterEach: function (browser, done) {
    require('nightwatch-video-recorder').stop(browser, done)
  }
}


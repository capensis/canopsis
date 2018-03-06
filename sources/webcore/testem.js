var StandardReporter = require('./testem-reporters/standard.js');

module.exports = {
  "routes": {
    "/en/static": "src",
    "/static": "src"
  },
  "proxies": {
    "/rest": {
      "target": "http://localhost:8082"
    },
    "/context": {
      "target": "http://localhost:8082"
    },
     "/entitylink": {
      "target": "http://localhost:8082"
    },
    "/keepalive": {
      "target": "http://localhost:8082"
    },
    "/autologin": {
      "target": "http://localhost:8082"
    },
    "/perfdata": {
      "target": "http://localhost:8082"
    },
    "/sessionstart": {
      "target": "http://localhost:8082"
    },
    "/account": {
      "target": "http://localhost:8082"
    }
  },
  "reporter": new StandardReporter(),
  "report_file": "testem-output.txt",
  "test_page": "en/static/canopsis/index.test.html"
};

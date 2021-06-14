Feature: get test suite's summary
  I need to be able to get test suite's summary
  Only admin should be able to get test suite's summary

  Scenario: GET unauthorized
    When I do GET /api/v4/cat/junit/test-suites/5b4461d4-cfea-41dd-97cd-a5b15c34c1e2/summary
    Then the response code should be 401

  Scenario: GET without permissions
    When I am noperms
    When I do GET /api/v4/cat/junit/test-suites/5b4461d4-cfea-41dd-97cd-a5b15c34c1e2/summary
    Then the response code should be 403

  Scenario: GET summary success
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/5b4461d4-cfea-41dd-97cd-a5b15c34c1e2/summary
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "5b4461d4-cfea-41dd-97cd-a5b15c34c1e2",
      "xml_feed": "noveo.app.PublishMissionTest-filename",
      "name": "noveo.app.PublishMissionTest",
      "hostname": "noveo.app.PublishMissionTest-hostname",
      "last_update": 1614782430,
      "time": 0.19,
      "total": 1,
      "disabled": 0,
      "errors": 0,
      "failures": 0,
      "skipped": 1,
      "system_out": "system out message",
      "system_err": "system err message"
    }
    """

  Scenario: GET summary, test-suite not found
    When I am admin
    When I do GET /api/v4/cat/junit/test-suites/5b4461d4-cfea-41dd-97cd-a5b15c34c1e2-not-found/summary
    Then the response code should be 404
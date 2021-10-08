Feature: Auth user
  I need to be able to auth user by different methods

  Scenario: given guest user should not allow access
    When I do GET /api/v4/alarms
    Then the response code should be 401

  Scenario: given auth user by api key should allow access
    When I am admin
    When I do GET /api/v4/alarms
    Then the response code should be 200

  Scenario: given auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/alarms
    Then the response code should be 403

  Scenario: given auth user by basic auth should allow access
    When I am authenticated with username "root" and password "test"
    When I do GET /api/v4/alarms
    Then the response code should be 200

  Scenario: given auth user by basic auth without permissions should not allow access
    When I am authenticated with username "nopermsuser" and password "test"
    When I do GET /api/v4/alarms
    Then the response code should be 403

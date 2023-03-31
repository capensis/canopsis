Feature: Manual meta alarms

  @concurrent
  Scenario: given get unauth request should not allow access
    When I do GET /api/v4/cat/manual-meta-alarms
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/cat/manual-meta-alarms
    Then the response code should be 403

  @concurrent
  Scenario: given create unauth request should not allow access
    When I do POST /api/v4/cat/manual-meta-alarms
    Then the response code should be 401

  @concurrent
  Scenario: given create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/cat/manual-meta-alarms
    Then the response code should be 403

  @concurrent
  Scenario: given invalid create request should return errors
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name is missing.",
        "alarms": "Alarms is missing."
      }
    }
    """
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "alarms": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarms": "Alarms should not be blank."
      }
    }
    """

  @concurrent
  Scenario: given create request with not exist alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-not-exist",
      "alarms": ["test-alarm-not-exist"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given create request with resolved alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-not-exist",
      "alarms": ["test-alarm-manual-correlation-fail-2"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given create request with existed name should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-exist",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-fail-check-exist until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-fail-check-exist"
      }
    ]
    """
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-exist",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

  @concurrent
  Scenario: given add unauth request should not allow access
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/add
    Then the response code should be 401

  @concurrent
  Scenario: given add request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/add
    Then the response code should be 403

  @concurrent
  Scenario: given invalid add request should return errors
    When I am admin
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/add:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarms": "Alarms is missing."
      }
    }
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/add:
    """json
    {
      "alarms": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarms": "Alarms should not be blank."
      }
    }
    """

  @concurrent
  Scenario: given add request with not exist meta alarm should return error
    When I am admin
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/add:
    """json
    {
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given add request with not exist alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-add-1",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-fail-check-add-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-fail-check-add-1"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/add:
    """json
    {
      "alarms": ["test-alarm-not-exist"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given add request with resolved alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-add-2",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-fail-check-add-2 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-fail-check-add-2"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/add:
    """json
    {
      "alarms": ["test-alarm-manual-correlation-fail-2"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given remove unauth request should not allow access
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/remove
    Then the response code should be 401

  @concurrent
  Scenario: given remove request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/remove
    Then the response code should be 403

  @concurrent
  Scenario: given invalid remove request should return errors
    When I am admin
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/remove:
    """json
    {}
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarms": "Alarms is missing."
      }
    }
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/remove:
    """json
    {
      "alarms": []
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "alarms": "Alarms should not be blank."
      }
    }
    """

  @concurrent
  Scenario: given remove request with not exist meta alarm should return error
    When I am admin
    When I do PUT /api/v4/cat/manual-meta-alarms/test-metaalarm-not-exist/remove:
    """json
    {
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given remove request with not exist alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-remove-1",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-fail-check-remove-1 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-fail-check-remove-1"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/remove:
    """json
    {
      "alarms": ["test-alarm-not-exist"]
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given remove request with resolved alarm should return error
    When I am admin
    When I do POST /api/v4/cat/manual-meta-alarms:
    """json
    {
      "name": "test-metaalarm-manual-correlation-fail-check-remove-2",
      "alarms": ["test-alarm-manual-correlation-fail-1"]
    }
    """
    Then the response code should be 204
    When I do GET /api/v4/cat/manual-meta-alarms?search=test-metaalarm-manual-correlation-fail-check-remove-2 until response code is 200 and body contains:
    """json
    [
      {
        "name": "test-metaalarm-manual-correlation-fail-check-remove-2"
      }
    ]
    """
    When I do PUT /api/v4/cat/manual-meta-alarms/{{ (index .lastResponse 0)._id }}/remove:
    """json
    {
      "alarms": ["test-alarm-manual-correlation-fail-2"]
    }
    """
    Then the response code should be 404

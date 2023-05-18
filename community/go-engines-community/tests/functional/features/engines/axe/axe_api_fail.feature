Feature: update alarm
  I need to be able to update alarm

  @concurrent
  Scenario: given ack request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ack
    Then the response code should be 401

  @concurrent
  Scenario: given ack request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ack
    Then the response code should be 403

  @concurrent
  Scenario: given ack invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ack
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing."
      }
    }
    """

  @concurrent
  Scenario: given ack request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ack:
    """json
    {
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given ack remove request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ackremove
    Then the response code should be 401

  @concurrent
  Scenario: given ack remove request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ackremove
    Then the response code should be 403

  @concurrent
  Scenario: given ack remove invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ackremove
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing."
      }
    }
    """

  @concurrent
  Scenario: given ack remove request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/ackremove:
    """json
    {
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given snooze request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/snooze
    Then the response code should be 401

  @concurrent
  Scenario: given snooze request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/snooze
    Then the response code should be 403

  @concurrent
  Scenario: given snooze invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/snooze
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing.",
        "duration.unit": "Unit is missing.",
        "duration.value": "Value is missing."
      }
    }
    """
    When I do PUT /api/v4/alarms/test-alarm-not-exist/snooze:
    """json
    {
      "duration": {
        "value": 1,
        "unit": "y"
      },
      "comment": "test-comment"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "duration": "Duration is invalid."
      }
    }
    """

  @concurrent
  Scenario: given snooze request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/snooze:
    """json
    {
      "comment": "test-comment",
      "duration": {
        "value": 3,
        "unit": "m"
      }
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given cancel request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/cancel
    Then the response code should be 401

  @concurrent
  Scenario: given cancel request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/cancel
    Then the response code should be 403

  @concurrent
  Scenario: given cancel invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/cancel
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing."
      }
    }
    """

  @concurrent
  Scenario: given cancel request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/cancel:
    """json
    {
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given uncancel request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/uncancel
    Then the response code should be 401

  @concurrent
  Scenario: given uncancel request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/uncancel
    Then the response code should be 403

  @concurrent
  Scenario: given uncancel invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/uncancel
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing."
      }
    }
    """

  @concurrent
  Scenario: given uncancel request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/uncancel:
    """json
    {
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given assoc ticket request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/assocticket
    Then the response code should be 401

  @concurrent
  Scenario: given assoc ticket request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/assocticket
    Then the response code should be 403

  @concurrent
  Scenario: given assoc ticket invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/assocticket
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "ticket": "Ticket is missing."
      }
    }    
    """

  @concurrent
  Scenario: given assoc ticket request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/assocticket:
    """json
    {
      "ticket": "test-ticket"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given comment request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/comment
    Then the response code should be 401

  @concurrent
  Scenario: given comment request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/comment
    Then the response code should be 403

  @concurrent
  Scenario: given comment invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/comment
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "comment": "Comment is missing."
      }
    }
    """

  @concurrent
  Scenario: given comment request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/comment:
    """json
    {
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

  @concurrent
  Scenario: given change state request and no auth user should not allow access
    When I do PUT /api/v4/alarms/test-alarm-not-exist/changestate
    Then the response code should be 401

  @concurrent
  Scenario: given change state request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/alarms/test-alarm-not-exist/changestate
    Then the response code should be 403

  @concurrent
  Scenario: given change state invalid request should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/changestate
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "state": "State is missing.",
        "comment": "Comment is missing."
      }
    }
    """
    When I do PUT /api/v4/alarms/test-alarm-not-exist/changestate:
    """json
    {
      "state": 5,
      "comment": "test-comment"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "state": "State must be one of [0 1 2 3]."
      }
    }
    """

  @concurrent
  Scenario: given change state request and not exist alarm should return error
    When I am admin
    When I do PUT /api/v4/alarms/test-alarm-not-exist/changestate:
    """json
    {
      "state": 0,
      "comment": "test-comment"
    }
    """
    Then the response code should be 404

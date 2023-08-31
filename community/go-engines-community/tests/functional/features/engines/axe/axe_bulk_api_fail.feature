Feature: update alarm
  I need to be able to update alarm

  @concurrent
  Scenario: given ack request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/ack
    Then the response code should be 401

  @concurrent
  Scenario: given ack request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/ack
    Then the response code should be 403

  @concurrent
  Scenario: given ack remove request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/ackremove
    Then the response code should be 401

  @concurrent
  Scenario: given ack remove request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/ackremove
    Then the response code should be 403

  @concurrent
  Scenario: given snooze request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/snooze
    Then the response code should be 401

  @concurrent
  Scenario: given snooze request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/snooze
    Then the response code should be 403

  @concurrent
  Scenario: given cancel request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/cancel
    Then the response code should be 401

  @concurrent
  Scenario: given cancel request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/cancel
    Then the response code should be 403

  @concurrent
  Scenario: given uncancel request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/uncancel
    Then the response code should be 401

  @concurrent
  Scenario: given uncancel request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/uncancel
    Then the response code should be 403

  @concurrent
  Scenario: given assoc ticket request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/assocticket
    Then the response code should be 401

  @concurrent
  Scenario: given assoc ticket request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/assocticket
    Then the response code should be 403

  @concurrent
  Scenario: given comment request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/comment
    Then the response code should be 401

  @concurrent
  Scenario: given comment request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/comment
    Then the response code should be 403

  @concurrent
  Scenario: given change state request and no auth user should not allow access
    When I do PUT /api/v4/bulk/alarms/changestate
    Then the response code should be 401

  @concurrent
  Scenario: given change state request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/alarms/changestate
    Then the response code should be 403

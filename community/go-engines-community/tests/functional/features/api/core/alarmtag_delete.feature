Feature: Delete a alarm tag
  I need to be able to delete a alarm tag
  Only admin should be able to delete a alarm tag

  @concurrent
  Scenario: given delete request should delete tag
    When I am admin
    When I do DELETE /api/v4/alarm-tags/test-alarm-tag-to-delete-1
    Then the response code should be 204
    When I do GET /api/v4/alarm-tags/test-alarm-tag-to-delete-1
    Then the response code should be 404

  @concurrent
  Scenario: given delete external tag request should return error
    When I am admin
    When I do DELETE /api/v4/alarm-tags/test-alarm-tag-to-delete-2
    Then the response code should be 404

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/alarm-tags/notexist
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/alarm-tags/notexist
    Then the response code should be 403

  @concurrent
  Scenario: given invalid delete request should return error
    When I am admin
    When I do DELETE /api/v4/alarm-tags/notexist
    Then the response code should be 404

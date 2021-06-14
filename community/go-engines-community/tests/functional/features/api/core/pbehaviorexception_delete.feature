Feature: Delete pbehavior exception
  I need to be able to delete a pbehavior exception

  Scenario: given delete request should delete exception
    When I am admin
    When I do DELETE /api/v4/pbehavior-exceptions/test-exception-to-delete
    Then the response code should be 204

  Scenario: given delete request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/pbehavior-exceptions/notexist
    Then the response code should be 403

  Scenario: given invalid delete request should return error
    When I am admin
    When I do DELETE /api/v4/pbehavior-exceptions/notexist
    Then the response code should be 404

  Scenario: given request to delete used exception should return
    When I am admin
    When I do DELETE /api/v4/pbehavior-exceptions/test-exception-to-delete-used-by-pbehavior
    Then the response code should be 400
Feature: delete user interface configuration
  I need to be able to delete user interface configuration
  Only admin should be able to delete user interface configuration

  Scenario: DELETE user interface configuration but unauthorized
    When I do DELETE /api/v4/internal/user_interface
    Then the response code should be 401

  Scenario: DELETE user interface configuration but without permissions
    When I am noperms
    When I do DELETE /api/v4/internal/user_interface
    Then the response code should be 403

  Scenario: DELETE user interface configuration with success
    When I am admin
    When I do DELETE /api/v4/internal/user_interface
    Then the response code should be 204
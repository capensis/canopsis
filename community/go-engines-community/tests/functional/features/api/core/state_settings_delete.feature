Feature: Delete state settings
  I need to be able to delete state settings
  Only admin should be able to delete state settings

  @concurrent
  Scenario: given delete request and no auth user should not allow access
    When I do DELETE /api/v4/state-settings/inherited-settings-to-delete
    Then the response code should be 401

  @concurrent
  Scenario: given delete request and auth user without permissions should not allow access
    When I am noperms
    When I do DELETE /api/v4/state-settings/inherited-settings-to-delete
    Then the response code should be 403

  @concurrent
  Scenario: given delete request should return ok
  When I am admin
  When I do DELETE /api/v4/state-settings/inherited-settings-to-delete
  Then the response code should be 204
  When I do GET /api/v4/state-settings/inherited-settings-to-delete
  Then the response code should be 404

  @concurrent
  Scenario: given delete request for not found theme should return error
  When I am admin
  When I do DELETE /api/v4/state-settings/settings-not-found
  Then the response code should be 404
  Then the response body should contain:
  """json
  {
    "error": "Not found"
  }
  """

  @concurrent
  Scenario: given delete request for default theme should return error
  When I am admin
  When I do DELETE /api/v4/state-settings/junit
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "error": "can't delete junit or service settings"
  }
  """

  @concurrent
  Scenario: given delete request for default theme should return error
  When I am admin
  When I do DELETE /api/v4/state-settings/service
  Then the response code should be 400
  Then the response body should contain:
  """json
  {
    "error": "can't delete junit or service settings"
  }
  """

Feature: update a PBehavior
  I need to be able to patch a pbehavior field individually
  Only admin should be able to patch a pbehavior

  @concurrent
  Scenario: given update request and no auth user should not allow access
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  @concurrent
  Scenario: given update request and auth user without view permission should not allow access
    When I am noperms
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  @concurrent
  Scenario: given no exist pbehavior id should return error
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-not-exist:
    """json
    {
      "name": "test-pbehavior-not-exist"
    }
    """
    Then the response code should be 404

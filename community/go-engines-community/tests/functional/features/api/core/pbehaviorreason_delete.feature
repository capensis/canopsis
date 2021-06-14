Feature: delete a reason

  Scenario: DELETE a Reason but unauthorized
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-to-delete
    Then the response code should be 401

  Scenario: DELETE a Reason but without permissions
    When I am noperms
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-to-delete
    Then the response code should be 403

  Scenario: DELETE a Reason with success
    When I am admin
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-to-delete
    Then the response code should be 204

  Scenario: DELETE a Reason with not found response
    When I am admin
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-not-found
    Then the response code should be 404
    Then the response body should be:
    """
    {
      "error": "Not found"
    }
    """

  Scenario: DELETE a Reason, which is linked to pbehaviors
    When I am admin
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-to-delete-linked-to-pbehavior
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "reason is linked to pbehavior"
    }
    """

  Scenario: DELETE a Reason, which is linked to actions
    When I am admin
    When I do DELETE /api/v4/pbehavior-reasons/test-reason-to-delete-linked-to-action
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "error": "reason is linked to action"
    }
    """
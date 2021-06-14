Feature: update a reason
  I need to be able to update a reason
  Only admin should be able to update a reason

  Scenario: PUT a valid reason but unauthorized
    When I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
      {
        "name": "updated_reason_name",
        "description": "updated_reason_description"
      }
    """
    Then the response code should be 401

  Scenario: PUT a valid reason but without permissions
    When I am noperms
    When I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
      {
        "name": "updated_reason_name",
        "description": "updated_reason_description"
      }
    """
    Then the response code should be 403

  Scenario: PUT a valid reason, linked pbehavior should be updated
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-which-reason-should-be-updated
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-pbehavior-which-reason-should-be-updated",
      "reason": {
        "_id": "test-reason-to-update",
        "name": "test-reason-to-update-name",
        "description": "test-reason-to-update-description"
      }
    }
    """
    Then I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
      {
        "name": "updated_reason_name",
        "description": "updated_reason_description"
      }
    """
    Then the response code should be 200
    Then I do GET /api/v4/pbehaviors/test-pbehavior-which-reason-should-be-updated
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "_id": "test-pbehavior-which-reason-should-be-updated",
      "reason": {
        "_id": "test-reason-to-update",
        "name": "updated_reason_name",
        "description": "updated_reason_description"
      }
    }
    """

  Scenario: PUT an invalid reason, where description is absent
    When I am admin
    When I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
      {
         "name": "updated_reason_name_2"
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
      {
        "errors": {
          "description": "Description is missing."
        }
      }
    """

  Scenario: PUT an invalid reason, where name already exists
    When I am admin
    When I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
      {
        "name": "test-reason-1-name",
        "description": "test-reason-1-description"
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
      {
        "errors": {
            "name": "Name already exists."
        }
      }
    """
  Scenario: PUT a valid reason without any changes
    When I am admin
    When I do PUT /api/v4/pbehavior-reasons/test-reason-to-update:
    """
    {
      "name":"updated_reason_name",
      "description":"updated_reason_description"
    }
    """
    Then the response code should be 200
    Then the response body should be:
    """
    {
      "_id": "test-reason-to-update",
      "name": "updated_reason_name",
      "description": "updated_reason_description"
    }
    """

Feature: create a reason
  I need to be able to create a reason
  Only admin should be able to create a reason

  Scenario: POST a valid reason but unauthorized
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "name": "new_reason_name",
        "description": "new_reason_description"
      }
    """
    Then the response code should be 401

  Scenario: POST a valid reason but without permissions
    When I am noperms
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "name": "new_reason_name",
        "description": "new_reason_description"
      }
    """
    Then the response code should be 403

  Scenario: POST a valid reason
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """
    {
      "name": "new_reason_name",
      "description": "new_reason_description"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "name": "new_reason_name",
      "description": "new_reason_description"
    }
    """

  Scenario: POST a valid reason
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "name": "new_reason_name_2",
        "description": "new_reason_description"
      }
    """
    When I do GET /api/v4/pbehavior-reasons/{{ .lastResponse._id}}
    Then the response code should be 200

  Scenario: POST a valid reason with custom id
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "_id": "custom-id",
        "name": "new_reason_name_custom_id",
        "description": "new_reason_description"
      }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-reasons/custom-id
    Then the response code should be 200

  Scenario: POST a valid reason with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "_id": "test-reason-to-update",
        "name": "new_reason_name_custom_id_2",
        "description": "new_reason_description"
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists"
      }
    }
    """

  Scenario: POST an invalid reason, where description is absent
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """
      {
        "name": "new_reason_name_invalid"
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

  Scenario: POST an invalid reason, where name already exists
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
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
            "name": "Name already exists"
        }
      }
    """

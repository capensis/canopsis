Feature: create a reason
  I need to be able to create a reason
  Only admin should be able to create a reason

  Scenario: POST a valid reason but unauthorized
    When I do POST /api/v4/pbehavior-reasons
    Then the response code should be 401

  Scenario: POST a valid reason but without permissions
    When I am noperms
    When I do POST /api/v4/pbehavior-reasons
    Then the response code should be 403

  Scenario: POST a valid reason
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "name": "test-reason-to-create-1-name",
      "description": "test-reason-to-create-1-description"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "name": "test-reason-to-create-1-name",
      "description": "test-reason-to-create-1-description"
    }
    """

  Scenario: POST a valid reason
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "name": "test-reason-to-create-2-name",
      "description": "test-reason-to-create-2-description"
    }
    """
    When I do GET /api/v4/pbehavior-reasons/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "name": "test-reason-to-create-2-name",
      "description": "test-reason-to-create-2-description"
    }
    """

  Scenario: POST a valid reason with custom id
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "_id": "test-reason-to-create-3",
      "name": "test-reason-to-create-3-name",
      "description": "test-reason-to-create-3-description"
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehavior-reasons/test-reason-to-create-3
    Then the response code should be 200

  Scenario: POST a valid reason with custom id that already exist should cause dup error
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "_id": "test-reason-to-check-unique",
      "name": "test-reason-to-create-4-name",
      "description": "test-reason-to-create-4-description"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: POST an invalid reason, where description is absent
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "name": "test-reason-to-create-4-name"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "description": "Description is missing."
      }
    }
    """

  Scenario: POST an invalid reason, where name already exists
    When I am admin
    When I do POST /api/v4/pbehavior-reasons:
    """json
    {
      "name": "test-reason-to-check-unique-name",
      "description": "test-reason-to-create-4-description"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
          "name": "Name already exists."
      }
    }
    """

Feature: Update an associative table
  I need to be able to update an associative table
  Only admin should be able to update an associative table

  Scenario: given create request should create associative table
    When I am admin
    Then I do POST /api/v4/associativetable:
    """
    {
      "name": "test-associativetable-to-create",
      "content": [
        {
          "title": "test-associativetable-to-create-content-nested-val-1",
          "names": [
            "test-associativetable-to-create-content-nested-val-2",
            "test-associativetable-to-create-content-nested-val-3"
          ],
          "_id": "test-associativetable-to-create-content-nested-val-4"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-associativetable-to-create",
      "content": [
        {
          "title": "test-associativetable-to-create-content-nested-val-1",
          "names": [
            "test-associativetable-to-create-content-nested-val-2",
            "test-associativetable-to-create-content-nested-val-3"
          ],
          "_id": "test-associativetable-to-create-content-nested-val-4"
        }
      ]
    }
    """

  Scenario: given update request should update associative table
    When I am admin
    Then I do POST /api/v4/associativetable:
    """
    {
      "name": "test-associativetable-to-update",
      "content": {
        "test-associativetable-to-update-key-updated": "test-associativetable-to-update-val-updated"
      }
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "name": "test-associativetable-to-update",
      "content": {
        "test-associativetable-to-update-key-updated": "test-associativetable-to-update-val-updated"
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do POST /api/v4/associativetable
    Then the response code should be 401

  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/associativetable
    Then the response code should be 403

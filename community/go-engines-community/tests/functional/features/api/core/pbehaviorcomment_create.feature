Feature: create a PBehavior comment
  I need to be able to create a PBehavior comment

  Scenario: Given new pbehavior should add comment to pbehavior
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-comment-create",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type": "test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-comment-create-pattern"
            }
          }
        ]
      ]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "{{ .lastResponse._id }}",
      "message": "Test message"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "message": "Test message"
    }
    """

  Scenario: Given pbehavior should add comment to pbehavior
    When I am admin
    When I do POST /api/v4/pbehavior-comments:
    """json
    {
      "pbehavior": "test-pbehavior-to-create-comment",
      "message": "test-create-comment-message"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-create-comment
    Then the response body should contain:
    """json
    {
      "comments": [
        {
          "author": {
            "_id": "root",
            "name": "root"
          }
        },
        {
          "author": {
            "_id": "root",
            "name": "root"
          }
        },
        {
          "author": {
            "_id": "root",
            "name": "root"
          },
          "message": "test-create-comment-message"
        }
      ]
    }
    """

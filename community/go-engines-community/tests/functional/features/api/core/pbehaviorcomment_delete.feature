Feature: delete a PBehavior comment
  I need to be able to delete a PBehavior comment

  Scenario: Given new pbehavior Should delete comment from pbehavior
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """json
    {
      "enabled": true,
      "name":" test-pbehavior-to-comment-delete",
      "tstart": 1591172881,
      "tstop": 1591536400,
      "color": "#FFFFFF",
      "type":"test-type-to-pbh-edit-1",
      "reason": "test-reason-to-pbh-edit",
      "entity_pattern": [
        [
          {
            "field": "name",
            "cond": {
              "type": "eq",
              "value": "test-pbehavior-to-comment-delete-pattern"
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
    When I do DELETE /api/v4/pbehavior-comments/{{ .lastResponse._id }}
    Then the response code should be 204

  Scenario: Given pbehavior Should delete comment from pbehavior
    When I am admin
    When I do DELETE /api/v4/pbehavior-comments/test-pbehavior-to-delete-comment-1
    Then the response code should be 204
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-delete-comment
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-delete-comment",
      "comments": [
        {
          "_id": "test-pbehavior-to-delete-comment-2",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "message": "test-pbehavior-to-delete-comment-2-message",
          "ts": 1592215337
        }
      ]
    }
    """

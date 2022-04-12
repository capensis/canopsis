Feature: create a PBehavior comment
  I need to be able to create a PBehavior comment

  Scenario: Given new pbehavior Should add comment to pbehavior
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled":true,
      "name":"test-pbehavior-to-comment-create",
      "tstart":1591172881,
      "tstop":1591536400,
      "color": "#FFFFFF",
      "type":"test-type-to-pbh-edit-1",
      "reason":"test-reason-1",
      "filter":{
        "$and":[
           {
              "name": "test filter"
           }
        ]
      },
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    When I do POST /api/v4/pbehavior-comments:
    """
    {
      "pbehavior": "{{ .lastResponse._id }}",
      "message": "Test message"
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "author": "root",
      "message": "Test message"
    }
    """

  Scenario: Given pbehavior Should add comment to pbehavior
    When I am admin
    When I do POST /api/v4/pbehavior-comments:
    """
    {
      "pbehavior": "test-pbehavior-to-create-comment",
      "message": "test-create-comment-message"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-create-comment
    Then the response body should contain:
    """
    {
      "comments": [
        {
          "author": "root"
        },
        {
          "author": "root"
        },
        {
          "author": "root",
          "message": "test-create-comment-message"
        }
      ]
    }
    """

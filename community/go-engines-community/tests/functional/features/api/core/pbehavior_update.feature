Feature: update a PBehavior
  I need to be able to update a PBehavior
  Only admin should be able to update a PBehavior

  Scenario: PUT a valid PBehavior but unauthorized
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update:
    """
      {
        "enabled":true,
        "name":"test-pbehavior-to-update-name",
        "tstart":1591172881,
        "tstop":1591536400,
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
        ]
      }
    """
    Then the response code should be 401

  Scenario: PUT a valid PBehavior but without permissions
    When I am noperms
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update:
    """
      {
        "enabled":true,
        "name":"test-pbehavior-to-update-name",
        "tstart":1591172881,
        "tstop":1591536400,
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
        ]
      }
    """
    Then the response code should be 403

  Scenario: PUT a valid PBehavior
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update:
    """
      {
        "enabled":true,
        "name":"test-pbehavior-to-update-name",
        "tstart":1591172881,
        "tstop":1591536400,
        "type":"test-type-to-pbh-edit-1",
        "reason":"test-reason-1",
        "filter":{
          "$and":[
             {
                "name": "test filter"
             }
          ]
        },
        "exdates":[],
        "exceptions": []
      }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "_id": "test-pbehavior-to-update",
        "author":"root",
        "created": 1592215337,
        "enabled":true,
        "exceptions": [],
        "reason": {
          "_id": "test-reason-1",
          "description": "test-reason-1-description",
          "name": "test-reason-1-name"
        },
        "filter":{
          "$and":[
             {
                "name": "test filter"
             }
          ]
        },
        "exdates":[],
        "comments": [
          {
            "_id": "test-comment-1",
            "author": "root",
            "ts": 1592215337,
            "message": "qwerty"
          },
          {
            "_id": "test-comment-2",
            "author": "root",
            "ts": 1592215337,
            "message": "asdasd"
          }
        ],
        "name":"test-pbehavior-to-update-name",
        "rrule": "",
        "tstart":1591172881,
        "tstop":1591536400,
        "type": {
          "_id": "test-type-to-pbh-edit-1",
          "description": "Pbh edit 1 State type",
          "icon_name": "test-to-pbh-edit-1-icon",
          "name": "Pbh edit 1 State",
          "priority": 10,
          "type": "active"
        }
      }
    """

  Scenario: PUT an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """
    {
      "tstart":1592172881,
      "tstop":1591536400
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """

  Scenario: PUT an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """
    {
      "exdates":[
        {
          "begin": 1592164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """

  Scenario: PUT an invalid PBehavior with not valid reason
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """
    {
      "reason":"notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """
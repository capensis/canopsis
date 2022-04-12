Feature: update a PBehavior
  I need to be able to update a PBehavior
  Only admin should be able to update a PBehavior

  Scenario: PUT a valid PBehavior but unauthorized
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  Scenario: PUT a valid PBehavior but without permissions
    When I am noperms
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  Scenario: PUT a valid PBehavior with tstop:null and type pause
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update:
    """json
    {
      "enabled":true,
      "name":"test-pbehavior-to-update-name",
      "tstart":1591172881,
      "tstop":null,
      "color": "#FFFFFF",
      "type":"test-default-pause-type",
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
    """json
    {
      "_id": "test-pbehavior-to-update",
      "author":"root",
      "created": 1592215337,
      "color": "#FFFFFF",
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
      "name":"test-pbehavior-to-update-name",
      "rrule": "",
      "tstart":1591172881,
      "tstop":null,
      "type": {
        "_id": "test-default-pause-type",
        "name": "Default Type Pause",
        "description": "Default Type Pause",
        "type": "pause",
        "priority": 3,
        "icon_name": "test-pause-icon"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update",
      "author":"root",
      "created": 1592215337,
      "color": "#FFFFFF",
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
      "tstop":null,
      "type": {
        "_id": "test-default-pause-type",
        "name": "Default Type Pause",
        "description": "Default Type Pause",
        "type": "pause",
        "priority": 3,
        "icon_name": "test-pause-icon"
      }
    }
    """

  Scenario: PUT a valid PBehavior
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update:
    """json
    {
      "enabled":true,
      "name":"test-pbehavior-to-update-name",
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
      "exdates":[],
      "exceptions": []
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "_id": "test-pbehavior-to-update",
      "author":"root",
      "created": 1592215337,
      "color": "#FFFFFF",
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
      },
      "last_alarm_date": 1592215337
    }
    """

  Scenario: PUT an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """json
    {
      "tstart":1592172881,
      "tstop":1591536400
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """

  Scenario: PUT an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """json
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
    """json
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """

  Scenario: PUT an invalid PBehavior with not valid reason
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-invalid:
    """json
    {
      "reason":"notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """

  Scenario: POST a valid PBehavior with the name that already exists should cause dup error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-update-1:
    """json
    {
      "enabled":true,
      "name": "test-pbehavior-to-check-unique-name",
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
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "name": "Name already exists."
      }
    }
    """

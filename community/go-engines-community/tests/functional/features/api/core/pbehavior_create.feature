Feature: create a PBehavior
  I need to be able to create a PBehavior
  Only admin should be able to create a PBehavior

  Scenario: POST a valid PBehavior but unauthorized
    When I do POST /api/v4/pbehaviors
    Then the response code should be 401

  Scenario: POST a valid PBehavior but without permissions
    When I am noperms
    When I do POST /api/v4/pbehaviors
    Then the response code should be 403

  Scenario: POST a valid PBehavior
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled":true,
      "name":"test-pbehavior-to-create-1",
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
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    Then the response body should contain:
    """
    {
      "enabled":true,
      "author":"root",
      "name":"test-pbehavior-to-create-1",
      "tstart":1591172881,
      "tstop":1591536400,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      },
      "reason": {
        "_id": "test-reason-1"
      },
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
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ],
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  Scenario: POST a valid PBehavior
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled":true,
      "name":"test-pbehavior-to-create-2",
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
      ],
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors/{{ .lastResponse._id}}
    Then the response code should be 200
    Then the response body should contain:
    """
    {
      "author": "root",
      "comments": null,
      "enabled": true,
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit",
          "created": 1592215037,
          "description": "test",
          "exdates": [
            {
              "begin": 15911648001,
              "end": 1591167901,
              "type": {
                "_id": "test-type-to-pbh-edit-1",
                "description": "Pbh edit 1 State type",
                "icon_name": "test-to-pbh-edit-1-icon",
                "name": "Pbh edit 1 State",
                "priority": 10,
                "type": "active"
              }
            }
          ],
          "name": "Exception to pbehavior edit"
        }
      ],
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1",
            "description": "Pbh edit 1 State type",
            "icon_name": "test-to-pbh-edit-1-icon",
            "name": "Pbh edit 1 State",
            "priority": 10,
            "type": "active"
          }
        }
      ],
      "filter": {
        "$and": [
          {
            "name": "test filter"
          }
        ]
      },
      "name": "test-pbehavior-to-create-2",
      "reason": {
        "_id": "test-reason-1",
        "description": "test-reason-1-description",
        "name": "test-reason-1-name"
      },
      "rrule": "",
      "tstart": 1591172881,
      "tstop": 1591536400,
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

  Scenario: POST a valid PBehavior with custom ID
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
      {
        "_id": "custom-id",
        "enabled":true,
        "name": "test-pbehavior-to-create-3",
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
        ],
        "exceptions": ["test-exception-to-pbh-edit"]
      }
    """
    Then the response code should be 201
    When I do GET /api/v4/pbehaviors/custom-id
    Then the response code should be 200

  Scenario: POST a valid PBehavior with the custom ID that already exists should cause dup error
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
      {
        "_id": "test-pbehavior-to-update",
        "enabled":true,
        "name": "test-pbehavior-to-create-4",
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
        ],
        "exceptions": ["test-exception-to-pbh-edit"]
      }
    """
    Then the response code should be 400
    Then the response body should be:
    """
    {
      "errors": {
        "_id": "ID already exists."
      }
    }
    """

  Scenario: POST a valid pause PBehavior without Stop
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "enabled":true,
      "name": "test-pbehavior-to-create-5",
      "tstart":1591172881,
      "type":"test-type-to-pbh-edit-3",
      "reason":"test-reason-1",
      "filter":{
        "$and":[
           {
              "name":"ccccc"
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
    Then the response body should contain:
    """
    {
      "enabled":true,
      "author":"root",
      "name": "test-pbehavior-to-create-5",
      "tstart":1591172881,
      "tstop":null,
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      },
      "reason": {
        "_id": "test-reason-1"
      },
      "filter":{
        "$and":[
           {
              "name":"ccccc"
           }
        ]
      },
      "exdates":[
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ],
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  Scenario: POST an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do POST /api/v4/pbehaviors:
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

  Scenario: POST an invalid PBehavior with not existed reason
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "reason": "notexist"
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

  Scenario: POST an invalid PBehavior, where tstart > tstop
    When I am admin
    When I do POST /api/v4/pbehaviors:
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

  Scenario: POST an invalid PBehavior with not existed type
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "type":"notexist"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "type": "Type doesn't exist."
      }
    }
    """

  Scenario: POST an invalid PBehavior, where filter is invalid
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "filter": "{}"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
      {
        "errors": {
          "filter": "Filter is invalid entity filter."
        }
      }
    """

  Scenario: POST an invalid PBehavior, where Stop is missing
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "tstart":1591172881,
      "type":"test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """

  Scenario: POST an invalid PBehavior, where ID is invalid
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "_id": "invalid/id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """
    {
      "_id": "invalidid?key=value"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """
    {
      "_id": "$invalidid"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """
    When I do POST /api/v4/pbehaviors:
    """
    {
      "_id": "invalid.id"
    }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
    {
      "errors": {
        "_id": "ID cannot contain /?.$ characters."
      }
    }
    """

  Scenario: POST a valid PBehavior with strange ids
    When I am admin
    When I do POST /api/v4/pbehaviors:
    """
    {
      "_id": "strange \\id&key=value!*@!'\"-_:;<>",
      "enabled":true,
      "name":"test-pbehavior-to-create-6",
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
      }
    }
    """
    Then the response code should be 201
    When I do DELETE /api/v4/pbehaviors/{{ .lastResponse._id}}
    Then the response code should be 204
Feature: update a pbehavior
  I need to be able to patch a pbehavior field individually
  Only admin should be able to patch a pbehavior

  @concurrent
  Scenario: given update request with start, stop and type should update only if all fields are valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": null,
      "tstop": 1591172881,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": 1591172681,
      "tstop": null,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": 1591172881,
      "tstop": 1591172681,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": 1591172881,
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": 1591172881,
      "tstop": 1591172681,
      "type": "test-type-to-pbh-edit-3"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop should be greater than Start."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-1:
    """json
    {
      "tstart": 1591172681,
      "tstop": null,
      "type": "test-type-to-pbh-edit-3"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-1
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "tstart": 1591172681,
      "tstop": null,
      "type": {
        "_id": "test-type-to-pbh-edit-3"
      }
    }
    """

  @concurrent
  Scenario: given update request with type should update only if type is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-2:
    """json
    {
      "type": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "type": {
        "_id": "test-type-to-pbh-edit-1"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-2:
    """json
    {
      "type": "test-type-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "type": "Type doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-2:
    """json
    {
      "type": "test-type-to-pbh-edit-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-2
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """

  @concurrent
  Scenario: given update request with type for pbehavior without stop should update only if type is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-3:
    """json
    {
      "type": "test-type-to-pbh-edit-1"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "tstop": "Stop is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-3:
    """json
    {
      "type": "test-type-to-pbh-edit-2"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-3
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "type": {
        "_id": "test-type-to-pbh-edit-2"
      }
    }
    """

  @concurrent
  Scenario: given update request with reason should update only if reason is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-4:
    """json
    {
      "reason": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-4:
    """json
    {
      "reason": "test-reason-not-exist"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "reason": "Reason doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-4:
    """json
    {
      "reason": "test-reason-to-pbh-edit"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-4
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "reason": {
        "_id": "test-reason-to-pbh-edit"
      }
    }
    """

  @concurrent
  Scenario: given update request with enabled should update only if enabled is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-5:
    """json
    {
      "enabled": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "enabled": true
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-5:
    """json
    {
      "enabled": false
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-5
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "enabled": false
    }
    """

  @concurrent
  Scenario: given update request with rrule should update only if rrule is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-6:
    """json
    {
      "rrule": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "rrule": "FREQ=DAILY"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-6:
    """json
    {
      "rrule": "invalid"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "rrule": "RRule is invalid recurrence rule."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-6:
    """json
    {
      "rrule": "FREQ=WEEKLY"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": "FREQ=WEEKLY"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": "FREQ=WEEKLY"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-6:
    """json
    {
      "rrule": ""
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": ""
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-6
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "rrule": ""
    }
    """

  @concurrent
  Scenario: given update request with exdates should update only if exdates is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          }
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": [{}]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates.0.begin": "Begin is missing.",
        "exdates.0.end": "End is missing.",
        "exdates.0.type": "Type is missing."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": [
        {
          "begin": 1591164001,
          "end": 1591167601,
          "type": "test-type-not-exist"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates": "Exdates doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": [
        {
          "begin": 1592164001,
          "end": 1591167601,
          "type": "test-type-to-pbh-edit-1"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exdates.0.end": "End should be greater than Begin."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": "test-type-to-pbh-edit-2"
        }
      ]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": {
            "_id": "test-type-to-pbh-edit-2"
          }
        }
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": [
        {
          "begin": 1591164002,
          "end": 1591167602,
          "type": {
            "_id": "test-type-to-pbh-edit-2"
          }
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-7:
    """json
    {
      "exdates": []
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": []
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-7
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exdates": []
    }
    """

  @concurrent
  Scenario: given update request with exceptions should update only if exceptions is empty or is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-8:
    """json
    {
      "exceptions": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-8:
    """json
    {
      "exceptions": [""]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-8:
    """json
    {
      "exceptions": ["test-exception-not-exist"]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "exceptions": "Exceptions doesn't exist."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-8:
    """json
    {
      "exceptions": []
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": []
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": []
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-8:
    """json
    {
      "exceptions": ["test-exception-to-pbh-edit"]
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-8
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "author": {
        "_id": "root",
        "name": "root"
      },
      "exceptions": [
        {
          "_id": "test-exception-to-pbh-edit"
        }
      ]
    }
    """

  @concurrent
  Scenario: given update request with color should update only if color is not empty and is valid
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-9:
    """json
    {
      "color": null
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#FFFFFF"
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-9:
    """json
    {
      "color": "test"
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "color": "Color is not valid."
      }
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-9:
    """json
    {
      "color": ""
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": ""
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": ""
    }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-second-9:
    """json
    {
      "color": "#AAAAAA"
    }
    """
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#AAAAAA"
    }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-second-9
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "color": "#AAAAAA"
    }
    """

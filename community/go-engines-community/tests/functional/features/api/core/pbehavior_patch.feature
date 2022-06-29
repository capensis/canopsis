Feature: update a PBehavior
  I need to be able to patch a PBehavior field individually
  Only admin should be able to patch a PBehavior

  Scenario: PATCH a valid PBehavior but unauthorized
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 401

  Scenario: PATCH a valid PBehavior but without permissions
    When I am noperms
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-update
    Then the response code should be 403

  Scenario: PATCH a non-existing PBehavior
    When I am admin
    When I do PATCH /api/v4/pbehaviors/non-existing-pbehavior:
    """
      {
        "name": "non-existing-pbehavior-name-new"
      }
    """
    Then the response code should be 404

  Scenario: PATCH a valid PBehavior with new name value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "name": "test-pbehavior-to-patch-1-name-new"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "name": "test-pbehavior-to-patch-1-name-new",
        "author": "root"
      }
    """

  Scenario: PATCH a valid PBehavior with new enabled value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "enabled": false
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "enabled": false
      }
    """

  Scenario: PATCH a valid PBehavior with new color value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "color": "#FFFFFF"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "color": "#FFFFFF"
      }
    """
  
  Scenario: PATCH a valid PBehavior with new filter value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "filter": "{\"$or\":[{\"name\":\"test-new-filter\"}]}"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "filter":{
          "$and":[
            {
              "name": "another test filter"
            }
          ]
        }
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "filter":{
          "$and":[
            {
              "name": "another test filter"
            }
          ]
        }
      }
    """
  
  Scenario: PATCH a valid PBehavior with new reason value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "reason": "test-reason-to-patch-pbehavior-new"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "reason": {
          "_id": "test-reason-to-patch-pbehavior-new",
          "name": "test-reason-to-patch-pbehavior-new-name",
          "description": "test-reason-to-patch-pbehavior-new-description"
        }
      }
    """

  Scenario: PATCH a valid PBehavior with new rrule value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "rrule": "FREQ=YEARLY"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "rrule": "FREQ=YEARLY"
      }
    """

  Scenario: PATCH a valid PBehavior with new exdates value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "exdates": [
            {
              "begin": 1111111117,
              "end": 1111111118,
              "type": "test-pause-type-to-patch-pbehavior"
            }
          ]
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "exdates": [
          {
            "begin": 1111111117,
            "end": 1111111118,
            "type": {
              "_id": "test-pause-type-to-patch-pbehavior",
              "name": "test-pause-type-to-patch-pbehavior-name",
              "description": "test-pause-type-to-patch-pbehavior-description",
              "type": "pause",
              "priority": 29,
              "icon_name": "test-pause-type-to-patch-pbehavior-icon"
            }
          }
        ]
      }
    """

  Scenario: PATCH a valid PBehavior with pause type and null stop to the new stop value less than start
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """
      {
          "tstop": 1592215336
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
      {
        "error": "invalid fields start, stop, type"
      }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """
      {
          "tstop": -300
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
      {
        "error": "invalid fields start, stop, type"
      }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-2:
    """
      {
          "tstop": 1592215339
      }
    """
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstop": 1592215339
      }
    """
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-2
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstop": 1592215339
      }
    """

  Scenario: PATCH a valid PBehavior with new exceptions value
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
          "exceptions": ["test-exception-to-patch-pbehavior-new"]
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "exceptions": [
          {
            "_id": "test-exception-to-patch-pbehavior-new",
            "created": 1592215337,
            "description": "test-exception-to-patch-pbehavior-new-description",
            "exdates": [
              {
                "begin": 1592215337,
                "end": 1592215337,
                "type": {
                  "_id": "test-pause-type-to-patch-pbehavior",
                  "description": "test-pause-type-to-patch-pbehavior-description",
                  "icon_name": "test-pause-type-to-patch-pbehavior-icon",
                  "name": "test-pause-type-to-patch-pbehavior-name",
                  "priority": 29,
                  "type": "pause"
                }
              }
            ],
            "name": "test-exception-to-patch-pbehavior-new-name"
          }
        ]
      }
    """

  Scenario: PATCH a valid PBehavior with invalid tstart/tstop and type pause
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstart": 1111111112,
        "tstop": 1111111111,
        "type": "test-pause-type-to-patch-pbehavior"
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
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstart": 1111111112,
        "tstop": 1111111113,
        "type": "test-pause-type-to-patch-pbehavior"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstart": 1111111112,
        "tstop": 1111111113,
        "type": {
          "_id": "test-pause-type-to-patch-pbehavior",
          "name": "test-pause-type-to-patch-pbehavior-name",
          "description": "test-pause-type-to-patch-pbehavior-description",
          "type": "pause",
          "priority": 29,
          "icon_name": "test-pause-type-to-patch-pbehavior-icon"
        }
      }
    """
  
  Scenario: PATCH a valid PBehavior having type pause with null tstop
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstop": null
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstop": null
      }
    """

  Scenario: PATCH a valid PBehavior having a valid tstart/tstop with type active
    When I am admin
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstart": 1111111112,
        "tstop": null,
        "type": {
          "type": "pause"
        }
      }
    """
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "type": "test-active-type-to-patch-pbehavior"
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

  Scenario: PATCH a valid PBehavior with valid tstart/tstop and type active
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstart": 1111111111,
        "tstop": 1111111112,
        "type": "test-active-type-to-patch-pbehavior"
      }
    """
    Then the response code should be 200
    When I do GET /api/v4/pbehaviors/test-pbehavior-to-patch-1
    Then the response code should be 200
    Then the response body should contain:
    """
      {
        "tstart": 1111111111,
        "tstop": 1111111112,
        "type": {
          "_id": "test-active-type-to-patch-pbehavior",
          "name": "test-active-type-to-patch-pbehavior-name",
          "description": "test-active-type-to-patch-pbehavior-description",
          "type": "active",
          "priority": 28,
          "icon_name": "test-active-type-to-patch-pbehavior-icon"
        }
      }
    """

  Scenario: PATCH a valid PBehavior having type active with null tstop
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstop": null
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
      {
        "error": "invalid fields start, stop, type"
      }
    """

  Scenario: PATCH a valid PBehavior having type active with invalid tstop
    When I am admin
    When I do PATCH /api/v4/pbehaviors/test-pbehavior-to-patch-1:
    """
      {
        "tstop": 1000000000
      }
    """
    Then the response code should be 400
    Then the response body should contain:
    """
      {
        "error": "invalid fields start, stop, type"
      }
    """

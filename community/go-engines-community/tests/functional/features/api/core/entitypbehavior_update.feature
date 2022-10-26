Feature: update a pbehavior
  I need to be able to update a pbehavior
  Only admin should be able to update a pbehavior

  Scenario: given update entity pbehavior request should return error
    When I am admin
    When I do PUT /api/v4/pbehaviors/test-pbehavior-to-entity-update-1:
    """json
    {
      "enabled": true,
      "name": "test-pbehavior-to-entity-update-1-name",
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
              "value": "test-pbehavior-to-entity-update-1-pattern"
            }
          }
        ]
      ],
      "exdates":[],
      "exceptions": []
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "Cannot update a pbehavior with origin."
      }
    }
    """

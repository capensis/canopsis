Feature: get a pbehavior by entity id
  I need to be able to get pbehavior for entity
  Only admin should be able to get a PBehavior

  Scenario: given get request should return pbehaviors
    When I am admin
    When I do GET /api/v4/entities/pbehaviors?id=test-pbehavior-get-by-entity-id-resource/test-pbehavior-get-by-entity-id-component
    Then the response code should be 200
    Then the response body should contain:
    """
    [
      {
        "_id": "test-pbehavior-to-get-by-entity-id",
        "author": "root",
        "comments": [],
        "created": 1592215337,
        "updated": 1592215337,
        "enabled": true,
        "exceptions": [],
        "exdates": [],
        "filter": {
          "$and": [
            {
              "name": "test-pbehavior-get-by-entity-id-resource"
            }
          ]
        },
        "name": "Pbehavior by entity id",
        "reason": {
          "_id": "test-reason-1",
          "name": "test-reason-1-name",
          "description": "test-reason-1-description"
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
    ]
    """

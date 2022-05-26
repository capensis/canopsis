Feature: get a pbehavior by entity id
  I need to be able to get pbehavior for entity
  Only admin should be able to get a PBehavior

  Scenario: given get request should return pbehaviors
    When I am admin
    When I do GET /api/v4/entities/pbehaviors?id=test-resource-to-pbehavior-get-by-entity-id-1/test-component-default
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "_id": "test-pbehavior-to-get-by-entity-id-1",
        "author": "root",
        "comments": [],
        "created": 1592215337,
        "updated": 1592215337,
        "enabled": true,
        "exceptions": [],
        "exdates": [],
        "entity_pattern": [
          [
            {
              "field": "name",
              "cond": {
                "type": "eq",
                "value": "test-resource-to-pbehavior-get-by-entity-id-1"
              }
            }
          ]
        ],
        "name": "test-pbehavior-to-get-by-entity-id-1-name",
        "reason": {
          "_id": "test-reason-to-pbh-edit",
          "name": "test-reason-to-pbh-edit-name",
          "description": "test-reason-to-pbh-edit-description"
        },
        "rrule": "",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": {
          "_id": "test-type-to-pbh-edit-1",
          "description": "Pbh edit 1 State type",
          "icon_name": "test-to-pbh-edit-1-icon",
          "name": "Pbh edit 1 State",
          "type": "active"
        }
      }
    ]
    """

  Scenario: given get request should return pbehaviors with old mongo query
    When I am admin
    When I do GET /api/v4/entities/pbehaviors?id=test-resource-to-pbehavior-get-by-entity-id-2/test-component-default
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "_id": "test-pbehavior-to-get-by-entity-id-2",
        "author": "root",
        "comments": [],
        "created": 1592215337,
        "updated": 1592215337,
        "enabled": true,
        "exceptions": [],
        "exdates": [],
        "old_mongo_query": {
          "$and": [{"name": "test-resource-to-pbehavior-get-by-entity-id-2"}]
        },
        "name": "test-pbehavior-to-get-by-entity-id-2-name",
        "reason": {
          "_id": "test-reason-to-pbh-edit",
          "name": "test-reason-to-pbh-edit-name",
          "description": "test-reason-to-pbh-edit-description"
        },
        "rrule": "",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": {
          "_id": "test-type-to-pbh-edit-1",
          "description": "Pbh edit 1 State type",
          "icon_name": "test-to-pbh-edit-1-icon",
          "name": "Pbh edit 1 State",
          "type": "active"
        }
      }
    ]
    """

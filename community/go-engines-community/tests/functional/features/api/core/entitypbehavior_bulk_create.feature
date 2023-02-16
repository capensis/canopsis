Feature: Bulk create pbehaviors
  I need to be able to create multiple pbehaviors
  Only admin should be able to create multiple pbehaviors

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/entity-pbehaviors
    Then the response code should be 401

  Scenario: given bulk create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/entity-pbehaviors
    Then the response code should be 403

  Scenario: given bulk create request should return multi status and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/entity-pbehaviors:
    """json
    [
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-create-1-origin",
        "name": "test-pbehavior-to-bulk-entity-create-1-name",
        "comment": "test-pbehavior-to-bulk-entity-create-1-comment",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "color": "#FFFFFF",
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-to-pbh-edit"
      },
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-create-1-origin",
        "name": "test-pbehavior-to-bulk-entity-create-1-another-name",
        "comment": "test-pbehavior-to-bulk-entity-create-1-comment",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "color": "#FFFFFF",
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-to-pbh-edit"
      },
      {},
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-create-2-origin",
        "name": "test-pbehavior-to-check-unique-name",
        "comment": "test-pbehavior-to-bulk-entity-create-2-comment",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "color": "#FFFFFF",
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-to-pbh-edit"
      },
      [],
      {
        "entity": "test-entity-to-bulk-entity-pbehavior-create-2/test-component-default",
        "origin": "test-pbehavior-to-bulk-entity-create-2-origin",
        "name": "test-pbehavior-to-bulk-entity-create-2-name",
        "comment": "test-pbehavior-to-bulk-entity-create-2-comment",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "color": "",
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-to-pbh-edit"
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-create-1-origin",
          "name": "test-pbehavior-to-bulk-entity-create-1-name",
          "comment": "test-pbehavior-to-bulk-entity-create-1-comment",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "#FFFFFF",
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-to-pbh-edit"
        }
      },
      {
        "status": 400,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-create-1-origin",
          "name": "test-pbehavior-to-bulk-entity-create-1-another-name",
          "comment": "test-pbehavior-to-bulk-entity-create-1-comment",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "#FFFFFF",
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-to-pbh-edit"
        },
        "errors": {
          "entity": "Pbehavior for origin already exists."
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "entity": "Entity is missing.",
          "origin": "Origin is missing.",
          "name": "Name is missing.",
          "tstart": "Start is missing.",
          "reason": "Reason is missing.",
          "type": "Type is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-create-2-origin",
          "name": "test-pbehavior-to-check-unique-name",
          "comment": "test-pbehavior-to-bulk-entity-create-2-comment",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "#FFFFFF",
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-to-pbh-edit"
        },
        "errors": {
          "name": "Name already exists."
        }
      },
      {
        "status": 400,
        "item": [],
        "error": "value doesn't contain object; it contains array"
      },
      {
        "status": 200,
        "item": {
          "entity": "test-entity-to-bulk-entity-pbehavior-create-2/test-component-default",
          "origin": "test-pbehavior-to-bulk-entity-create-2-origin",
          "name": "test-pbehavior-to-bulk-entity-create-2-name",
          "comment": "test-pbehavior-to-bulk-entity-create-2-comment",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "",
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-to-pbh-edit"
        }
      }
    ]
    """
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-bulk-entity-create&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "enabled": true,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "name": "test-pbehavior-to-bulk-entity-create-1-name",
          "origin": "test-pbehavior-to-bulk-entity-create-1-origin",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "#FFFFFF",
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          },
          "reason": {
            "_id": "test-reason-to-pbh-edit"
          },
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-bulk-entity-pbehavior-create-1/test-component-default"
                }
              }
            ]
          ],
          "comments": [
            {
              "author": {
                "_id": "root",
                "name": "root"
              },
              "message": "test-pbehavior-to-bulk-entity-create-1-comment"
            }
          ]
        },
        {
          "enabled": true,
          "author": {
            "_id": "root",
            "name": "root"
          },
          "name": "test-pbehavior-to-bulk-entity-create-2-name",
          "origin": "test-pbehavior-to-bulk-entity-create-2-origin",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "color": "",
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          },
          "reason": {
            "_id": "test-reason-to-pbh-edit"
          },
          "entity_pattern": [
            [
              {
                "field": "_id",
                "cond": {
                  "type": "eq",
                  "value": "test-entity-to-bulk-entity-pbehavior-create-2/test-component-default"
                }
              }
            ]
          ],
          "comments": [
            {
              "author": {
                "_id": "root",
                "name": "root"
              },
              "message": "test-pbehavior-to-bulk-entity-create-2-comment"
            }
          ]
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

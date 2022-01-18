Feature: Bulk create pbehaviors
  I need to be able to create multiple pbehaviors
  Only admin should be able to create multiple pbehaviors

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/pbehaviors
    Then the response code should be 401

  Scenario: given bulk create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/pbehaviors
    Then the response code should be 403

  Scenario: given bulk create request should return multi status and should be handled independently
    When I am admin
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-create-1-1",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-create-1-1-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter": {
          "$and": [
            {
              "name": "test-pbehavior-to-bulk-create-1-1-filter"
            }
          ]
        },
        "exdates": [
          {
            "begin": 1591164001,
            "end": 1591167601,
            "type": "test-type-to-pbh-edit-1"
          }
        ],
        "exceptions": ["test-exception-to-pbh-edit"]
      },
      {
        "_id": "test-pbehavior-to-bulk-create-1-1",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-create-1-1-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter": {
          "$and": [
            {
              "name": "test-pbehavior-to-bulk-create-1-1-filter"
            }
          ]
        },
        "exdates": [
          {
            "begin": 1591164001,
            "end": 1591167601,
            "type": "test-type-to-pbh-edit-1"
          }
        ],
        "exceptions": ["test-exception-to-pbh-edit"]
      },
      {},
      {
        "_id": "test-pbehavior-to-check-unique"
      },
      {
        "name": "test-pbehavior-to-check-unique-name"
      },
      [],
      {
        "_id": "test-pbehavior-to-bulk-create-1-2",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-create-1-2-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter": {
          "$and": [
            {
              "name": "test-pbehavior-to-bulk-create-1-2-filter"
            }
          ]
        },
        "exdates": [
          {
            "begin": 1591164001,
            "end": 1591167601,
            "type": "test-type-to-pbh-edit-1"
          }
        ],
        "exceptions": ["test-exception-to-pbh-edit"]
      }
    ]
    """
    Then the response code should be 207
    Then the response body should contain:
    """json
    [
      {
        "status": 200,
        "id": "test-pbehavior-to-bulk-create-1-1",
        "item": {
          "_id": "test-pbehavior-to-bulk-create-1-1",
          "enabled": true,
          "name": "test-pbehavior-to-bulk-create-1-1-name",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-1",
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-bulk-create-1-1-filter"
              }
            ]
          },
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": "test-type-to-pbh-edit-1"
            }
          ],
          "exceptions": ["test-exception-to-pbh-edit"]
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-pbehavior-to-bulk-create-1-1",
          "enabled": true,
          "name": "test-pbehavior-to-bulk-create-1-1-name",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-1",
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-bulk-create-1-1-filter"
              }
            ]
          },
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": "test-type-to-pbh-edit-1"
            }
          ],
          "exceptions": ["test-exception-to-pbh-edit"]
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {},
        "errors": {
          "enabled": "Enabled is missing.",
          "name": "Name is missing.",
          "filter": "Filter is missing.",
          "tstart": "Start is missing.",
          "reason": "Reason is missing.",
          "type": "Type is missing."
        }
      },
      {
        "status": 400,
        "item": {
          "_id": "test-pbehavior-to-check-unique"
        },
        "errors": {
          "_id": "ID already exists."
        }
      },
      {
        "status": 400,
        "item": {
          "name": "test-pbehavior-to-check-unique-name"
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
        "id": "test-pbehavior-to-bulk-create-1-2",
        "item": {
          "enabled": true,
          "name": "test-pbehavior-to-bulk-create-1-2-name",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": "test-type-to-pbh-edit-1",
          "reason": "test-reason-1",
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-bulk-create-1-2-filter"
              }
            ]
          },
          "exdates": [
            {
              "begin": 1591164001,
              "end": 1591167601,
              "type": "test-type-to-pbh-edit-1"
            }
          ],
          "exceptions": ["test-exception-to-pbh-edit"]
        }
      }
    ]
    """
    When I do GET /api/v4/pbehaviors?search=test-pbehavior-to-bulk-create-1&sort_by=name
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-pbehavior-to-bulk-create-1-1",
          "enabled": true,
          "author": "root",
          "name": "test-pbehavior-to-bulk-create-1-1-name",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          },
          "reason": {
            "_id": "test-reason-1"
          },
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-bulk-create-1-1-filter"
              }
            ]
          },
          "exdates": [
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
        },
        {
          "enabled": true,
          "author": "root",
          "name": "test-pbehavior-to-bulk-create-1-2-name",
          "tstart": 1591172881,
          "tstop": 1591536400,
          "type": {
            "_id": "test-type-to-pbh-edit-1"
          },
          "reason": {
            "_id": "test-reason-1"
          },
          "filter": {
            "$and": [
              {
                "name": "test-pbehavior-to-bulk-create-1-2-filter"
              }
            ]
          },
          "exdates": [
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
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 2
      }
    }
    """

Feature: Bulk create pbehaviors
  I need to be able to create multiple pbehaviors
  Only admin should be able to create multiple pbehaviors

  Scenario: given bulk create request should return ok
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
    Then the response code should be 201
    Then the response body should contain:
    """json
    [
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

  Scenario: given invalid bulk create request should return errors
    When I am admin
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "0.enabled": "Enabled is missing.",
        "0.name": "Name is missing.",
        "0.filter": "Filter is missing.",
        "0.tstart": "Start is missing.",
        "0.reason": "Reason is missing.",
        "0.type": "Type is missing."
      }
    }
    """

  Scenario: given bulk create request with one invalid and one valid data should return errors
    When I am admin
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "enabled": true,
        "name": "test-pbehavior-to-bulk-create-2-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter": {
          "$and": [
            {
              "name": "test-pbehavior-to-bulk-create-2-filter"
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
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "1.enabled": "Enabled is missing.",
        "1.name": "Name is missing.",
        "1.filter": "Filter is missing.",
        "1.tstart": "Start is missing.",
        "1.reason": "Reason is missing.",
        "1.type": "Type is missing."
      }
    }
    """

  Scenario: given bulk create request with already exists id or name should return error
    When I am admin
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-check-unique"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "0._id": "ID already exists."
      }
    }
    """
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "name": "test-pbehavior-to-check-unique-name"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
          "0.name": "Name already exists."
      }
    }
    """

  Scenario: given bulk create request with multiple items with the same id or name should return error
    When I am admin
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-create-2"
      },
      {
        "_id": "test-pbehavior-to-bulk-create-2"
      },
      {
        "_id": "test-pbehavior-to-bulk-create-2"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "1._id": "ID already exists.",
        "2._id": "ID already exists."
      }
    }
    """
    When I do POST /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "name": "test-pbehavior-to-bulk-create-2"
      },
      {
        "name": "test-pbehavior-to-bulk-create-2"
      },
      {
        "name": "test-pbehavior-to-bulk-create-2"
      }
    ]
    """
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "1.name": "Name already exists.",
        "2.name": "Name already exists."
      }
    }
    """

  Scenario: given bulk create request and no auth user should not allow access
    When I do POST /api/v4/bulk/pbehaviors
    Then the response code should be 401

  Scenario: given bulk create request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/bulk/pbehaviors
    Then the response code should be 403

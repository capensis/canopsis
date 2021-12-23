Feature: Bulk update a pbehaviors
  I need to be able to update multiple pbehaviors
  Only admin should be able to update multiple pbehaviors

  Scenario: given bulk update request should update pbehavior
    When I am admin
    Then I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-update-1",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-1-name-updated",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter":{
          "$and":[
            {
              "name": "test filter"
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
        "_id": "test-pbehavior-to-bulk-update-2",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-2-name-updated",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter":{
          "$and":[
            {
              "name": "test filter"
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
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-update-1",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-1-name-updated",
        "tstart": 1591172881,
        "tstop": 1591536400,
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
        ],
        "author": "root",
        "created": 1592215337
      },
      {
        "_id": "test-pbehavior-to-bulk-update-2",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-2-name-updated",
        "tstart": 1591172881,
        "tstop": 1591536400,
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
        ],
        "author": "root",
        "created": 1592215337
      }
    ]
    """

  Scenario: given bulk update request with not exist ids should return not found error
    When I am admin
    When I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-not-found",
        "enabled": true,
        "name": "test-pbehavior-not-found-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter":{
          "$and":[
            {
              "name": "test filter"
            }
          ]
        }
      },
      {
        "_id": "test-pbehavior-to-bulk-update-2",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-2-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter":{
          "$and":[
            {
              "name": "test filter"
            }
          ]
        }
      }
    ]
    """
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given invalid bulk update request should return errors
    When I am admin
    When I do PUT /api/v4/bulk/pbehaviors:
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
        "0._id": "ID is missing.",
        "0.enabled": "Enabled is missing.",
        "0.name": "Name is missing.",
        "0.filter": "Filter is missing.",
        "0.tstart": "Start is missing.",
        "0.reason": "Reason is missing.",
        "0.type": "Type is missing."
      }
    }
    """

  Scenario: given bulk update request with one valid item and one invalid item should return errors
    When I am admin
    When I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-update-1",
        "enabled": true,
        "name": "test-pbehavior-to-bulk-update-1-name",
        "tstart": 1591172881,
        "tstop": 1591536400,
        "type": "test-type-to-pbh-edit-1",
        "reason": "test-reason-1",
        "filter":{
          "$and":[
            {
              "name": "test filter"
            }
          ]
        }
      },
      {}
    ]
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "1._id": "ID is missing.",
        "1.enabled": "Enabled is missing.",
        "1.name": "Name is missing.",
        "1.filter": "Filter is missing.",
        "1.tstart": "Start is missing.",
        "1.reason": "Reason is missing.",
        "1.type": "Type is missing."
      }
    }
    """

  Scenario: given bulk update request with already exists name should return error
    When I am admin
    When I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-update-1",
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

  Scenario: given bulk update request with multiple items with the same name should return error
    When I am admin
    Then I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "name": "test-pbehavior-to-bulk-update-2-name"
      },
      {
        "name": "test-pbehavior-to-bulk-update-2-name"
      },
      {
        "name": "test-pbehavior-to-bulk-update-2-name"
      }
    ]
    """
    Then the response code should be 400
    """json
    {
      "errors": {
          "1.name": "Name already exists.",
          "2.name": "Name already exists."
      }
    }
    """

  Scenario: given bulk update request with multiple items with the same id should return error
    When I am admin
    Then I do PUT /api/v4/bulk/pbehaviors:
    """json
    [
      {
        "_id": "test-pbehavior-to-bulk-update-1"
      },
      {
        "_id": "test-pbehavior-to-bulk-update-1"
      },
      {
        "_id": "test-pbehavior-to-bulk-update-1"
      }
    ]
    """
    Then the response code should be 400
    """json
    {
      "errors": {
          "1._id": "ID already exists.",
          "2._id": "ID already exists."
      }
    }
    """

  Scenario: given bulk update request and no auth user should not allow access
    When I do PUT /api/v4/bulk/pbehaviors
    Then the response code should be 401

  Scenario: given bulk update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/bulk/pbehaviors
    Then the response code should be 403

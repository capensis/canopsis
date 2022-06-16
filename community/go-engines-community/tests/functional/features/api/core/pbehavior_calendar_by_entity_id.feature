Feature: get a pbehaviors' calendar
  I need to be able to get a pbehaviors' calendar
  Only admin should be able to get a pbehaviors' calendar

  Scenario: given get request should return pbehaviors' intervals
    When I am admin
    When I do GET /api/v4/entities/pbehavior-calendar?_id=test-resource-to-pbh-calendar-get-by-entity/test-component-default&from={{ parseTime "01-10-2022 22:00" }}&to={{ parseTime "02-10-2022 22:00" }}
    Then the response code should be 200
    Then the response body should contain:
    """json
    [
      {
        "_id": "test-pbehavior-to-calendar-get-by-entity-2",
        "title": "test-pbehavior-to-calendar-get-by-entity-2-name",
        "color": "#FFFFFF",
        "from": {{ parseTime "02-10-2022 10:00" }},
        "to": {{ parseTime "02-10-2022 11:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "_id": "test-pbehavior-to-calendar-get-by-entity-1",
        "title": "test-pbehavior-to-calendar-get-by-entity-1-name",
        "color": "#FFFFFF",
        "from": {{ parseTime "02-10-2022 10:00" }},
        "to": {{ parseTime "02-10-2022 12:00" }},
        "type": {
          "_id": "test-default-maintenance-type"
        }
      },
      {
        "_id": "test-pbehavior-to-calendar-get-by-entity-2",
        "title": "test-pbehavior-to-calendar-get-by-entity-2-name",
        "color": "#FFFFFF",
        "from": {{ parseTime "02-10-2022 11:00" }},
        "to": {{ parseTime "02-10-2022 12:00" }},
        "type": {
          "_id": "test-default-active-type"
        }
      },
      {
        "_id": "test-pbehavior-to-calendar-get-by-entity-2",
        "title": "test-pbehavior-to-calendar-get-by-entity-2-name",
        "color": "#FFFFFF",
        "from": {{ parseTime "01-10-2022 22:00" }},
        "to": {{ parseTime "02-10-2022 10:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      },
      {
        "_id": "test-pbehavior-to-calendar-get-by-entity-2",
        "title": "test-pbehavior-to-calendar-get-by-entity-2-name",
        "color": "#FFFFFF",
        "from": {{ parseTime "02-10-2022 12:00" }},
        "to": {{ parseTime "02-10-2022 22:00" }},
        "type": {
          "_id": "test-default-inactive-type"
        }
      }
    ]
    """

  Scenario: given get request with not exist id should return error
    When I am admin
    When I do GET /api/v4/entities/pbehavior-calendar?_id=not-exist&from={{ parseTime "01-10-2022 22:00" }}&to={{ parseTime "06-10-2022 22:00" }}
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  Scenario: given invalid get request should return errors
    When I am admin
    When I do GET /api/v4/entities/pbehavior-calendar
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "_id": "ID is missing.",
        "from": "From is missing.",
        "to": "To is missing."
      }
    }
    """
    When I do GET /api/v4/entities/pbehavior-calendar?from={{ parseTime "03-10-2022 22:00" }}&to={{ parseTime "02-10-2022 22:00" }}
    Then the response code should be 400
    Then the response body should contain:
    """json
    {
      "errors": {
        "to": "To should be greater than From."
      }
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/entities/pbehavior-calendar
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/entities/pbehavior-calendar
    Then the response code should be 403

Feature: Update a widget template
  I need to be able to update a widget template
  Only admin should be able to update a widget template

  Scenario: given update request should return ok
    When I am admin
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-update-1:
    """json
    {
      "title": "test-widgettemplate-to-update-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "extra_details"
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
      "title": "test-widgettemplate-to-update-1-title",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        },
        {
          "value": "v.component"
        },
        {
          "value": "extra_details"
        }
      ],
      "created": 1605263992
    }
    """

  Scenario: given update request and no auth user should not allow access
    When I do PUT /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 401

  Scenario: given update request and auth user without permissions should not allow access
    When I am noperms
    When I do PUT /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 403

  Scenario: given update request with not exist id should return not found error
    When I am admin
    When I do PUT /api/v4/widget-templates/test-widgettemplate-notexist:
    """json
    {
      "title": "test-widgettemplate-to-update-notexist",
      "type": "alarm_columns",
      "columns": [
        {
          "value": "v.resource"
        }
      ]
    }
    """
    Then the response code should be 404

  Scenario: given update request with another type should return bad request error
    When I am admin
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-update-1:
    """json
    {
      "title": "test-widgettemplate-to-update-1-title",
      "type": "entity_columns",
      "columns": [
        {
          "value": "_id"
        }
      ]
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {"type": "Type cannot be changed"}
    }
    """

  Scenario: given update request with missing fields should return bad request error
    When I am admin
    When I do PUT /api/v4/widget-templates/test-widgettemplate-to-update-1:
    """json
    {
    }
    """
    Then the response code should be 400
    Then the response body should be:
    """json
    {
      "errors": {
        "title": "Title is missing.",
        "type": "Type is missing."
      }
    }
    """

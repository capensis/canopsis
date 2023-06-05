Feature: Get a widget template
  I need to be able to get a widget template
  Only admin should be able to get a widget template

  Scenario: given search request should return templates
    When I am admin
    When I do GET /api/v4/widget-templates?search=test-widgettemplate-to-get&sort_by=title
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-widgettemplate-to-get-1",
          "title": "test-widgettemplate-to-get-1-title",
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
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-widgettemplate-to-get-2",
          "title": "test-widgettemplate-to-get-2-title",
          "type": "entity_columns",
          "columns": [
            {
              "value": "_id"
            },
            {
              "value": "type"
            }
          ],
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        },
        {
          "_id": "test-widgettemplate-to-get-3",
          "title": "test-widgettemplate-to-get-3-title",
          "type": "alarm_more_infos",
          "content": "{{ `{{ alarm.v.display_name }}` }}",
          "author": {
            "_id": "root",
            "name": "root"
          },
          "created": 1605263992,
          "updated": 1605263992
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 3
      }
    }
    """

  Scenario: given filter request should return templates by type
    When I am admin
    When I do GET /api/v4/widget-templates?search=test-widgettemplate-to-get&sort_by=title&type=alarm_columns
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-widgettemplate-to-get-1",
          "type": "alarm_columns"
        }
      ],
      "meta": {
        "page": 1,
        "page_count": 1,
        "per_page": 10,
        "total_count": 1
      }
    }
    """

  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/widget-templates
    Then the response code should be 401

  Scenario: given get all request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widget-templates
    Then the response code should be 403

  Scenario: given get request should return template
    When I am admin
    When I do GET /api/v4/widget-templates/test-widgettemplate-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-widgettemplate-to-get-1",
      "title": "test-widgettemplate-to-get-1-title",
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
      "author": {
        "_id": "root",
        "name": "root"
      },
      "created": 1605263992,
      "updated": 1605263992
    }
    """

  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 401

  Scenario: given get request and auth user without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 403

  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/widget-templates/test-widgettemplate-notexist
    Then the response code should be 404

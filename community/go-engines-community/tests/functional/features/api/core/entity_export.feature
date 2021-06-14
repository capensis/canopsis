Feature: Export entities
  I need to be able to export a entities

  Scenario: given export request should return entities
    When I am admin
    When I do POST /api/v4/entity-export?search=test-entity-to-export&active_columns[]=_id&active_columns[]=name&active_columns[]=type&active_columns[]=infos
    Then the response code should be 200
    When I do GET /api/v4/entity-export/{{ .lastResponse._id }}
    Then the response code should be 200
    When I wait 1s
    When I do GET /api/v4/entity-export/{{ .lastResponse._id }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """
    _id,name,type,infos
    test-entity-to-export-resource-1/test-entity-to-export-component,test-entity-to-export-resource-1,resource,"{""test-entity-to-export-resource-1-info-1"":{""name"":""test-entity-to-export-resource-1-info-1-name"",""description"":""test-entity-to-export-resource-1-info-1-description"",""value"":""test-entity-to-export-resource-1-info-1-value""}}"
    test-entity-to-export-resource-2/test-entity-to-export-component,test-entity-to-export-resource-2,resource,{}

    """

  Scenario: given export request should return empty response
    When I am admin
    When I do POST /api/v4/entity-export?search=not-found&active_columns[]=_id&active_columns[]=name&active_columns[]=type&active_columns[]=infos
    Then the response code should be 200
    When I do GET /api/v4/entity-export/{{ .lastResponse._id }}
    Then the response code should be 200
    When I wait 1s
    When I do GET /api/v4/entity-export/{{ .lastResponse._id }}/download
    Then the response code should be 200
    Then the response raw body should be:
    """
    _id,name,type,infos

    """


  Scenario: given not exit export task should return not found error
    When I am admin
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 404
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 404

  Scenario: given unauth request should not allow access
    When I do POST /api/v4/entity-export
    Then the response code should be 401
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 401
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 401

  Scenario: given request and auth user without permissions should not allow access
    When I am noperms
    When I do POST /api/v4/entity-export
    Then the response code should be 403
    When I do GET /api/v4/entity-export/not-found
    Then the response code should be 403
    When I do GET /api/v4/entity-export/not-found/download
    Then the response code should be 403
Feature: Get a view group
  I need to be able to get a view group
  Only admin should be able to get a view group

  @concurrent
  Scenario: given search request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "created": 1611229670,
          "updated": 1611229670
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

  @concurrent
  Scenario: given get request should return viewgroup
    When I am admin
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "_id": "test-viewgroup-to-get-1",
      "title": "test-viewgroup-to-get-1-title",
      "is_private": false,
      "author": {
        "_id": "root",
        "name": "root",
        "display_name": "root John Doe admin@canopsis.net"
      },
      "created": 1611229670,
      "updated": 1611229670
    }
    """

  @concurrent
  Scenario: given with views request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true&with_tabs=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "tabs": [],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "tabs": [
                {
                  "_id": "test-tab-to-viewgroup-get-1",
                  "title": "test-tab-to-viewgroup-get-1-title",
                  "is_private": false,
                  "author": {
                    "_id": "root",
                    "name": "root",
                    "display_name": "root John Doe admin@canopsis.net"
                  },
                  "created": 1611229670,
                  "updated": 1611229670
                },
                {
                  "_id": "test-tab-to-viewgroup-get-2",
                  "title": "test-tab-to-viewgroup-get-2-title",
                  "is_private": false,
                  "author": {
                    "_id": "root",
                    "name": "root",
                    "display_name": "root John Doe admin@canopsis.net"
                  },
                  "created": 1611229670,
                  "updated": 1611229670
                }
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true&with_tabs=true&with_widgets=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "is_private": false,
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "tabs": [],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "tabs": [
                {
                  "_id": "test-tab-to-viewgroup-get-1",
                  "title": "test-tab-to-viewgroup-get-1-title",
                  "author": {
                    "_id": "root",
                    "name": "root",
                    "display_name": "root John Doe admin@canopsis.net"
                  },
                  "is_private": false,
                  "widgets": [
                    {
                      "_id": "test-widget-to-viewgroup-get-1",
                      "author": {
                        "_id": "root",
                        "name": "root",
                        "display_name": "root John Doe admin@canopsis.net"
                      },
                      "is_private": false,
                      "created": 1611229670,
                      "updated": 1611229670,
                      "grid_parameters": {
                        "desktop": {"x": 0, "y": 0}
                      },
                      "parameters": {
                        "test-widget-to-viewgroup-get-1-parameter-1": {
                          "test-widget-to-viewgroup-get-1-parameter-1-subparameter": "test-widget-to-viewgroup-get-1-parameter-1-subvalue"
                        },
                        "test-widget-to-viewgroup-get-1-parameter-2": [
                          {
                            "test-widget-to-viewgroup-get-1-parameter-2-subparameter": "test-widget-to-viewgroup-get-1-parameter-2-subvalue"
                          }
                        ]
                      },
                      "filters": [
                        {
                          "_id": "test-widgetfilter-to-viewgroup-get-1",
                          "title": "test-widgetfilter-to-viewgroup-get-1-title",
                          "is_user_preference": false,
                          "is_private": false,
                          "author": {
                            "_id": "nopermsuser",
                            "name": "nopermsuser",
                            "display_name": "nopermsuser   "
                          },
                          "created": 1611229670,
                          "updated": 1611229670,
                          "alarm_pattern": [
                            [
                              {
                                "field": "v.component",
                                "cond": {
                                  "type": "eq",
                                  "value": "test-widgetfilter-to-viewgroup-get-1-pattern"
                                }
                              }
                            ]
                          ]
                        },
                        {
                          "_id": "test-widgetfilter-to-viewgroup-get-2",
                          "title": "test-widgetfilter-to-viewgroup-get-2-title",
                          "is_user_preference": false,
                          "is_private": false,
                          "author": {
                            "_id": "root",
                            "name": "root",
                            "display_name": "root John Doe admin@canopsis.net"
                          },
                          "created": 1611229670,
                          "updated": 1611229670,
                          "alarm_pattern": [
                            [
                              {
                                "field": "v.component",
                                "cond": {
                                  "type": "eq",
                                  "value": "test-widgetfilter-to-viewgroup-get-2-pattern"
                                }
                              }
                            ]
                          ]
                        }
                      ],
                      "title": "test-widget-to-viewgroup-get-1-title",
                      "type": "test-widget-to-viewgroup-get-1-type"
                    },
                    {
                      "_id": "test-widget-to-viewgroup-get-2",
                      "author": {
                        "_id": "root",
                        "name": "root",
                        "display_name": "root John Doe admin@canopsis.net"
                      },
                      "is_private": false,
                      "created": 1611229670,
                      "updated": 1611229670,
                      "grid_parameters": {
                        "desktop": {"x": 0, "y": 1}
                      },
                      "filters": [],
                      "parameters": {},
                      "title": "test-widget-to-viewgroup-get-2-title",
                      "type": "test-widget-to-viewgroup-get-2-type"
                    }
                  ],
                  "created": 1611229670,
                  "updated": 1611229670
                },
                {
                  "_id": "test-tab-to-viewgroup-get-2",
                  "title": "test-tab-to-viewgroup-get-2-title",
                  "author": {
                    "_id": "root",
                    "name": "root",
                    "display_name": "root John Doe admin@canopsis.net"
                  },
                  "is_private": false,
                  "widgets": [],
                  "created": 1611229670,
                  "updated": 1611229670
                }
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "is_private": false,
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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

  @concurrent
  Scenario: given with flags request should return view groups
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_flags=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "is_private": false,
          "deletable": false
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "is_private": false,
          "deletable": true
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
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-get&with_views=true&with_flags=true
    Then the response code should be 200
    Then the response body should be:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-get-1",
          "title": "test-viewgroup-to-get-1-title",
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "is_private": false,
          "deletable": false,
          "views": [
            {
              "_id": "test-view-to-viewgroup-get-2",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-2-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-2-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-2-tag"
              ],
              "updated": 1611229670
            },
            {
              "_id": "test-view-to-viewgroup-get-1",
              "author": {
                "_id": "root",
                "name": "root",
                "display_name": "root John Doe admin@canopsis.net"
              },
              "created": 1611229670,
              "description": "test-view-to-viewgroup-get-1-description",
              "enabled": true,
              "is_private": false,
              "group": {
                "_id": "test-viewgroup-to-get-1",
                "author": {
                  "_id": "root",
                  "name": "root",
                  "display_name": "root John Doe admin@canopsis.net"
                },
                "created": 1611229670,
                "title": "test-viewgroup-to-get-1-title",
                "updated": 1611229670
              },
              "title": "test-view-to-viewgroup-get-1-title",
              "periodic_refresh": {
                "enabled": true,
                "value": 1,
                "unit": "s"
              },
              "tags": [
                "test-view-to-viewgroup-get-1-tag"
              ],
              "updated": 1611229670
            }
          ],
          "created": 1611229670,
          "updated": 1611229670
        },
        {
          "_id": "test-viewgroup-to-get-2",
          "title": "test-viewgroup-to-get-2-title",
          "author": {
            "_id": "root",
            "name": "root",
            "display_name": "root John Doe admin@canopsis.net"
          },
          "is_private": false,
          "deletable": true,
          "views": [],
          "created": 1611229670,
          "updated": 1611229670
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

  @concurrent
  Scenario: given get all request and no auth user should not allow access
    When I do GET /api/v4/view-groups
    Then the response code should be 401

  @concurrent
  Scenario: given get all request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-groups
    Then the response code should be 403

  @concurrent
  Scenario: given get request and no auth user should not allow access
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 401

  @concurrent
  Scenario: given get request and auth user by api key without permissions should not allow access
    When I am noperms
    When I do GET /api/v4/view-groups/test-viewgroup-to-get-1
    Then the response code should be 403

  @concurrent
  Scenario: given get request with not exist id should return not found error
    When I am admin
    When I do GET /api/v4/view-groups/test-viewgroup-not-found
    Then the response code should be 404
    Then the response body should be:
    """json
    {
      "error": "Not found"
    }
    """

  @concurrent
  Scenario: given get request shouldn't return private views without with_private flag
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-private-get&with_views=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-private-get-2",
          "title": "test-viewgroup-to-private-get-2-title",
          "is_private": false,
          "views": [
            {
              "_id": "test-view-viewgroup-to-private-get-2",
              "is_private": false
            }
          ]
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

  @concurrent
  Scenario: given get request should return private views with with_private flag only for owner
    When I am admin
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-private-get&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-private-get-2",
          "title": "test-viewgroup-to-private-get-2-title",
          "is_private": false,
          "views": [
            {
              "_id": "test-view-viewgroup-to-private-get-2",
              "is_private": false
            }
          ]
        },
        {
          "_id": "test-viewgroup-to-private-get-1",
          "title": "test-viewgroup-to-private-get-1-title",
          "is_private": true,
          "views": [
            {
              "_id": "test-view-viewgroup-to-private-get-1",
              "is_private": true
            }
          ]
        },
        {
          "_id": "test-viewgroup-to-private-get-4",
          "title": "test-viewgroup-to-private-get-4-title",
          "is_private": true,
          "views": [
            {
              "_id": "test-view-viewgroup-to-private-get-4",
              "is_private": true
            }
          ]
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
    When I am manager
    When I do GET /api/v4/view-groups?search=test-viewgroup-to-private-get&with_views=true&with_private=true
    Then the response code should be 200
    Then the response body should contain:
    """json
    {
      "data": [
        {
          "_id": "test-viewgroup-to-private-get-2",
          "title": "test-viewgroup-to-private-get-2-title",
          "is_private": false,
          "views": [
            {
              "_id": "test-view-viewgroup-to-private-get-2",
              "is_private": false
            }
          ]
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

# Associative Table APIs

## CRUD on associative table system

This endpoint create/update/delete/returns a configuration keeped in the
database (associative table).

#### Url

  `GET` /api/v2/associativetable/*name*

  `PUT` /api/v2/associativetable/*name*

  `POST` /api/v2/associativetable/*name*

  `DELETE` /api/v2/associativetable/*name*

#### Params

- **name** : the name that identify a configuration document

#### GET example

/api/v2/associativetable/pe_links_builder

Response:
```{json}
    {
        "basic_link_builder": {
            "base_url": "http://example.com/screenshot?file={name}",
            "category": "screenshot"
        }
    }
```

#### POST example

/api/v2/associativetable/pe_links_builder

Payload json:
```{json}
    {
        "basic_link_builder" : {
            "base_url": "http://example.com/screenshot?file={name}",
            "category": "screenshot"
        }
    }
```

Response:
```{json}
true
```

#### PUT example

/api/v2/associativetable/pe_links_builder

Payload json:
```{json}
    {
        "basic_link_builder" : {
            "base_url": "http://example.com/capture?file={name}",
            "category": "capture"
        }
    }
```

Response:
```{json}
true
```

#### DELETE example

/api/v2/associativetable/pe_links_builder

Response:
```{json}
true
```

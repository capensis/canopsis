## Migration de @route vers bottle pur

Exemples de code actuel / ancien :

```python
@route(ws.application.post, name='path', payload=['param1', 'param2']):
def route_func(param1, param2):
    ...
```

Nouveau code :

```python
from bottle import request

from canopsis.webcore.utils import gen_json_error, gen_json, HTTP_ERROR

@ws.application.post('/path')
def route_func():
    """
    :param param1 <type>: <explain>
    :param param2 <type>: <explain>
    :param get_param1 <type>: <explain>
    :param get_param2 <type>: <explain>
    """
    body = json.loads(request.body.read())

    try:
        param1 = <type>(body.get['param1'])
        get_param1 = <type>(request.params['get_param1'])
    except KeyError as exc:
        return gen_json_error(
            'missing required param: {}'.format(exc), HTTP_ERROR
        )

    param2 = <type>(body.get('param2', 'default_value'))
    get_param2 = <type>(request.params.get('get_param2', 'default_value'))


    try:
        res = do_some_work(param1, param2)
    except (SomeExceptions,) as exc:
        return gen_json_error('error: {}'.format(exc), HTTP_ERROR)

    return gen_json(res)
```
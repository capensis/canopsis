from canopsis.context.manager import Context


def transform_id_to_name():

    context = Context()

    storage = context['ctx_storage']

    elements = storage.get_elements()

    for element in elements:
        if 'id' in element:
            element[Context.NAME] = element.pop('id')
            _id = element['_id']
            storage.put_element(_id=_id, element=element)

if __name__ == '__main__':
    transform_id_to_name()

def cfg_to_array(value, sep=','):
    """
    Convert a configuration line to an array.

    :param str value: a value to parse
    :param str sep: separator
    :rtype: array
    """
    values = value.split(',')

    f_values = []

    for val in values:
        val.strip()
        if len(val) > 0:
            f_values.append(val)

    return f_values

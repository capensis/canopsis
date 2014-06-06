__all__ = (
	'get_aggregation_value',
	'add_aggregation',
	'get_aggregations',
	'AggregationError')

_AGGREGATIONS = dict()


class AggregationError(Exception):
	pass


def get_aggregations():
	"""
	Get aggregation functions by name.
	"""

	result = _AGGREGATIONS.copy()

	return result


def get_aggregation_value(name, points):
	"""
	Get aggregation value where with related name and points.

	Points must have only real values.
	"""

	aggregation = _AGGREGATIONS.get(name, None)
	if aggregation is None:
		raise NotImplementedError("No aggregation {0} exists".format(name))

	result = aggregation(points)

	return result


def add_aggregation(name, function, push=False):
	"""
	Set an aggregation function to this AGGREGATIONS module variable.

	- push :
		if False, raise an AggregationError if an aggregation has already been
	added with the same name.
	- push : change of aggregation if name already exists.

	Added aggregations are available through module properties.
	"""

	if push is False and name in _AGGREGATIONS:
		raise AggregationError("name {0} already exists".format(name))

	_AGGREGATIONS[name] = function


def _mean(points):
	return sum(points) / len(points)
add_aggregation('MEAN', _mean)


def _last(points):
	return points[-1]
add_aggregation('LAST', _last)


def _first(points):
	return points[0]
add_aggregation('FIRST', _first)


def _delta(points):
	return (max(points) - min(points)) / 2
add_aggregation('DELTA', _delta)


def _sum(points):
	return sum(points)
add_aggregation('SUM', _sum)


def _max(points):
	return max(points)
add_aggregation('MAX', _max)


def _min(points):
	return min(points)
add_aggregation('MIN', _min)

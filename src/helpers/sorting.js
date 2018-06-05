import get from 'lodash/get';

export function stringSortObject(arrayToSort, propertyObject) {
  arrayToSort.sort((a, b) => get(a, propertyObject).localeCompare(get(b, propertyObject)));
}

export function numericSortObject(arrayToSort, propertyObject, direction = 'asc') {
  let compareFunct;
  if (direction === 'asc') {
    compareFunct = (a, b) => get(a, propertyObject) - (get(b, propertyObject));
  } else if (direction === 'desc') {
    compareFunct = (b, a) => get(a, propertyObject) - (get(b, propertyObject));
  }
  arrayToSort.sort(compareFunct);
}

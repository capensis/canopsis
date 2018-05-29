import get from 'lodash/get';

export function stringSortObject(arrayToSort, propertyObject) {
  arrayToSort.sort((a, b) => get(a, propertyObject).localeCompare(get(b, propertyObject)));
}

export function numericSortObject(arrayToSort, propertyObject) {
  arrayToSort.sort((a, b) => get(a, propertyObject) - (get(b, propertyObject)));
}

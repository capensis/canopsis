import { get } from 'lodash';

/**
 * Sort an array of objects by a property whith string value
 * @param arrayToSort The array to sort
 * @param propertyObject The property to compare
 * @returns {*[]} An new sorted array
 */
export function stringSortObject(arrayToSort, propertyObject) {
  return [...arrayToSort].sort((a, b) => get(a, propertyObject).localeCompare(get(b, propertyObject)));
}

/**
 * Sort an array of objects by a property whith numeric value
 * @param arrayToSort The array to sort
 * @param propertyObject The property to compare
 * @param direction The direction of sorting, either 'asc' or 'desc'
 * @returns {*[]} An new sorted array
 */
export function numericSortObject(arrayToSort, propertyObject, direction = 'asc') {
  let compareFunct;
  if (direction === 'asc') {
    compareFunct = (a, b) => get(a, propertyObject) - (get(b, propertyObject));
  } else if (direction === 'desc') {
    compareFunct = (b, a) => get(a, propertyObject) - (get(b, propertyObject));
  }
  return [...arrayToSort].sort(compareFunct);
}

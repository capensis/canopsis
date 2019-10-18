import { get, isNull, isNaN } from 'lodash';

/**
 * Custom sort method for v-data-table vuetify component with ignoring null values.
 * If sorting is enabled items with null value will be always on the end of table.
 *
 * @param items - items of data table component
 * @param index - index for sorting
 * @param isDescending
 * @returns {*}
 */
export function dataTableCustomSortWithNullIgnoring(items, index, isDescending) {
  if (isNull(index)) {
    return items;
  }

  return items.sort((a, b) => {
    let sortA = get(a, index);
    let sortB = get(b, index);

    if (isNull(sortA)) {
      if (isNull(sortB)) {
        return 0;
      }

      return 1;
    }

    if (isNull(sortB)) {
      return -1;
    }

    if (isDescending) {
      [sortA, sortB] = [sortB, sortA];
    }

    if (!isNaN(sortA) && !isNaN(sortB)) {
      return sortA - sortB;
    }

    [sortA, sortB] = [sortA, sortB].map(s => (s || '').toString().toLocaleLowerCase());

    if (sortA > sortB) {
      return 1;
    }

    if (sortA < sortB) {
      return -1;
    }

    return 0;
  });
}

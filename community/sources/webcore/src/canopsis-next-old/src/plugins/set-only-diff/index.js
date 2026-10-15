import { isObject, isArray } from 'lodash';
import { detailedDiff } from 'deep-object-diff';

/**
 * Recursively deletes properties from an object using Vue.delete.
 *
 * @param {Vue} Vue - The Vue instance.
 * @param {Object} source - The source object to delete properties from.
 * @param {Object} deleted - The object containing properties to delete.
 */
const deepDelete = (Vue, source, deleted) => {
  Object.entries(deleted).forEach(([key, value]) => {
    if (!value) {
      Vue.delete(source, key);

      return;
    }

    deepDelete(Vue, source[key], value);
  });
};

/**
 * Recursively updates properties in an object using Vue.set.
 *
 * @param {Vue} Vue - The Vue instance.
 * @param {Object} source - The source object to update properties in.
 * @param {Object} updated - The object containing properties to update.
 */
const deepUpdate = (Vue, source, updated) => {
  Object.entries(updated).forEach(([key, value]) => {
    if (source[key] && (isObject(value) || isArray(value))) {
      deepUpdate(Vue, source[key], value);

      return;
    }

    Vue.set(source, key, value);
  });
};

/**
 * TODO: use it in the future
 */
export default {
  install(Vue) {
    Vue.setOnlyDiff = (source, key, payload) => {
      if (!source[key]) {
        Vue.set(source, key, payload);

        return;
      }

      const { added, updated, deleted } = detailedDiff(source[key], payload);

      deepUpdate(Vue, source[key], added);
      deepUpdate(Vue, source[key], updated);
      deepDelete(Vue, source[key], deleted);
    };
  },
};

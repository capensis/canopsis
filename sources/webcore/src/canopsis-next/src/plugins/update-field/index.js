import { get } from 'lodash';

import { setField } from '@/helpers/immutable';

export default {
  install(Vue) {
    Vue.prototype.$updateField = function updateField(path, value, mutate = false) {
      const { prop = 'value', event = 'input' } = this.$options.model || {};

      if (mutate) {
        return this.$emit(event, value, path);
      }

      return this.$emit(event, path.length ? setField(this[prop], path, value) : value);
    };

    Vue.prototype.$updateFieldModel = function updateField(source, path = [], value = undefined) {
      if (!path.length) {
        throw new Error(`Incorrect path: ${path}`);
      }

      const copiedPath = [...path];
      const fieldKey = copiedPath.pop();

      this.$set(copiedPath.length ? get(source, copiedPath) : source, fieldKey, value);
    };
  },
};

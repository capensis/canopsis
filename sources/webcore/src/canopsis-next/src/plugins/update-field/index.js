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
  },
};

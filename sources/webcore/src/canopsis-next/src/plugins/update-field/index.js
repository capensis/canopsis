import { setField } from '@/helpers/immutable';

export default {
  install(Vue) {
    Vue.prototype.$updateField = function updateField(path, value) {
      const { prop = 'value', event = 'input' } = this.$options.model || {};

      this.$emit(event, path.length ? setField(this[prop], path, value) : value);
    };
  },
};

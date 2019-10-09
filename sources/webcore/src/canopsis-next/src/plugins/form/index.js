import { setIn } from '@/helpers/immutable';

export default {
  install(Vue) {
    Vue.prototype.$updateField = function updateField(path, value) {
      const { event = 'input', prop = 'value' } = this.$options.model || {};

      this.$emit(event, setIn(this[prop], path, value));
    };
  },
};

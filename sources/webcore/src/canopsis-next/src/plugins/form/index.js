import { setIn } from '@/helpers/immutable';

export default {
  install(Vue) {
    Object.defineProperty(Vue.prototype, '$form', {
      get() {
        return {
          updateField: (path, value) => {
            const { event = 'input', prop = 'value' } = this.$options.model || {};

            this.$emit(event, setIn(this[prop], path, value));
          },

          updateModel: (value) => {
            const { event = 'input' } = this.$options.model || {};

            this.$emit(event, value);
          },
        };
      },
    });
  },
};

import { setIn, unsetIn } from '@/helpers/immutable';

export default {
  install(Vue) {
    Object.defineProperty(Vue.prototype, '$form', {
      get() {
        return {
          updateModel: (value) => {
            const { event = 'input' } = this.$options.model || {};

            this.$emit(event, value);
          },

          updateField: (path, value) => {
            const { prop = 'value' } = this.$options.model || {};

            this.$form.updateModel(path.length ? setIn(this[prop], path, value) : value);
          },

          removeField: (path) => {
            const { prop = 'value' } = this.$options.model || {};

            this.$form.updateModel(unsetIn(this[prop], path));
          },
        };
      },
    });
  },
};

import { get } from 'lodash';

import { setIn, unsetIn } from '@/helpers/immutable';

export default {
  install(Vue) {
    Object.defineProperty(Vue.prototype, '$form', {
      get() {
        return {
          uf: (path, value, basePath) => {
            const newPath = basePath ? path.concat(basePath) : path;

            this.$form.updateModel(value, newPath);
          },

          um: (path, value, basePath) => {
            const newPath = basePath ? path.concat(basePath) : path;
            const key = newPath.pop();
            const source = get(this, newPath);

            this.$set(source, key, value);
          },

          updateModel: (...args) => {
            const { event = 'input' } = this.$options.model || {};

            this.$emit(event, ...args);
          },

          updateField: (path, value) => {
            const { prop = 'value' } = this.$options.model || {};

            this.$form.updateModel(setIn(this[prop], path, value));
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

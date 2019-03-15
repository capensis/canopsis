import uid from '@/helpers/uid';

export const modelPropKeyComputed = uid('_modelProp');
export const modelEventKeyComputed = uid('_modelEventKey');

/**
 * @mixin Form mixin
 */
export default {
  computed: {
    [modelPropKeyComputed]() {
      if (this.$options.model && this.$options.model.prop) {
        return this.$options.model.prop;
      }

      return 'value';
    },

    [modelEventKeyComputed]() {
      if (this.$options.model && this.$options.model.event) {
        return this.$options.model.event;
      }

      return 'input';
    },
  },
};

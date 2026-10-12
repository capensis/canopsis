import { uid } from '@/helpers/uid';

export const modelPropKeyComputed = uid('_modelProp');
export const modelEventKeyComputed = uid('_modelEventKey');

/**
 * @mixin Form mixin
 */
export default {
  computed: {
    [modelPropKeyComputed]() {
      const event = this.$options.emits?.find(emit => emit.startsWith('update:'));

      if (event) {
        return event.replace('update:', '');
      }

      return 'value';
    },

    [modelEventKeyComputed]() {
      const event = this.$options.emits?.find(emit => emit.startsWith('update:'));

      if (event) {
        return event;
      }

      return 'input';
    },
  },
};

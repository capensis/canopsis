import uid from '../../../helpers/uid';

export const eventKeyComputed = uid('_eventKey');
export const formKeyComputed = uid('_formKey');

/**
 * @mixin Form mixin
 */
export default {
  computed: {
    [formKeyComputed]() {
      if (this.$options.model && this.$options.model.prop) {
        return this.$options.model.prop;
      }

      return 'value';
    },

    [eventKeyComputed]() {
      if (this.$options.model && this.$options.model.event) {
        return this.$options.model.event;
      }

      return 'input';
    },
  },
};

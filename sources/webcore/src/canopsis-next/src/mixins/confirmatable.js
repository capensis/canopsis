import { cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

export default ({ field = 'form', method = 'submit' } = {}) => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');
  const clickOutsideHandlerMethodKey = uid('click-outside');

  return {
    inject: ['$clickOutside'],
    created() {
      this.$clickOutside.register(this[clickOutsideHandlerMethodKey]);
    },
    mounted() {
      this[originalField] = cloneDeep(this[field]);
      this[confirmationModalIdField] = uid('modal');
    },
    beforeDestroy() {
      this.$clickOutside.unregister(this[clickOutsideHandlerMethodKey]);
    },
    methods: {
      [clickOutsideHandlerMethodKey]() {
        const equal = isEqual(this[field], this[originalField]);
        const getterKey = `${this.$modals.moduleName}/hasModalById`;

        if (!equal && !this.$store.getters[getterKey](this[confirmationModalIdField])) {
          this.$modals.show({
            id: this[confirmationModalIdField],
            name: MODALS.formConfirmation,
            config: {
              action: async (submitted) => {
                if (submitted) {
                  return this[method]();
                }

                return this.$modals.hide();
              },
            },
          });
        }

        return equal;
      },
    },
  };
};

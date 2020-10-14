import { get, cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

export default ({
  field = 'form',
  method = 'submit',
  modalName = MODALS.clickOutsideConfirmation,
  comparator = isEqual,
} = {}) => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');
  const clickOutsideHandlerMethodKey = uid('click-outside');

  return {
    inject: ['$clickOutside'],
    created() {
      this.$clickOutside.register(this[clickOutsideHandlerMethodKey]);

      this[originalField] = cloneDeep(this[field]);
      this[confirmationModalIdField] = uid('modal');
    },
    beforeDestroy() {
      this.$clickOutside.unregister(this[clickOutsideHandlerMethodKey]);
    },
    methods: {
      [clickOutsideHandlerMethodKey]() {
        const equal = comparator.call(this, this[field], this[originalField]);
        const statePath = [this.$modals.moduleName, 'byId', this[confirmationModalIdField]];

        if (!equal && !get(this.$store.state, statePath)) {
          this.$modals.show({
            id: this[confirmationModalIdField],
            name: modalName,
            dialogProps: {
              persistent: true,
            },
            config: {
              action: async (confirmed) => {
                if (confirmed) {
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

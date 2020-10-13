import { cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

export default ({
  field = 'form',
  method = 'submit',
  modalName = MODALS.formConfirmation,
  comparator = isEqual,
} = {}) => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');
  const confirmationMethodKey = uid('confirmation');

  return {
    created() {
      this[originalField] = cloneDeep(this[field]);
      const sourceMethod = this[method];

      if (sourceMethod) {
        this[method] = async (...args) => {
          this[confirmationMethodKey](() => sourceMethod.apply(this, args));
        };
      }
    },
    methods: {
      [confirmationMethodKey](action) {
        const equal = comparator.call(this, this[field], this[originalField]);

        if (!equal) {
          this.$modals.show({
            id: this[confirmationModalIdField],
            name: modalName,
            dialogProps: {
              persistent: true,
            },
            config: {
              action: async () => {
                action();
              },
            },
          });
        } else {
          action();
        }
      },
    },
  };
};

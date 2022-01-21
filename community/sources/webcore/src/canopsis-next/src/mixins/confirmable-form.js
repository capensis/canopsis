import { cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

/**
 * Confirm user action on form mixin creator
 *
 * @param {string} [field = 'form']
 * @param {string} [method = 'submit']
 * @param {string} [modalName = MODALS.clickOutsideConfirmation]
 * @param {Function} [comparator = isEqual]
 * @param {boolean} [cloning = false]
 * @returns {{created(): void, methods: {}, beforeDestroy(): void, inject: [string]}|*}
 */
export const confirmableFormMixinCreator = ({
  field = 'form',
  method = 'submit',
  modalName = MODALS.confirmation,
  comparator = isEqual,
  cloning = false,
} = {}) => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');
  const confirmationMethodKey = uid('confirmation');

  return {
    created() {
      if (cloning) {
        this[originalField] = cloneDeep(this[field]);
      }

      const sourceMethod = this[method];

      if (sourceMethod) {
        this[method] = (...args) => this[confirmationMethodKey](() => sourceMethod.apply(this, args));
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
              action,
            },
          });
        } else {
          action();
        }
      },
    },
  };
};

export default confirmableFormMixinCreator;

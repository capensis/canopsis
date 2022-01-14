import { get, cloneDeep, isEqual, isFunction } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

/**
 * Confirm modal click outside mixin creator
 *
 * @param {string} [field = 'form']
 * @param {string} [method = 'submit']
 * @param {string} [closeMethod = '$modals.hide']
 * @param {string} [modalName = MODALS.clickOutsideConfirmation]
 * @param {Function} [comparator = isEqual]
 * @returns {{created(): void, methods: {}, beforeDestroy(): void, inject: [string]}|*}
 */
export const confirmableModalMixinCreator = ({
  field = 'form',
  method = 'submit',
  closeMethod = '$modals.hide',
  modalName = MODALS.clickOutsideConfirmation,
  comparator = isEqual,
} = {}) => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');
  const clickOutsideHandlerMethodKey = uid('click-outside');
  const clickOutsideCloseMethodKey = uid('close-method');

  return {
    provide() {
      return {
        $closeModal: () => {
          if (this[clickOutsideHandlerMethodKey]()) {
            this[clickOutsideCloseMethodKey]();
          }
        },
      };
    },
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
      [clickOutsideCloseMethodKey]() {
        const close = get(this, closeMethod);

        if (isFunction(close)) {
          close();
        }
      },

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
              action: (confirmed) => {
                if (confirmed) {
                  return this[method]();
                }

                return this[clickOutsideCloseMethodKey]();
              },
            },
          });
        }

        return equal;
      },
    },
  };
};

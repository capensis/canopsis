import { cloneDeep, isEqual } from 'lodash';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

export default (field = 'form') => {
  const originalField = Symbol('originalField');
  const confirmationModalIdField = Symbol('confirmationModalIdField');

  return {
    inject: ['$clickOutside'],
    created() {
      this.$clickOutside.register(this.clickOutsideHandler);
    },
    mounted() {
      this[originalField] = cloneDeep(this[field]);
      this[confirmationModalIdField] = uid('modal');
    },
    beforeDestroy() {
      this.$clickOutside.unregister(this.clickOutsideHandler);
    },
    methods: {
      clickOutsideHandler() {
        const equal = isEqual(this[field], this[originalField]);
        const getterKey = `${this.$modals.moduleName}/hasModalById`;

        if (!equal && !this.$store.getters[getterKey](this[confirmationModalIdField])) {
          this.$modals.show({
            id: this[confirmationModalIdField],
            name: MODALS.confirmation,
            config: {
              text: 'Changes will not be saved. Are you sure?',
              action: () => this.$modals.hide(),
            },
          });
        }

        return equal;
      },
    },
  };
};

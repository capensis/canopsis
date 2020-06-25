<template lang="pug">
  td
    v-checkbox(
      v-bind="checkbox.bind",
      v-on="checkbox.on",
      :disabled="disabled",
      color="primary",
      hideDetails
    )
</template>

<script>
import { USERS_RIGHTS_MASKS } from '@/constants';

import { getRightMasks, getCheckboxValue } from '@/helpers/right';

export default {
  model: {
    prop: 'group',
    event: 'change',
  },
  props: {
    group: {
      type: Object,
      required: true,
    },
    role: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    changedRole: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    checkbox() {
      const bind = this.group.rights.reduce((acc, right, index) => {
        const masks = getRightMasks(right);

        masks.forEach((mask, maskIndex) => {
          const checkboxValue = getCheckboxValue(right, this.role, this.changedRole, mask);

          if (!acc.indeterminate && checkboxValue !== acc.inputValue && !(index === 0 && maskIndex === 0)) {
            acc.indeterminate = true;
          }

          acc.inputValue = acc.inputValue || checkboxValue;
        });

        return acc;
      }, {
        inputValue: false,
        indeterminate: false,
      });

      const on = {
        change: value => this.group.rights.forEach((right) => {
          const masks = getRightMasks(right);

          masks.forEach(mask =>
            getCheckboxValue(right, this.role, this.changedRole, mask) !== value
            && this.changeCheckboxValue(value, right, mask));
        }),
      };

      return {
        bind,
        on,
      };
    },
  },
  methods: {
    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} right
     * @param {number} [mask = USERS_RIGHTS_MASKS.default]
     */
    changeCheckboxValue(value, right, mask = USERS_RIGHTS_MASKS.default) {
      this.$emit('change', value, this.role, right, mask);
    },
  },
};
</script>

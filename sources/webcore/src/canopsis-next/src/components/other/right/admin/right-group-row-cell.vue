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
import { get, isUndefined } from 'lodash';
import { USERS_RIGHTS_MASKS } from '@/constants';

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
        const checkboxValue = this.getCheckboxValue(right);
        if (!acc.indeterminate && checkboxValue !== acc.inputValue && index !== 0) {
          acc.indeterminate = true;
        }

        acc.inputValue = acc.inputValue || checkboxValue;

        return acc;
      }, { inputValue: false, indeterminate: false });

      const on = {
        change: value => this.group.rights.forEach(right => this.changeCheckboxValue(value, right)),
      };

      return {
        bind,
        on,
      };
    },
  },
  methods: {
    getCheckboxValue(right, rightMask = USERS_RIGHTS_MASKS.default) {
      const { role, changedRole } = this;

      const checkSum = get(role, ['rights', right._id, 'checksum'], 0);

      const changedCheckSum = get(changedRole, [right._id]);
      const currentCheckSum = isUndefined(changedCheckSum) ? checkSum : changedCheckSum;
      const rightType = currentCheckSum & rightMask;

      return rightType === rightMask;
    },

    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {Object} right
     * @param {number} mask
     */
    changeCheckboxValue(value, right, mask) {
      this.$emit('change', value, this.role, right, mask);
    },
  },
};
</script>

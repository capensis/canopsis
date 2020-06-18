<template lang="pug">
  td
    v-checkbox-functional(
      v-for="(checkbox, index) in checkboxes",
      v-bind="checkbox.bind",
      v-on="checkbox.on",
      :key="index",
      :disabled="disabled",
      hideDetails
    )
</template>

<script>
import { get, isUndefined } from 'lodash';
import { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES_TO_MASKS } from '@/constants';

export default {
  model: {
    prop: 'right',
    event: 'change',
  },
  props: {
    right: {
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
    checkboxes() {
      const masks = USERS_RIGHTS_TYPES_TO_MASKS[this.right.type];

      if (masks) {
        return masks.map(mask => ({
          bind: {
            inputValue: this.getCheckboxValue(mask),
            label: mask,
          },
          on: {
            change: value => this.changeCheckboxValue(value, mask),
          },
        }));
      }

      return [
        {
          key: 'right',
          bind: {
            inputValue: this.getCheckboxValue(),
          },
          on: {
            change: value => this.changeCheckboxValue(value),
          },
        },
      ];
    },
  },
  methods: {
    getCheckboxValue(rightMask = USERS_RIGHTS_MASKS.default) {
      const { right, role, changedRole } = this;

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
     * @param {number} mask
     */
    changeCheckboxValue(value, mask) {
      this.$emit('change', value, this.role, this.right, mask);
    },
  },
};
</script>

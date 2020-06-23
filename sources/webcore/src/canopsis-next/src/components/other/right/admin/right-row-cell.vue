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
import { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES_TO_MASKS } from '@/constants';

import { getCheckboxValue } from '@/helpers/right';

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
        return masks.map((mask) => {
          const [label] = Object.entries(USERS_RIGHTS_MASKS)
            .find(([key, value]) => key !== 'default' && value === mask) || [''];

          return {
            bind: {
              label,

              inputValue: getCheckboxValue(this.right, this.role, this.changedRole, mask),
            },
            on: {
              change: value => this.changeCheckboxValue(value, mask),
            },
          };
        });
      }

      return [
        {
          bind: {
            inputValue: getCheckboxValue(this.right, this.role, this.changedRole),
          },
          on: {
            change: value => this.changeCheckboxValue(value),
          },
        },
      ];
    },
  },
  methods: {
    /**
     * Change checkbox value
     *
     * @param {boolean} value
     * @param {number} mask
     */
    changeCheckboxValue(value, mask = USERS_RIGHTS_MASKS.default) {
      this.$emit('change', value, this.role, this.right, mask);
    },
  },
};
</script>

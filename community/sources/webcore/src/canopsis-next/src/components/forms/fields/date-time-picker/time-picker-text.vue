<template lang="pug">
  v-text-field.time-picker__select(
    :value="value",
    :label="label",
    placeholder="−−:−−",
    hide-details,
    disabled
  )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: null,
    },
    stepsInHours: {
      type: Number,
      default: 4,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.time-picker__select {
  display: inline-block;
  width: 56px;

  &.v-select--is-menu-active .v-input__slot {
    background: #686868;
  }

  & ::v-deep {
    .v-select__slot {
      padding: 0 5px;
    }
  }
}
</style>

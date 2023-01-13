<template lang="pug">
  v-layout(row)
    v-flex.pr-3(xs4)
      c-number-field(
        v-field="value.retry_count",
        :label="$t('common.retryCount')",
        :min="0",
        :name="countFieldName",
        :required="isRequired",
        :disabled="disabled"
      )
    v-flex(xs8)
      c-duration-field(
        v-field="value.retry_delay",
        :units-label="$t('common.unit')",
        :required="isRequired",
        :name="name",
        :disabled="isDurationDisabled",
        clearable
      )
</template>

<script>
import { isNumber } from 'lodash';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'retry',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    countFieldName() {
      return `${this.name}.count`;
    },

    isDurationDisabled() {
      return this.disabled || this.value.retry_count === 0;
    },

    isRequired() {
      const { retry_delay: retryDelay, retry_count: retryCount } = this.value;

      return this.required
        || isNumber(retryCount)
        || isNumber(retryDelay?.value)
        || Boolean(retryDelay?.unit);
    },
  },
};
</script>

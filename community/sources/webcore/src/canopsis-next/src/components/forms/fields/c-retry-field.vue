<template lang="pug">
  v-layout
    v-flex.pr-3(xs4)
      c-number-field(
        v-field="retry.count",
        :label="$t('common.retryCount')",
        :min="1",
        :name="countFieldName",
        :required="isRequired"
      )
    v-flex(xs8)
      c-duration-field(
        v-field="retry",
        :units-label="$t('common.unit')",
        :required="isRequired",
        :name="name",
        clearable
      )
</template>

<script>
import { isNumber } from 'lodash';

export default {
  inject: ['$validator'],
  model: {
    prop: 'retry',
    event: 'input',
  },
  props: {
    retry: {
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
  },
  computed: {
    countFieldName() {
      return `${this.name}.count`;
    },

    isRequired() {
      return this.required || isNumber(this.retry.count) || isNumber(this.retry.value) || Boolean(this.retry.unit);
    },
  },
};
</script>

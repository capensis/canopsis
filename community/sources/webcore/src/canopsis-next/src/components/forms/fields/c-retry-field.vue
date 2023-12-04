<template>
  <v-layout>
    <v-flex
      class="pr-3"
      xs4
    >
      <c-number-field
        v-field="value.retry_count"
        :label="$t('common.retryCount')"
        :min="0"
        :name="countFieldName"
        :required="isRequired"
        :disabled="disabled"
      />
    </v-flex>
    <v-flex xs8>
      <c-duration-field
        :duration="value.retry_delay"
        :units-label="$t('common.unit')"
        :required="isRequired"
        :name="name"
        :disabled="isDurationDisabled"
        clearable
        @input="updateDelay"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { isNumber } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
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
        || (isNumber(retryCount) && retryCount > 0)
        || isNumber(retryDelay?.value)
        || Boolean(retryDelay?.unit);
    },
  },
  methods: {
    updateDelay(duration) {
      if (duration.unit || duration.value) {
        this.updateField('retry_delay', duration);
      } else {
        this.updateModel({
          ...this.value,
          retry_delay: duration,
          retry_count: undefined,
        });
      }
    },
  },
};
</script>

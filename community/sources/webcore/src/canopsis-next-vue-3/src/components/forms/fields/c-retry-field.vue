<template>
  <c-layout class="c-retry-field" gap="3">
    <div class="c-retry-field__count">
      <c-number-field
        :value="value.retry_count"
        :label="$t('common.retryCount')"
        :min="0"
        :name="countFieldName"
        :required="isRequired"
        :disabled="disabled"
        @update:value="updateField('retry_count', $event)"
      />
    </div>

    <div class="c-retry-field__delay">
      <c-duration-field
        :duration="value.retry_delay"
        :units-label="$t('common.unit')"
        :required="isRequired"
        :name="name"
        :disabled="isDurationDisabled"
        clearable
        @update:duration="updateDelay"
      />
    </div>
  </c-layout>
</template>

<script>
import { isNumber } from 'lodash';

import { formMixin } from '@/mixins/form';

export default {
  mixins: [formMixin],
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

  emits: ['update:value'],

  computed: {
    countFieldName() {
      return `${this.name}.count`;
    },

    isDurationDisabled() {
      return this.disabled || this.value.retry_count === 0;
    },

    isRequired() {
      const { retry_delay: retryDelay, retry_count: retryCount } = this.value;

      return (
        this.required ||
        (isNumber(retryCount) && retryCount > 0) ||
        isNumber(retryDelay?.value) ||
        Boolean(retryDelay?.unit)
      );
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

<style lang="scss" scoped>
.c-retry-field {
  &__count {
    flex: 0 0 calc(100% / 3 - 8px);
    max-width: calc(100% / 3 - 8px);
    min-width: 0;
  }

  &__delay {
    flex: 0 0 calc(100% * 2 / 3 - 4px);
    max-width: calc(100% * 2 / 3 - 4px);
    min-width: 0;
  }
}
</style>

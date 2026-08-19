<template>
  <c-layout class="availability-value-filter-field" gap="3" align="end">
    <availability-value-filter-method-field
      v-model="method"
      class="availability-value-filter-field__method"
    />
    <template v-if="method">
      <c-number-field
        v-if="isPercentType"
        :value="modelValue.value"
        class="availability-value-filter-field__value"
        prefix="%"
        hide-details
        @update:value="updateField('value', $event)"
      />
      <c-splitted-duration-field
        v-else
        :model-value="modelValue.value"
        :max-value="maxSeconds"
        class="availability-value-filter-field__value"
        @update:model-value="updateField('value', $event)"
      />
      <c-action-btn
        :tooltip="$t('common.clear')"
        type="delete"
        icon="clear"
        small
        @click="clearValue"
      />
    </template>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { AVAILABILITY_SHOW_TYPE } from '@/constants';

import { useModelField } from '@/hooks/form';

import AvailabilityValueFilterMethodField from './availability-value-filter-method-field.vue';

export default {
  components: { AvailabilityValueFilterMethodField },

  props: {
    modelValue: {
      type: Object,
      required: false,
    },
    showType: {
      type: Number,
      default: AVAILABILITY_SHOW_TYPE.percent,
    },
    maxSeconds: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateModel, updateField } = useModelField(props, emit);
    const isPercentType = computed(() => props.showType === AVAILABILITY_SHOW_TYPE.percent);

    const method = computed({
      get() {
        return props.modelValue?.method;
      },
      set(newMethod) {
        updateModel({
          ...props.modelValue,
          value: props.modelValue?.value ?? 0,
          method: newMethod,
        });
      },
    });

    const clearValue = () => updateModel(undefined);

    return {
      method,
      isPercentType,
      updateField,
      clearValue,
    };
  },
};
</script>

<style lang="scss">
.availability-value-filter-field {
  flex-grow: 0;

  &__method {
    width: 140px;
    flex-grow: 0;
    flex-shrink: 0;
  }

  &__value {
    width: 160px;
  }
}
</style>

<template>
  <c-combobox
    :model-value="inputValue"
    :items="eventTypes"
    :label="label ?? $t('common.eventType')"
    :return-object="false"
    :name="name"
    :required="required"
    item-title="text"
    item-value="value"
    v-bind="$attrs"
    @update:model-value="updateValue"
  >
    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps" />
    </template>
  </c-combobox>
</template>

<script>
import { computed } from 'vue';

import { EVENT_TYPES } from '@/constants';

export default {
  inheritAttrs: false,

  props: {
    value: {
      type: [String, Array],
      default: undefined,
    },
    modelValue: {
      type: [String, Array],
      default: undefined,
    },
    label: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'event_type',
    },
    types: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue', 'update:value'],

  setup(props, { emit }) {
    const inputValue = computed(() =>
      props.modelValue !== undefined ? props.modelValue : props.value,
    );

    const eventTypes = computed(() =>
      props.types?.length
        ? props.types
        : Object.values(EVENT_TYPES).map(value => ({
            value,
            text: value,
          })),
    );

    const updateValue = (value) => {
      emit('update:modelValue', value);
      emit('update:value', value);
    };

    return {
      inputValue,
      eventTypes,
      updateValue,
    };
  },
};
</script>

<template>
  <c-radio-group
    v-model="currentValue"
    :label="label"
    :name="name"
    class="mt-0"
  >
    <v-radio
      v-for="item in items"
      :key="item.value"
      :label="item.label"
      :value="item.value"
    />
  </c-radio-group>
</template>

<script>
import { computed } from 'vue';

import { IDLE_RULE_ALARM_CONDITIONS } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: IDLE_RULE_ALARM_CONDITIONS.lastEvent,
    },
    value: {
      type: String,
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'alarm_condition',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const currentValue = computed({
      get: () => props.modelValue ?? props.value,
      set: (value) => {
        emit('update:modelValue', value);
      },
    });

    const items = computed(() =>
      Object.values(IDLE_RULE_ALARM_CONDITIONS).map(value => ({
        value,
        label: t(`idleRules.alarmConditions.${value}`),
      })),
    );

    return {
      currentValue,
      items,
    };
  },
};
</script>

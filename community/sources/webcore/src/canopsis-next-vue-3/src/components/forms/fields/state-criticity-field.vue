<template>
  <c-layout>
    <v-btn-toggle
      :model-value="selectedValue"
      :mandatory="mandatory"
      class="state-criticity-field"
      density="compact"
      divided
      @update:model-value="updateValue"
    >
      <v-btn
        v-for="{ color, state, text, textColor } in availableStates"
        :key="state"
        :value="state"
        :style="getButtonStyle({ color, state, textColor })"
        class="state-criticity-field__button"
        variant="flat"
      >
        {{ text }}
      </v-btn>
    </v-btn-toggle>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STATES, ALARM_STATES_TEXTS } from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Number,
      default: undefined,
    },
    value: {
      type: Number,
      default: null,
    },
    mandatory: {
      type: Boolean,
      default: false,
    },
    stateValues: {
      type: Object,
      default: () => ALARM_STATES,
    },
  },

  emits: ['update:modelValue', 'update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const selectedValue = computed(() => props.modelValue ?? props.value);

    const availableStates = computed(() => Object.entries(props.stateValues).map(([key, state]) => ({
      text: t(`modals.createChangeStateEvent.states.${key}`),
      state,
      color: getAlarmStateColor(state),
      textColor: `var(--v-state-${ALARM_STATES_TEXTS[state]}-text-base)`,
    })));

    const updateValue = (value) => {
      emit('update:modelValue', value);
      emit('update:value', value);
    };

    const getButtonStyle = ({ color, state, textColor }) => (
      {
        backgroundColor: color,
        color: textColor,
        opacity: selectedValue.value === state ? 1 : 0.45,
      }
    );

    return {
      selectedValue,
      availableStates,
      updateValue,
      getButtonStyle,
    };
  },
};
</script>

<style lang="scss" scoped>
.state-criticity-field {
  :deep(.v-btn--active) {
    opacity: 1;
  }
}

.state-criticity-field__button {
  min-width: 88px;
  font-weight: 500;
}
</style>

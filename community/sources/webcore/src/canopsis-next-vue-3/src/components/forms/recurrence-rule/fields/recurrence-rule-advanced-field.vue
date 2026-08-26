<template>
  <c-text-field
    :model-value="modelValue"
    :label="$attrs.label"
    :name="name"
    :validation-rules="validationRules"
    persistent-hint
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #append>
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        icon="help"
        max-width="250"
        left
      />
    </template>
  </c-text-field>
</template>

<script setup>
import { isNaN, isNumber } from 'lodash';
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  helpText: {
    type: String,
    default: '',
  },
  name: {
    type: String,
    default: 'regex',
  },
  min: {
    type: Number,
    default: undefined,
  },
  max: {
    type: Number,
    default: undefined,
  },
  negative: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const invalidMessage = computed(() => {
  const isMinZero = props.min === 0;

  if (!props.negative || isMinZero) {
    return t('recurrenceRule.errors.invalidRange', {
      min: props.negative && isMinZero ? -props.max : props.min,
      max: props.max,
    });
  }

  return t('recurrenceRule.errors.invalidRangeNegative', { min: props.min, max: props.max });
});

const isValueValid = (value = '') =>
  !value.length ||
  value.split(',').every((item) => {
    const preparedValue = props.negative ? Math.abs(+item) : +item;

    if (isNaN(preparedValue) || !isNumber(preparedValue)) {
      return false;
    }

    return isNumber(props.min) && isNumber(props.max)
      ? preparedValue >= props.min && preparedValue <= props.max
      : true;
  });

const validationRules = [value => isValueValid(value) || invalidMessage.value];
</script>

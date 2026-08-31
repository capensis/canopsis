<template>
  <c-radio-group :model-value="type" :label="t('common.end')" @update:model-value="updateType">
    <v-radio :value="types.never" color="primary">
      <template #label>
        {{ t('recurrenceRule.never') }}
      </template>
    </v-radio>

    <v-radio :value="types.date" class="mb-0" color="primary">
      <template #label>
        <c-layout gap="2" align-center>
          <span>{{ t('recurrenceRule.on') }}</span>

          <date-time-splitted-picker-field
            :value="modelValue.until"
            :rules="dateTimeSplittedRules"
            :placeholder="t('common.date')"
            :disabled="!isDateType"
            name="until"
            @update:value="updateModelField('until', $event)"
          />
        </c-layout>
      </template>
    </v-radio>

    <v-radio :value="types.after" class="mb-0" color="primary">
      <template #label>
        <c-layout gap="2" align-center>
          <span>{{ t('recurrenceRule.after') }}</span>

          <c-number-field
            :value="modelValue.count"
            :disabled="!isAfterType"
            :required="isAfterType"
            class="recurrence-rule-end-field__count"
            name="count"
            @update:value="updateModelField('count', $event)"
          />

          <span class="text-lowercase">
            {{ t('recurrenceRule.occurrence', modelValue.count || 1) }}
          </span>
        </c-layout>
      </template>
    </v-radio>
  </c-radio-group>
</template>

<script setup>
import { isNumber, omit } from 'lodash';
import { computed, ref } from 'vue';

import { RECURRENCE_RULE_END_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import DateTimeSplittedPickerField from '@/components/forms/fields/date-time-picker/date-time-splitted-picker-field.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  types: {
    type: Object,
    default: () => RECURRENCE_RULE_END_TYPES,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const type = ref(
  props.modelValue.until
    ? RECURRENCE_RULE_END_TYPES.date
    : isNumber(props.modelValue.count)
      ? RECURRENCE_RULE_END_TYPES.after
      : RECURRENCE_RULE_END_TYPES.never,
);

const isDateType = computed(() => type.value === props.types.date);
const isAfterType = computed(() => type.value === props.types.after);
const dateTimeSplittedRules = computed(() => ({ required: isDateType.value }));

const updateModelField = (field, value) => emit('update:modelValue', {
    ...props.modelValue,
    [field]: value,
  });

const updateType = (value) => {
  type.value = value;

  const deletingFields = [];

  if (!isDateType.value) {
    deletingFields.push('until');
  }

  if (!isAfterType.value) {
    deletingFields.push('count');
  }

  emit('update:modelValue', omit(props.modelValue, deletingFields));
};
</script>

<style lang="scss" scoped>
.recurrence-rule-end-field__count {
  width: 96px;
}
</style>

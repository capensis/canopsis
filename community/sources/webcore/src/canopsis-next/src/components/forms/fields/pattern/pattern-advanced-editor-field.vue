<template>
  <c-json-field
    :value="value"
    :label="t('pattern.advancedEditor')"
    :readonly="disabled"
    :name="name"
    validate-on="button"
    rows="10"
    @update:value="updatePatternsFromJSON"
  >
    <template #append="{ jsonString, errorMessage: jsonErrorMessage }">
      <v-fade-transition>
        <pattern-fix-with-ai-btn v-if="errorMessage || jsonErrorMessage" :json-string="jsonString" />
      </v-fade-transition>
    </template>
  </c-json-field>
</template>

<script setup>
import { isArray } from 'lodash';
import { computed, toRef } from 'vue';

import {
  isExtraInfosRuleType,
  isInfosRuleType,
  isValidPatternRule,
  isObjectRuleType,
} from '@/helpers/entities/pattern/form';

import { useFieldError, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import PatternFixWithAiBtn from './pattern-fix-with-ai-btn.vue';

const props = defineProps({
  value: {
    type: Array,
    required: true,
  },
  attributes: {
    type: Array,
    required: true,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  name: {
    type: String,
    default: undefined,
  },
});

const emit = defineEmits(['update:value']);

const { t } = useI18n();
const formContext = useFormContext();

const fieldName = toRef(props, 'name');
const errorMessage = useFieldError(fieldName);

const disabledAttributes = computed(() => props.attributes.reduce((acc, { value, options }) => {
  if (options?.disabled) {
    acc.push(value);
  }

  return acc;
}, []));

const setErrorMessage = message => formContext.setFieldError(fieldName.value, message);

const isDisabledField = field => disabledAttributes.value.includes(field);

const isValidRuleField = ({ field }) => props.attributes.some(({ value, options }) => {
  if (isDisabledField(value)) {
    return false;
  }

  if (
    isInfosRuleType(options?.type) ||
    isExtraInfosRuleType(options?.type) ||
    isObjectRuleType(options?.type)
  ) {
    return field.startsWith(`${value}.`);
  }

  return value === field;
});

const isValidPatternRuleField = rule => isValidPatternRule(rule) && isValidRuleField(rule);

const isValidPatternRules = rules => isArray(rules) && rules.length > 0 && rules.every(isValidPatternRuleField);

const isValidPatterns = patterns => isArray(patterns) && patterns.every(isValidPatternRules);

const updatePatternsFromJSON = (patterns) => {
  if (isValidPatterns(patterns)) {
    setErrorMessage(undefined);
    emit('update:value', patterns);

    return;
  }

  setErrorMessage(t('pattern.errors.invalidPatterns'));
};
</script>

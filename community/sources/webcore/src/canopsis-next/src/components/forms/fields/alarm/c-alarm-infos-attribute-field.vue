<template>
  <c-layout column>
    <c-combobox
      :model-value="rule"
      :items="rules"
      :disabled="disabled"
      :label="label || t('common.rule')"
      :name="ruleName"
      :loading="pending"
      :clearable="!required"
      :required="ruleRequired"
      item-title="name"
      item-value="_id"
      return-object
      @update:model-value="rule = $event"
    />
    <c-combobox
      :model-value="value.dictionary"
      :items="infosName"
      :disabled="disabled"
      :label="label || t('common.dictionary')"
      :name="dictionaryName"
      :loading="pending"
      :clearable="!required"
      :required="dictionaryRequired"
      @update:model-value="updateField('dictionary', $event)"
    />
  </c-layout>
</template>

<script setup>
import { computed, useModel } from 'vue';

import { useI18n } from '@/hooks/i18n';

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
  rules: {
    type: Array,
    default: () => [],
  },
  label: {
    type: String,
    default: '',
  },
  name: {
    type: String,
    default: 'infos',
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  pending: {
    type: Boolean,
    default: false,
  },
  required: {
    type: Boolean,
    default: false,
  },
});

defineEmits(['update:value']);

const { t } = useI18n();
const value = useModel(props, 'value');

const updateField = (field, fieldValue) => {
  value.value = {
    ...value.value,
    [field]: fieldValue,
  };
};

const rule = computed({
  get: () => props.rules.find(item => item._id === value.value.rule) ?? null,
  set: selectedRule => updateField('rule', selectedRule?._id),
});

const ruleName = computed(() => `${props.name}.ruleId`);
const dictionaryName = computed(() => `${props.name}.dictionary`);
const ruleRequired = computed(() => props.required || Boolean(value.value.dictionary));
const dictionaryRequired = computed(() => props.required || Boolean(rule.value));
const infosName = computed(() => rule.value?.infos?.map(info => info.name) ?? []);
</script>

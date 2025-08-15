<template>
  <v-layout :column="column">
    <v-flex :xs6="row">
      <v-combobox
        v-if="combobox"
        v-validate="'required'"
        :value="value.dictionary"
        :items="items"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :name="dictionaryName"
        :error-messages="errors.collect(dictionaryName)"
        :loading="pending"
        :hide-details="row"
        item-text="value"
        item-value="value"
        return-object
        @input="updateDictionary"
      />
      <v-text-field
        v-else
        v-validate="'required'"
        :value="value.dictionary"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :error-messages="errors.collect(dictionaryName)"
        :name="dictionaryName"
        :hide-details="row"
        @input="updateDictionary"
      />
    </v-flex>
    <v-flex
      :class="{ 'pl-3': row }"
      :xs6="row"
    >
      <v-select
        v-field="value.field"
        v-validate="'required'"
        :items="fieldItems"
        :disabled="disabled || !value.dictionary"
        :label="label || $t('common.field')"
        :name="fieldName"
        :error-messages="errors.collect(fieldName)"
        :hide-details="row"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERN_RULE_INFOS_FIELDS, PATTERN_FIELD_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    items: {
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
    divider: {
      type: String,
      default: '.',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    row: {
      type: Boolean,
      required: false,
    },
    column: {
      type: Boolean,
      required: false,
    },
    pending: {
      type: Boolean,
      required: false,
    },
  },
  setup(props, { emit }) {
    // Composables
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    // Computed properties
    const dictionaryName = computed(() => `${props.name}.dictionary`);
    const fieldName = computed(() => `${props.name}.field`);
    const fieldItems = computed(() => [
      {
        text: t('common.name'),
        value: PATTERN_RULE_INFOS_FIELDS.name,
      },
      {
        text: t('common.value'),
        value: PATTERN_RULE_INFOS_FIELDS.value,
      },
    ]);

    const updateDictionary = (infos = {}) => updateModel({
      ...props.value,

      dictionary: infos?.value ?? infos,
      definedType: infos?.definedType,
      fieldType: infos?.definedType ?? PATTERN_FIELD_TYPES.string,
    });

    return {
      dictionaryName,
      fieldName,
      fieldItems,
      updateDictionary,
    };
  },
};
</script>

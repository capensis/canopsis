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
        @input="updateInfosDictionary"
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
        @input="updateInfosDictionary"
      />
    </v-flex>
    <v-flex
      :class="{ 'pl-3': row }"
      :xs6="row"
    >
      <v-select
        v-validate="'required'"
        :value="value.field"
        :items="fieldItems"
        :disabled="disabled || !value.dictionary"
        :label="label || $t('common.field')"
        :name="fieldName"
        :error-messages="errors.collect(fieldName)"
        :hide-details="row"
        @input="updateInfosField"
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
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

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

    /**
     * Updates `value.dictionary` and related type fields based on the provided infos.
     *
     * When a combobox item is selected, `infos` is an object with `value` and optional
     * `definedType`. When a plain text is entered, `infos` is a string.
     *
     * Triggers model update (emits `input`).
     *
     * @param {Object|string} [infos={}]
     * @param {string} [infos.value] Dictionary name when an object is provided
     * @param {string} [infos.definedType] Field defined type coming from the dictionary item
     * @returns {void}
     */
    const updateInfosDictionary = (infos = {}) => updateModel({
      ...props.value,

      dictionary: infos?.value ?? infos,
      definedType: infos?.definedType,
      fieldType: infos?.definedType ?? PATTERN_FIELD_TYPES.string,
    });

    /**
     * Updates `value.field` and keeps `definedType` only when field is `value`.
     *
     * @param {'name'|'value'} field Field key from `PATTERN_RULE_INFOS_FIELDS`
     * @returns {void}
     */
    const updateInfosField = field => updateModel({
      ...props.value,
      fieldType: field === PATTERN_RULE_INFOS_FIELDS.value ? props.value.definedType : null,
      field,
    });

    return {
      dictionaryName,
      fieldName,
      fieldItems,
      updateInfosDictionary,
      updateInfosField,
    };
  },
};
</script>

<template>
  <v-layout column>
    <v-radio-group
      v-field="form.selectedType"
      v-validate="'required'"
      :error-messages="errors.collect('type')"
      class="mt-0"
      name="type"
    >
      <v-radio
        v-for="type in types"
        :key="type.value"
        :value="type.value"
        color="success"
      >
        <template #label>
          <v-layout justify-space-between>
            <span>
              {{ type.text }}
            </span>
            <span>
              {{ type.description }}
            </span>
          </v-layout>
        </template>
      </v-radio>
    </v-radio-group>
    <v-flex v-if="isCustomTypeSelected">
      <v-layout class="ml-8" column>
        <v-flex>
          <v-radio-group
            v-field="form.selectedSeparator"
            class="mt-0"
            hide-details
            @change="validateCustomSeparator"
          >
            <v-radio
              v-for="separator in separatorsWithCustom"
              :key="separator.value"
              v-bind="separator"
              color="success"
            >
              <template v-if="!separator.label" #label>
                <v-text-field
                  v-field="form.customSeparator"
                  v-validate="rule"
                  :error-messages="errors.collect('separator')"
                  :disabled="!isCustomSeparatorSelected"
                  name="separator"
                />
              </template>
            </v-radio>
          </v-radio-group>
        </v-flex>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

const CUSTOM_SEPARATOR_VALUE = 'custom';

export default {
  inject: ['$validator'],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    tableSeparator: {
      type: String,
      default: '',
    },

  },
  setup(props) {
    const { t } = useI18n();

    const validator = useValidator();

    const isCustomTypeSelected = computed(() => (
      props.form.selectedType === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));

    const isCustomSeparatorSelected = computed(() => props.form.selectedSeparator === CUSTOM_SEPARATOR_VALUE);

    const types = computed(() => Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES).map(type => ({
      value: type,
      text: t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${type}.text`),
      description: t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${type}.description`),
    })));

    const separatorsWithCustom = computed(() => [
      ...Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS).map(separator => ({
        value: separator,
        label: separator,
      })).filter(separator => separator.value !== CSV_SEPARATORS_TO_SYMBOLS[props.tableSeparator]),

      { value: CUSTOM_SEPARATOR_VALUE, class: 'radio-string-array-separator-custom' },
    ]);

    const rule = computed(() => (
      isCustomSeparatorSelected.value
        ? { required: true, forbidden_separator: true }
        : {}
    ));

    /**
     * Validates the custom separator field
     */
    const validateCustomSeparator = () => validator.validate('separator');

    return {
      isCustomTypeSelected,
      isCustomSeparatorSelected,
      types,
      separatorsWithCustom,
      rule,
      validateCustomSeparator,
    };
  },
};
</script>

<style lang="scss">
.radio-string-array-separator-custom {
  align-items: start;

  .v-input {
    margin: 0;
    padding: 0;
  }
}
</style>

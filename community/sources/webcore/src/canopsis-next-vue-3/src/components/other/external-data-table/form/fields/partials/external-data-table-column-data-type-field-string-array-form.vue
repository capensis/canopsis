<template>
  <c-layout column>
    <c-radio-group
      :model-value="form.selectedType"
      class="mt-0"
      name="type"
      required
      @update:model-value="updateField('selectedType', $event)"
    >
      <v-radio
        v-for="type in types"
        :key="type.value"
        :value="type.value"
        color="success"
      >
        <template #label>
          <c-layout justify="space-between">
            <span>
              {{ type.text }}
            </span>
            <span>
              {{ type.description }}
            </span>
          </c-layout>
        </template>
      </v-radio>
    </c-radio-group>
    <template v-if="isCustomTypeSelected">
      <c-layout class="ml-8" column>
        <c-radio-group
          :model-value="form.selectedSeparator"
          class="mt-0"
          hide-details
          @update:model-value="updateSelectedSeparator"
        >
          <v-radio
            v-for="separator in separatorsWithCustom"
            :key="separator.value"
            v-bind="separator"
            color="success"
          >
            <template v-if="!separator.label" #label>
              <c-text-field
                :model-value="form.customSeparator"
                :disabled="!isCustomSeparatorSelected"
                :required="isCustomSeparatorSelected"
                name="separator"
                @update:model-value="updateField('customSeparator', $event)"
              />
            </template>
          </v-radio>
        </c-radio-group>
      </c-layout>
    </template>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

const CUSTOM_SEPARATOR_VALUE = 'custom';

export default {
  emits: ['update:form'],
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);

    const isCustomTypeSelected = computed(
      () =>
        props.form.selectedType === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom,
    );

    const isCustomSeparatorSelected = computed(
      () => props.form.selectedSeparator === CUSTOM_SEPARATOR_VALUE,
    );

    const types = computed(() =>
      Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES).map(type => ({
        value: type,
        text: t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${type}.text`),
        description: t(
          `externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${type}.description`,
        ),
      })),
    );

    const separatorsWithCustom = computed(() => [
      ...Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS)
        .map(separator => ({
          value: separator,
          label: separator,
        }))
        .filter(separator => separator.value !== CSV_SEPARATORS_TO_SYMBOLS[props.tableSeparator]),

      { value: CUSTOM_SEPARATOR_VALUE, class: 'radio-string-array-separator-custom' },
    ]);

    /**
     * Validates the custom separator field
     */
    const validateCustomSeparator = () => validator.validate('separator');

    /**
     * Updates selected separator and validates custom separator field.
     *
     * @param {*} value
     */
    const updateSelectedSeparator = (value) => {
      updateField('selectedSeparator', value);
      validateCustomSeparator();
    };

    return {
      isCustomTypeSelected,
      isCustomSeparatorSelected,
      types,
      separatorsWithCustom,
      updateField,
      updateSelectedSeparator,
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

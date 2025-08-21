<template>
  <v-menu
    v-model="isOpen"
    ref="menuElement"
    :close-on-content-click="false"
    max-width="400"
    min-width="350"
    offset-y
  >
    <template #activator="{ on }">
      <c-chip
        class="px-2"
        color="grey"
        text-color="blue darken-1"
        outlined
        v-on="on"
      >
        <span v-if="!value.string_array_type" class="grey--text">
          {{ $t('externalData.tableColumnDataTypesAdditionalChips.stringArray.selectSeparator') }}
        </span>
        <span v-else-if="isCustomValueTypeSelected">
          <span class="grey--text mr-2">
            {{ chipPrefix }}:
          </span>
          <span>
            {{ value.string_array_separator }}
          </span>
        </span>
        <span v-else>
          <span class="grey--text">
            {{ chipPrefix }}
          </span>
        </span>
      </c-chip>
    </template>
    <v-card>
      <v-card-text>
        <v-layout column>
          <v-radio-group
            v-model="selectedType"
            v-validate="'required'"
            :error-messages="errors.collect('type')"
            class="mt-0"
            name="type"
            @change="callMenuResize"
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
                  v-model="selectedSeparator"
                  class="mt-0"
                  hide-details
                >
                  <v-radio
                    v-for="separator in separatorsWithCustom"
                    :key="separator.value"
                    v-bind="separator"
                    color="success"
                  >
                    <template v-if="!separator.label" #label>
                      <v-text-field
                        v-model="customSeparator"
                        v-validate="'required|forbidden_separator'"
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
      </v-card-text>
      <v-card-actions class="pa-4 pt-0 justify-end">
        <v-btn
          depressed
          text
          @click="cancel"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          color="success"
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script>
import { ref, computed, onBeforeMount } from 'vue';

import {
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';

const CUSTOM_SEPARATOR_VALUE = 'custom';

// TODO: rewrite it
const getDefaultSeparator = (value, tableSeparator) => {
  if (
    value.string_array_separator
    && !Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS)
      .includes(value.string_array_separator)
  ) {
    return CUSTOM_SEPARATOR_VALUE;
  }

  if (!value.string_array_separator) {
    return CSV_SEPARATORS_TO_SYMBOLS[tableSeparator]
      !== EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS.comma
      ? EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS.comma
      : EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS.semicolon;
  }

  return value.string_array_separator;
};

export default {
  $_veeValidate: {
    validator: 'new',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    tableSeparator: {
      type: String,
      default: '',
    },
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const validator = useValidator();

    const menuElement = ref(null);
    const isOpen = ref(false);

    const defaultSeparator = getDefaultSeparator(props.value, props.tableSeparator);

    const selectedType = ref(props.value.string_array_type ?? null);
    const selectedSeparator = ref(defaultSeparator);
    const customSeparator = ref(defaultSeparator === CUSTOM_SEPARATOR_VALUE ? props.value.string_array_separator : '');

    const isCustomTypeSelected = computed(() => (
      selectedType.value === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));

    const isCustomValueTypeSelected = computed(() => (
      props.value.string_array_type === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));

    const isCustomSeparatorSelected = computed(() => selectedSeparator.value === CUSTOM_SEPARATOR_VALUE);

    const chipPrefix = computed(() => (
      isCustomValueTypeSelected.value
        ? t('externalData.tableColumnDataTypesAdditionalChips.stringArray.separator')
        : t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${props.value.string_array_type}.text`)
    ));

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

    const submit = async () => {
      const isValid = await validator.validateAll();

      if (isValid) {
        updateModel({
          string_array_type: selectedType.value,
          string_array_separator: selectedSeparator.value === CUSTOM_SEPARATOR_VALUE
            ? customSeparator.value
            : selectedSeparator.value,
        });
      }
    };

    const cancel = () => isOpen.value = false;

    const callMenuResize = () => menuElement.value?.onResize?.();

    const extendValidator = () => validator.extend('forbidden_separator', {
      getMessage: () => t('externalData.tableColumnDataTypesAdditionalChips.forbiddenSeparator'),
      validate: value => value !== CSV_SEPARATORS_TO_SYMBOLS[props.tableSeparator],
    });

    onBeforeMount(extendValidator);

    return {
      menuElement,
      isOpen,
      selectedType,
      selectedSeparator,
      customSeparator,
      chipPrefix,
      isCustomTypeSelected,
      isCustomValueTypeSelected,
      isCustomSeparatorSelected,
      types,
      separatorsWithCustom,

      submit,
      cancel,
      callMenuResize,
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

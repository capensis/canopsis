<template>
  <v-menu
    v-model="isOpen"
    :close-on-content-click="false"
    max-width="400"
    min-width="350"
    offset-y
  >
    <template #activator="{ on }">
      <c-chip
        class="px-2"
        color="grey"
        text-color="blue"
        outlined
        v-on="on"
      >
        <span>
          JSON
        </span>
      </c-chip>
    </template>
    <v-card class="pa-4">
      <v-card-text>
        <v-layout column>
          <v-radio-group
            v-model="selectedType"
            class="mt-0"
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
                >
                  <v-radio
                    v-for="separator in separatorsWithCustom"
                    :key="separator.value"
                    :label="separator.text"
                    :value="separator.value"
                    color="success"
                  >
                    <template v-if="!separator.text" #label>
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
      <v-card-actions class="justify-end">
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
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form';
import { useValidator } from '@/hooks/validator/validator';

const CUSTOM_SEPARATOR_VALUE = 'custom';

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

    const isOpen = ref(false);

    const defaultSeparator = props.value.string_array_separator
     || !Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS)
       .includes(props.value.string_array_separator)
      ? CUSTOM_SEPARATOR_VALUE
      : props.value.string_array_separator ?? EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS.comma;

    const selectedType = ref(props.value.string_array_type ?? null);
    const selectedSeparator = ref(defaultSeparator);
    const customSeparator = ref(defaultSeparator === CUSTOM_SEPARATOR_VALUE ? props.value.string_array_separator : '');

    const isCustomTypeSelected = computed(() => (
      selectedType.value === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));
    const isCustomSeparatorSelected = computed(() => selectedSeparator.value === CUSTOM_SEPARATOR_VALUE);

    const types = computed(() => Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES).map(type => ({
      value: type,
      text: t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${type}.text`),
    })));

    const separatorsWithCustom = computed(() => [
      ...Object.values(EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_SEPARATORS).map(separator => ({
        value: separator,
        text: separator,
      })),

      { value: CUSTOM_SEPARATOR_VALUE },
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

    const extendValidator = () => {
      validator.extend('forbidden_separator', {
        getMessage: () => t('externalData.tableColumnDataTypesAdditionalChips.forbiddenSeparator'),
        validate: value => value === props.tableSeparator,
      });
    };

    onBeforeMount(extendValidator);

    return {
      isOpen,
      selectedType,
      selectedSeparator,
      customSeparator,
      isCustomTypeSelected,
      isCustomSeparatorSelected,
      types,
      separatorsWithCustom,

      submit,
      cancel,
    };
  },
};
</script>

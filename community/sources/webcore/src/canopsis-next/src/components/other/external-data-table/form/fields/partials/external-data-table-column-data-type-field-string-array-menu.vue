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
        <external-data-table-column-data-type-field-string-array-form
          v-model="form"
          :table-separator="tableSeparator"
          @input="callMenuResize"
        />
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
import { ref, computed, onBeforeMount, onBeforeUnmount } from 'vue';

import {
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR,
} from '@/constants';

import { getDefaultSeparator } from '@/helpers/entities/external-data-table/form';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

import ExternalDataTableColumnDataTypeFieldStringArrayForm from './external-data-table-column-data-type-field-string-array-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ExternalDataTableColumnDataTypeFieldStringArrayForm,
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
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const validator = useValidator();

    const menuElement = ref(null);
    const isOpen = ref(false);

    const defaultSeparator = getDefaultSeparator(props.value, props.tableSeparator);

    const form = ref({
      selectedType: props.value.string_array_type ?? null,
      selectedSeparator: defaultSeparator,
      customSeparator: defaultSeparator === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR
        ? props.value.string_array_separator
        : '',
    });

    const isCustomValueTypeSelected = computed(() => (
      form.value.selectedType === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));

    const isCustomSeparatorSelected = computed(() => (
      form.value.selectedSeparator === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR
    ));

    const chipPrefix = computed(() => (
      isCustomValueTypeSelected.value
        ? t('externalData.tableColumnDataTypesAdditionalChips.stringArray.separator')
        : t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${props.value.string_array_type}.text`)
    ));

    /**
     * Handles form submission by validating the form data and updating the model
     * Closes the menu if validation passes
     */
    const submit = async () => {
      const isValid = await validator.validateAll();

      if (isValid) {
        updateModel({
          string_array_type: form.value.selectedType,
          string_array_separator:
            isCustomSeparatorSelected.value
              ? form.value.customSeparator
              : form.value.selectedSeparator,
        });
        isOpen.value = false;
      }
    };

    /**
     * Handles form cancellation by closing the menu without saving changes
     */
    const cancel = () => {
      isOpen.value = false;
    };

    /**
     * Triggers menu resize to adjust position and size when content changes
     */
    const callMenuResize = () => menuElement.value?.onResize?.();

    /**
     * Extends the validator with a custom 'forbidden_separator' rule
     * Prevents using the same separator as the table's CSV separator
     */
    const extendValidator = () => validator.extend('forbidden_separator', {
      getMessage: () => t('externalData.tableColumnDataTypesAdditionalChips.forbiddenSeparator'),
      validate: value => (
        !isCustomSeparatorSelected.value
        || value !== CSV_SEPARATORS_TO_SYMBOLS[props.tableSeparator]
      ),
    });

    const {
      attachRequiredRule,
      detachRequiredRule,
    } = useValidationAttachRequired(props.name);

    onBeforeMount(() => {
      extendValidator();
      attachRequiredRule();
    });

    onBeforeUnmount(detachRequiredRule);

    return {
      menuElement,
      isOpen,
      form,
      isCustomValueTypeSelected,
      chipPrefix,
      submit,
      cancel,
      callMenuResize,
    };
  },
};
</script>

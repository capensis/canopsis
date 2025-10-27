<template>
  <v-menu
    v-model="isOpen"
    ref="menuElement"
    :close-on-content-click="false"
    :disabled="disabled"
    max-width="400"
    min-width="350"
    offset-y
  >
    <template #activator="{ on }">
      <slot
        v-bind="{ on }"
        name="selection"
      />
    </template>
    <v-form @submit.prevent="submit">
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
            :disabled="isDisabled"
            :loading="submitting"
            color="success"
            type="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-menu>
</template>

<script>
import { ref, computed, onBeforeMount } from 'vue';

import {
  CSV_SEPARATORS_TO_SYMBOLS,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR,
} from '@/constants';

import { getDefaultSeparator } from '@/helpers/entities/external-data-table/form';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';
import { useSubmittableForm } from '@/hooks/submittable-form';

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
    disabled: {
      type: Boolean,
      default: false,
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

    const isCustomSeparatorSelected = computed(() => (
      form.value.selectedSeparator === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR
    ));

    /**
     * Handles form submission by validating the form data and updating the model
     * Closes the menu if validation passes
     */
    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const isValid = await validator.validateAll();

        if (isValid) {
          updateModel({
            ...props.value,

            string_array_type: form.value.selectedType,
            string_array_separator:
              isCustomSeparatorSelected.value
                ? form.value.customSeparator
                : form.value.selectedSeparator,
          });
          isOpen.value = false;
        }
      },
    });

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

    onBeforeMount(extendValidator);

    return {
      menuElement,
      isOpen,
      form,
      submit,
      isDisabled,
      submitting,
      cancel,
      callMenuResize,
    };
  },
};
</script>

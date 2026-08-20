<template>
  <v-menu
    v-model="isOpen"
    ref="menuElement"
    :close-on-content-click="false"
    :disabled="disabled"
    min-width="480"
    max-width="520"
    offset="4"
  >
    <template #activator="{ props }">
      <slot v-bind="{ props }" name="selection" />
    </template>
    <v-form @submit.prevent="submit">
      <v-card>
        <v-card-text>
          <external-data-table-column-data-type-field-string-array-form
            v-model:form="form"
            :table-separator="tableSeparator"
            @update:form="callMenuResize"
          />
        </v-card-text>
        <v-card-actions class="pa-4 pt-0 justify-end">
          <v-btn variant="text" @click="cancel">
            {{ $t('common.cancel') }}
          </v-btn>
          <v-btn
            :disabled="submitting"
            :loading="submitting"
            color="primary"
            variant="elevated"
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
import { ref, computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR } from '@/constants';

import { getDefaultSeparator } from '@/helpers/entities/external-data-table/form';

import { useForm } from '@/hooks/form';
import { useModelField } from '@/hooks/form/model-field';

import ExternalDataTableColumnDataTypeFieldStringArrayForm from './external-data-table-column-data-type-field-string-array-form.vue';

export default {
  components: {
    ExternalDataTableColumnDataTypeFieldStringArrayForm,
  },
  emits: ['update:value'],
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
    const { updateModel } = useModelField(props, emit);

    const menuElement = ref(null);
    const isOpen = ref(false);

    const defaultSeparator = getDefaultSeparator(props.value, props.tableSeparator);

    const { values: form, handleSubmit, isSubmitting: submitting } = useForm({
      initialValues: {
        selectedType: props.value.string_array_type ?? null,
        selectedSeparator: defaultSeparator,
        customSeparator:
          defaultSeparator === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR
            ? props.value.string_array_separator
            : '',
      },
    });

    const isCustomSeparatorSelected = computed(
      () =>
        form.value.selectedSeparator ===
        EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_CUSTOM_SEPARATOR,
    );

    const submit = handleSubmit(() => {
      updateModel({
        ...props.value,

        string_array_type: form.value.selectedType,
        string_array_separator: isCustomSeparatorSelected.value
          ? form.value.customSeparator
          : form.value.selectedSeparator,
      });
      isOpen.value = false;
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

    return {
      menuElement,
      isOpen,
      form,
      submit,
      submitting,
      cancel,
      callMenuResize,
    };
  },
};
</script>

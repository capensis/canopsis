<template>
  <c-layout gap="3" column>
    <external-data-table-general-info-form v-if="!values._id" :form="externalDataTable" />
    <c-id-field v-else :value="values._id" disabled />

    <component
      v-for="component in components"
      :is="component.is"
      :key="component.key"
      v-bind="component.bind"
      required
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES, FORM_MAX_STRING_LENGTH } from '@/constants';

import { useFormContext } from '@/hooks/form';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import ExternalDataTableGeneralInfoForm from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';

export default {
  components: { ExternalDataTableGeneralInfoForm },

  props: {
    externalDataTable: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const form = useFormContext();
    const { values, errors, setValue } = form;

    const componentsMap = {
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string]: 'c-text-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean]: 'c-enabled-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number]: 'c-number-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray]: 'c-array-text-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime]: DateTimePickerField,
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp]: DateTimePickerField,
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp]: 'c-text-field',
    };

    const getErrorMessages = (name) => {
      const message = errors.value?.[name];

      return message ? [message] : [];
    };

    const components = computed(() => props.externalDataTable.column_configs.map((column) => {
        const is = componentsMap[column.type];
        const rules = {
          required: true,
        };
        const defaultBind = {
          name: column.name,
          label: column.name,
          errorMessages: getErrorMessages(column.name),
          required: true,
        };

        if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray) {
          return {
            is,
            key: column.name,
            bind: {
              ...defaultBind,
              maxLength: FORM_MAX_STRING_LENGTH,
              values: values.value[column.name],
              'onUpdate:values': value => setValue(column.name, value),
            },
          };
        }

        if (
          [
            EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string,
            EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp,
          ].includes(column.type)
        ) {
          defaultBind.maxLength = FORM_MAX_STRING_LENGTH;
        }

        if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number) {
          defaultBind.inputmode = 'decimal';
        }

        if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean) {
          return {
            is,
            key: column.name,
            bind: {
              ...defaultBind,
              modelValue: values.value[column.name],
              'onUpdate:modelValue': value => setValue(column.name, value),
            },
          };
        }

        if (
          [
            EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime,
            EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp,
          ].includes(column.type)
        ) {
          return {
            is,
            key: column.name,
            bind: {
              ...defaultBind,
              rules,
              value: values.value[column.name],
              'onUpdate:value': value => setValue(column.name, value),
            },
          };
        }

        if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number) {
          return {
            is,
            key: column.name,
            bind: {
              ...defaultBind,
              value: values.value[column.name],
              'onUpdate:value': value => setValue(column.name, value),
            },
          };
        }

        return {
          is,
          key: column.name,
          bind: {
            ...defaultBind,
            modelValue: values.value[column.name],
            'onUpdate:modelValue': value => setValue(column.name, value),
          },
        };
      }),
    );

    return {
      values,
      components,
    };
  },
};
</script>

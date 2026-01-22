<template>
  <v-layout column>
    <external-data-table-general-info-form v-if="!form._id" :form="externalDataTable" />
    <c-id-field
      v-else
      :value="form._id"
      disabled
    />

    <component
      v-for="component in components"
      :is="component.is"
      v-validate="component.rules"
      v-bind="component.bind"
      :key="component.key"
      required
      v-on="component.on"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES } from '@/constants';

import { useValidator } from '@/hooks/validator/validator';
import { useModelField } from '@/hooks/form/model-field';

import DateTimePickerField from '@/components/forms/fields/date-time-picker/date-time-picker-field.vue';
import ExternalDataTableGeneralInfoForm
  from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataTableGeneralInfoForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    externalDataTable: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);

    const componentsMap = {
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string]: 'v-text-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.boolean]: 'c-enabled-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number]: 'c-number-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray]: 'c-array-text-field',
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.datetime]: DateTimePickerField,
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.timestamp]: DateTimePickerField,
      [EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp]: 'v-text-field',
    };

    const components = computed(() => (
      props.externalDataTable.column_configs.map((column) => {
        const is = componentsMap[column.type];
        const defaultBind = {
          name: column.name,
          label: column.name,
          errorMessages: validator.errors.collect(column.name),
          required: true,
        };

        const rules = {
          required: true,
        };

        if (column.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray) {
          return {
            is,
            rules,
            key: column.name,
            bind: {
              ...defaultBind,
              maxLength: 255,
              values: props.form[column.name],
              class: 'mt-2',
            },
            on: {
              change: value => updateField(column.name, value),
            },
          };
        }

        if ([
          EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.string,
          EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.regexp,
        ].includes(column.type)
        ) {
          rules.max = 255; // TODO: move into constants
        }

        if (EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number) {
          // defaultBind['data-vv-validate-on'] = 'blur';
          defaultBind.inputmode = 'decimal';
        }

        return {
          is,
          rules,
          key: column.name,
          bind: {
            ...defaultBind,
            value: props.form[column.name],
          },
          on: {
            input: value => updateField(column.name, value),
          },
        };
      })
    ));

    return {
      components,
    };
  },
};
</script>

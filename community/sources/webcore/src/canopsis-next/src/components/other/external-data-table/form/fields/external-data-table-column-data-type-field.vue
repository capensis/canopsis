<template>
  <v-layout class="gap-2" column>
    <v-flex>
      <c-select-chip
        v-field="value.type"
        :items="types"
        class="px-2"
        color="grey"
        text-color="blue darken-1"
        outlined
      >
        <template #selection-empty>
          <span class="grey--text">Select data type</span>
        </template>
      </c-select-chip>
    </v-flex>
    <v-flex v-if="additionalComponent">
      <component
        :is="additionalComponent.is"
        v-bind="additionalComponent.bind"
        v-on="additionalComponent.on"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { mapValues } from 'lodash';
import { computed } from 'vue';

import { EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form';

import ExternalDataTableColumnDataTypeFieldNumberForm from './partials/external-data-table-column-data-type-field-number-form.vue';
import ExternalDataTableColumnDataTypeFieldStringArrayForm
  from './partials/external-data-table-column-data-type-field-string-array-form.vue';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    tableSeparator: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t, te } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const typesMap = computed(() => mapValues(EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES, (type) => {
      const key = `externalData.tableColumnDataTypes.${type}`;
      const item = {
        value: type,
        text: t(`${key}.text`),
      };

      if (te(`${key}.tooltip`)) {
        item.tooltip = t(`${key}.tooltip`);
      }

      return item;
    }));

    const types = computed(() => Object.values(typesMap.value));
    const additionalComponent = computed(() => {
      if ([
        EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number,
        EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray,
      ].includes(props.value.type)) {
        return {
          is: props.value.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number
            ? ExternalDataTableColumnDataTypeFieldNumberForm
            : ExternalDataTableColumnDataTypeFieldStringArrayForm,
          bind: {
            value: props.value,
            tableSeparator: props.tableSeparator,
          },
          on: {
            input: value => updateModel(value),
          },
        };
      }

      return null;
    });

    return {
      typesMap,
      types,
      additionalComponent,
    };
  },
};
</script>

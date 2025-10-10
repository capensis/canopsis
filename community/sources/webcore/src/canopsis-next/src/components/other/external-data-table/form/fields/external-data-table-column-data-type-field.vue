<template>
  <v-layout column>
    <v-flex>
      <c-select-chip
        v-field="value.type"
        :items="types"
        :disabled="disabled"
        class="px-2"
        color="grey"
        text-color="blue darken-1"
        outlined
      >
        <template #selection-empty>
          <span class="grey--text">{{ $t('externalData.selectDataType') }}</span>
        </template>
      </c-select-chip>
    </v-flex>
    <v-flex v-if="isNumber || isStringArray">
      <external-data-table-column-data-type-field-number-form
        v-if="isNumber"
        v-field="value"
        :table-separator="tableSeparator"
        :disabled="disabled"
      />
      <external-data-table-column-data-type-field-string-array-menu
        v-if="isStringArray && !disabled"
        v-field="value"
        :table-separator="tableSeparator"
        :disabled="disabled"
      >
        <template #selection="{ on }">
          <c-chip
            :color="stringArrayChip.color"
            :text-color="stringArrayChip.textColor"
            class="mt-2 px-2"
            outlined
            v-on="on"
          >
            <span v-if="!value.string_array_type" :class="stringArrayChip.prefixColorClass">
              {{ $t('externalData.tableColumnDataTypesAdditionalChips.stringArray.selectSeparator') }}
            </span>
            <span v-else-if="isCustomValueTypeSelected">
              <span :class="stringArrayChip.prefixColorClass" class="mr-2">
                {{ stringArrayChip.prefix }}:
              </span>
              <span>
                {{ value.string_array_separator }}
              </span>
            </span>
            <span v-else>
              <span :class="stringArrayChip.prefixColorClass">
                {{ stringArrayChip.prefix }}
              </span>
            </span>
          </c-chip>
        </template>
      </external-data-table-column-data-type-field-string-array-menu>
    </v-flex>
  </v-layout>
</template>

<script>
import { mapValues } from 'lodash';
import { computed, onBeforeMount, onBeforeUnmount } from 'vue';
import { Validator } from 'vee-validate';

import {
  EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

import ExternalDataTableColumnDataTypeFieldNumberForm from './partials/external-data-table-column-data-type-field-number-form.vue';
import ExternalDataTableColumnDataTypeFieldStringArrayMenu
  from './partials/external-data-table-column-data-type-field-string-array-menu.vue';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  components: {
    ExternalDataTableColumnDataTypeFieldNumberForm,
    ExternalDataTableColumnDataTypeFieldStringArrayMenu,
  },
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
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledTypes: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const { t, te } = useI18n();

    const validator = useValidator();

    const typesMap = computed(() => mapValues(EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES, (type) => {
      const key = `externalData.tableColumnDataTypes.${type}`;
      const item = {
        value: type,
        text: t(`${key}.text`),
      };

      if (te(`${key}.tooltip`)) {
        item.tooltip = t(`${key}.tooltip`);
      }

      if (props.disabledTypes.includes(type)) {
        item.disabled = true;
        item.disabledMessage = t(`${key}.disabledMessage`);
      }

      return item;
    }));

    const types = computed(() => Object.values(typesMap.value));
    const isNumber = computed(() => props.value.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number);
    const isStringArray = computed(() => props.value.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray);

    const isCustomValueTypeSelected = computed(() => (
      props.value.string_array_type === EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom
    ));

    const stringArrayChip = computed(() => {
      const result = {
        prefix: isCustomValueTypeSelected.value
          ? t('externalData.tableColumnDataTypesAdditionalChips.stringArray.separator')
          : t(`externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${props.value.string_array_type}.text`),
        color: 'grey',
        textColor: 'blue darken-1',
        prefixColorClass: 'grey--text',
      };

      if (validator.errors.has(props.value.name)) {
        result.color = 'error';
        result.textColor = 'error';
        result.prefixColorClass = 'error--text';
      }

      return result;
    });

    const {
      attachRequiredRule,
      detachRequiredRule,
    } = useValidationAttachRequired(props.value.name);

    onBeforeMount(() => attachRequiredRule(() => !isStringArray.value || !!props.value.string_array_type));
    onBeforeUnmount(detachRequiredRule);

    return {
      typesMap,
      types,
      isNumber,
      isStringArray,
      isCustomValueTypeSelected,
      stringArrayChip,
    };
  },
};
</script>

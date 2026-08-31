<template>
  <c-layout column>
    <c-select-chip
      :value="value.type"
      :items="types"
      :disabled="disabled"
      class="px-2"
      color="grey"
      text-color="blue darken-1"
      outlined
      @update:value="updateField('type', $event)"
    >
      <template #selection-empty>
        <span class="grey--text">
          {{ $t('externalData.selectDataType') }}
        </span>
      </template>
    </c-select-chip>
    <template v-if="isNumber || isStringArray">
      <external-data-table-column-data-type-field-number-form
        v-if="isNumber"
        :value="value"
        :table-separator="tableSeparator"
        :disabled="disabled"
        @update:value="updateModel"
      />
      <external-data-table-column-data-type-field-string-array-menu
        v-if="isStringArray && !disabled"
        :value="value"
        :table-separator="tableSeparator"
        :disabled="disabled"
        @update:value="updateModel"
      >
        <template #selection="{ props }">
          <span class="d-inline-flex align-self-start">
            <c-chip
              v-bind="props"
              :color="stringArrayChip.color"
              :text-color="stringArrayChip.textColor"
              class="mt-2 px-2"
              outlined
            >
            <span v-if="!value.string_array_type" :class="stringArrayChip.prefixColorClass">
              {{
                $t('externalData.tableColumnDataTypesAdditionalChips.stringArray.selectSeparator')
              }}
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
          </span>
        </template>
      </external-data-table-column-data-type-field-string-array-menu>
    </template>
  </c-layout>
</template>

<script>
import { mapValues } from 'lodash';
import { computed, onBeforeMount, onBeforeUnmount } from 'vue';

import {
  EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES,
  EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES,
} from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';
import { useValidator } from '@/hooks/validator/validator';

import ExternalDataTableColumnDataTypeFieldNumberForm from './partials/external-data-table-column-data-type-field-number-form.vue';
import ExternalDataTableColumnDataTypeFieldStringArrayMenu from './partials/external-data-table-column-data-type-field-string-array-menu.vue';

export default {
  emits: ['update:value'],
  components: {
    ExternalDataTableColumnDataTypeFieldNumberForm,
    ExternalDataTableColumnDataTypeFieldStringArrayMenu,
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
  },
  setup(props, { emit }) {
    const { t, te } = useI18n();
    const validator = useValidator();
    const { updateField, updateModel } = useModelField(props, emit);

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
      }),
    );

    const types = computed(() => Object.values(typesMap.value));
    const isNumber = computed(
      () => props.value.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.number,
    );
    const isStringArray = computed(
      () => props.value.type === EXTERNAL_DATA_TABLE_COLUMN_DATA_TYPES.stringArray,
    );

    const isCustomValueTypeSelected = computed(
      () => props.value.string_array_type ===
        EXTERNAL_DATA_TABLE_COLUMN_STRING_ARRAY_DATA_TYPE_TYPES.custom,
    );

    const stringArrayChip = computed(() => {
      const result = {
        prefix: '',
        color: 'grey',
        textColor: 'blue darken-1',
        prefixColorClass: 'grey--text',
      };

      if (isCustomValueTypeSelected.value) {
        result.prefix = t('externalData.tableColumnDataTypesAdditionalChips.stringArray.separator');
      } else if (props.value.string_array_type) {
        result.prefix = t(
          `externalData.tableColumnDataTypesAdditionalChips.stringArray.types.${props.value.string_array_type}.text`,
        );
      }

      if (validator.errors.has(props.value.name)) {
        result.color = 'error';
        result.textColor = 'error';
        result.prefixColorClass = 'error--text';
      }

      return result;
    });

    const { attachRequiredRule, detachRequiredRule } = useValidationAttachRequired(
      props.value.name,
    );

    onBeforeMount(() => attachRequiredRule(() => !isStringArray.value || !!props.value.string_array_type),
    );
    onBeforeUnmount(detachRequiredRule);

    return {
      typesMap,
      types,
      isNumber,
      isStringArray,
      isCustomValueTypeSelected,
      stringArrayChip,
      updateField,
      updateModel,
    };
  },
};
</script>

<template>
  <v-layout class="gap-2" column>
    <v-flex v-for="chip in chips" :key="chip.key">
      <c-select-chip
        v-bind="chip.bind"
        class="px-2"
        color="grey"
        text-color="blue darken-1"
        outlined
        v-on="chip.on"
      >
        <template #selection-empty>
          <span class="grey--text">{{ chip.selectionEmpty }}</span>
        </template>
        <template #selection-prefix>
          <span class="grey--text">{{ chip.selectionPrefix }}:</span>
        </template>
      </c-select-chip>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_DECIMAL_SEPARATOR,
  EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_THOUSANDS_SEPARATOR,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

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
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const getDisabledMessage = (separator, competingSeparator, messageKey) => {
      if (separator === props.tableSeparator) {
        return t('externalData.tableColumnDataTypesAdditionalChips.number.separatorDisabledByTableSeparator');
      }

      if (separator === competingSeparator) {
        return t(`externalData.tableColumnDataTypesAdditionalChips.number.${messageKey}`);
      }

      return null;
    };

    const decimalSeparators = computed(() => [
      ...Object.values(EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_DECIMAL_SEPARATOR).map((separator) => {
        const disabledMessage = getDisabledMessage(separator, props.value.thousands_separator, 'decimalSeparatorDisabled');

        return {
          disabledMessage,
          disabled: !!disabledMessage,
          value: separator,
          text: separator,
        };
      }),
    ]);

    const thousandsSeparators = computed(() => [
      ...Object.values(EXTERNAL_DATA_TABLE_COLUMN_NUMBER_DATA_TYPE_THOUSANDS_SEPARATOR).map((separator) => {
        const disabledMessage = getDisabledMessage(separator, props.value.decimal_separator, 'thousandsSeparatorDisabled');

        return {
          disabledMessage,
          disabled: !!disabledMessage,
          value: separator,
          text: separator,
        };
      }),
    ]);

    const chips = computed(() => [
      {
        key: 'decimal_separator',
        selectionEmpty: t('externalData.tableColumnDataTypesAdditionalChips.number.selectDecimalSeparator'),
        selectionPrefix: t('externalData.tableColumnDataTypesAdditionalChips.number.decimalSeparator'),
        bind: { value: props.value.decimal_separator, items: decimalSeparators.value },
        on: { input: separator => updateField('decimal_separator', separator) },
      },
      {
        key: 'thousands_separator',
        selectionEmpty: t('externalData.tableColumnDataTypesAdditionalChips.number.selectThousandsSeparator'),
        selectionPrefix: t('externalData.tableColumnDataTypesAdditionalChips.number.thousandsSeparator'),
        bind: { value: props.value.thousands_separator, items: thousandsSeparators.value },
        on: { input: separator => updateField('thousands_separator', separator) },
      },
    ]);

    return {
      chips,
    };
  },
};
</script>

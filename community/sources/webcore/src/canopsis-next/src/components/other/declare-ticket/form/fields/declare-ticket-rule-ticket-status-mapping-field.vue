<template>
  <c-form-block-row :label="$t('declareTicket.ticketStatusMapping')" :depth="1" indented>
    <v-layout class="gap-3" column>
      <c-label
        :label="$t('declareTicket.ticketStatusMapping')"
        :help-text="$t('declareTicket.ticketStatusMappingHelpText')"
        :error="errors.has(name)"
        required
      />

      <c-text-pairs-field
        v-field="value"
        :name="name"
        :disabled="disabled"
        :value-items="canopsisValueItems"
        :text-label="$t('declareTicket.sourceValue')"
        :value-label="$t('declareTicket.canopsisValue')"
        :add-button-label="$t('declareTicket.addMappingPair')"
        :required-error-message="$t('declareTicket.ticketStatusMappingHelpText')"
        text-required
        value-required
      />
    </v-layout>
  </c-form-block-row>
</template>

<script>
import { computed, watch } from 'vue';

import { DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES_WITHOUT_UNKNOWN } from '@/constants';

import { textPairToForm } from '@/helpers/text-pairs';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'status_mapping',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const canopsisValueItems = computed(() => (
      Object.values(DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES_WITHOUT_UNKNOWN).map(value => ({
        value,
        text: t(`declareTicket.status.${value}`),
      }))
    ));

    /**
     * Appends an empty mapping pair (text/value) to the ticket status mapping array.
     */
    const addMappingPair = () => addItemIntoArray(textPairToForm());

    /**
     * Checks whether at least one mapping pair has the "closed" Canopsis status.
     *
     * @returns {boolean}
     */
    const hasClosedItem = () => (
      props.value.some(item => item.value === DECLARE_TICKET_RULE_STATUS_MAPPING_VALUES_WITHOUT_UNKNOWN.closed)
    );

    const { asyncValidateRequiredRule } = useValidationAttachRequiredForField(props.name, hasClosedItem, false);

    watch(() => props.value, (newValue = [], oldValue = []) => {
      if (newValue.length && newValue.length !== oldValue.length) {
        return;
      }

      asyncValidateRequiredRule();
    });

    return {
      canopsisValueItems,
      addMappingPair,
      removeItemFromArray,
    };
  },
};
</script>

<template>
  <c-information-block
    :title="$t('declareTicket.ticketStatusMapping')"
    :help-text="$t('declareTicket.ticketStatusMappingHelpText')"
    help-icon="help"
    help-icon-color="grey darken-1"
  >
    <v-layout class="gap-3" column>
      <strong class="mt-3 grey--text text--darken-1">
        {{ $t('declareTicket.ticketStatusMappingUnmappedToOpen') }}
      </strong>

      <c-alert :value="errors.has(name)" type="error">
        {{ $t('declareTicket.ticketStatusMappingHelpText') }}
      </c-alert>

      <c-text-pairs-field
        v-field="value"
        :name="name"
        :disabled="disabled"
        :value-items="canopsisValueItems"
        :text-label="$t('declareTicket.sourceValue')"
        :value-label="$t('declareTicket.canopsisValue')"
        :add-button-label="$t('declareTicket.addMappingPair')"
        text-required
        value-required
      />
    </v-layout>
  </c-information-block>
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

    watch(() => props.value, asyncValidateRequiredRule);

    return {
      canopsisValueItems,
      addMappingPair,
      removeItemFromArray,
    };
  },
};
</script>

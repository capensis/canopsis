<template>
  <v-chip
    v-bind="chip.bind"
    v-on="chip.on"
  >
    <v-progress-circular
      v-if="instructionsPending && chip.hasPending"
      color="primary"
      size="20"
      width="2"
      indeterminate
    />
    <span v-else v-html="chip.text" />
  </v-chip>
</template>

<script>
import { omit } from 'lodash';
import { computed } from 'vue';

import { REMEDIATION_INSTRUCTION_FILTER_FIELDS } from '@/constants';

import { getSelectionText } from '@/helpers/vuetify';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'filter',
    event: 'input',
  },
  props: {
    field: {
      type: String,
      default: '',
    },
    filter: {
      type: Object,
      default: () => ({}),
    },
    instructions: {
      type: Array,
      default: () => [],
    },
    instructionsPending: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    /**
     * Clears all remediation instruction filter fields by updating the model to an empty object.
     */
    const clearFilterType = () => updateModel({});

    /**
     * Clears instruction_type, instruction_statuses, and instruction_ids from the filter.
     */
    const clearType = () => updateModel(omit(props.filter, ['instruction_type', 'instruction_statuses', 'instruction_ids']));

    /**
     * Clears instruction_statuses from the filter.
     */
    const clearStatuses = () => updateModel(omit(props.filter, ['instruction_statuses']));

    /**
     * Clears instruction_ids from the filter.
     */
    const clearIds = () => updateModel(omit(props.filter, ['instruction_ids']));

    const chip = computed(() => {
      switch (props.field) {
        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionFilterType:
          return {
            bind: { close: true },
            on: { 'click:close': clearFilterType },
            text: t(`alarm.instructionsFilter.filters.${props.filter.instruction_filter_type}`),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionType:
          return {
            bind: { close: true },
            on: { 'click:close': clearType },
            text: t(`remediation.instruction.types.${props.filter.instruction_type}`),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionStatuses:
          return {
            bind: { close: true },
            on: { 'click:close': clearStatuses },
            text: t('alarm.instructionsFilter.statusSelector', {
              statuses: props.filter.instruction_statuses.map(status => t(`alarm.instructionsFilter.statuses.${status}`)).join(', '),
            }),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionIds:
          return {
            hasPending: true,
            bind: { close: true },
            on: { 'click:close': clearIds },
            text: t('alarm.instructionsFilter.nameSelector', {
              names: props.filter.instruction_ids.map(id => getSelectionText(props.instructions, id, '_id', 'name')).join(', '),
            }),
          };

        default:
          return {};
      }
    });

    return {
      chip,
    };
  },
};
</script>

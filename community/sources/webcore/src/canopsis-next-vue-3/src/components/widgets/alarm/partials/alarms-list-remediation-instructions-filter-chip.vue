<template>
  <v-chip
    :closable="chip.closable"
    size="small"
    @click:close="chip.close"
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

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    field: {
      type: String,
      default: '',
    },
    modelValue: {
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

  emits: ['update:modelValue'],

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
    const clearType = () => updateModel(omit(props.modelValue, ['instruction_type', 'instruction_statuses', 'instruction_ids']));

    /**
     * Clears instruction_statuses from the filter.
     */
    const clearStatuses = () => updateModel(omit(props.modelValue, ['instruction_statuses']));

    /**
     * Clears instruction_ids from the filter.
     */
    const clearIds = () => updateModel(omit(props.modelValue, ['instruction_ids']));

    const chip = computed(() => {
      switch (props.field) {
        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionFilterType:
          return {
            closable: true,
            close: clearFilterType,
            text: t(`alarm.instructionsFilter.filters.${props.modelValue.instruction_filter_type}`),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionType:
          return {
            closable: true,
            close: clearType,
            text: t(`remediation.instruction.types.${props.modelValue.instruction_type}`),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionStatuses:
          return {
            closable: true,
            close: clearStatuses,
            text: t('alarm.instructionsFilter.statusSelector', {
              statuses: props.modelValue.instruction_statuses
                .map(status => t(`alarm.instructionsFilter.statuses.${status}`))
                .join(', '),
            }),
          };

        case REMEDIATION_INSTRUCTION_FILTER_FIELDS.instructionIds:
          return {
            hasPending: true,
            closable: true,
            close: clearIds,
            text: t('alarm.instructionsFilter.nameSelector', {
              names: props.modelValue.instruction_ids
                .map(id => getSelectionText(props.instructions, id, '_id', 'name'))
                .join(', '),
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

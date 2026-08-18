<template>
  <div class="alarms-list-remediation-instructions-filter-fields">
    <c-select
      v-if="field === fieldsMap.instructionFilterType"
      :model-value="modelValue.instruction_filter_type"
      :items="instructionFilterTypes"
      :label="$t('alarm.instructionsFilter.filter')"
      item-title="text"
      item-value="value"
      clearable
      @update:model-value="updateInstructionFilterType"
    />
    <c-select
      v-else-if="field === fieldsMap.instructionType"
      :model-value="modelValue.instruction_type"
      :items="instructionTypes"
      :label="$t('alarm.instructionsFilter.type')"
      item-title="text"
      item-value="value"
      clearable
      @update:model-value="updateInstructionType"
    />
    <c-select
      v-else-if="field === fieldsMap.instructionStatuses"
      :model-value="modelValue.instruction_statuses"
      :items="instructionStatuses"
      :label="$t('alarm.instructionsFilter.status')"
      item-title="text"
      item-value="value"
      multiple
      chips
      clearable
      @update:model-value="updateInstructionStatuses"
    />
    <c-lazy-search-field
      v-else-if="field === fieldsMap.instructionIds"
      :value="modelValue.instruction_ids ?? []"
      :items="instructions"
      :label="$t('alarm.instructionsFilter.name')"
      :loading="instructionsPending"
      :has-more="hasMoreInstructions"
      item-text="name"
      item-value="_id"
      multiple
      clearable
      autocomplete
      @fetch="$emit('fetch:instructions')"
      @fetch:more="$emit('fetch:more-instructions')"
      @update:search="$emit('update:instructions-search', $event)"
      @update:value="updateInstructionIds"
    >
      <template #selection="{ item, props }">
        <v-chip v-bind="props" size="small" closable>
          {{ getInstructionNameById(item) }}
        </v-chip>
      </template>
    </c-lazy-search-field>
  </div>
</template>

<script>
import { isNil, omit } from 'lodash';
import { computed } from 'vue';

import {
  REMEDIATION_INSTRUCTION_FILTER_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
  REMEDIATION_INSTRUCTION_FILTER_MANUAL_STATUSES,
  REMEDIATION_INSTRUCTION_FILTER_AUTO_STATUSES,
  REMEDIATION_INSTRUCTION_FILTER_FIELDS,
} from '@/constants';

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
    hasMoreInstructions: {
      type: Boolean,
      default: false,
    },
  },

  emits: [
    'update:modelValue',
    'fetch:instructions',
    'fetch:more-instructions',
    'update:instructions-search',
  ],

  setup(props, { emit }) {
    const fieldsMap = REMEDIATION_INSTRUCTION_FILTER_FIELDS;

    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const instructionFilterTypes = computed(() =>
      Object.values(REMEDIATION_INSTRUCTION_FILTER_TYPES).map(type => ({
        value: type,
        text: t(`alarm.instructionsFilter.filters.${type}`),
      })),
    );

    const instructionTypes = computed(() =>
      [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.auto].map(type => ({
        value: type,
        text: t(`remediation.instruction.types.${type}`),
      })),
    );

    const manualInstructionStatuses = computed(() =>
      REMEDIATION_INSTRUCTION_FILTER_MANUAL_STATUSES.map(status => ({
        value: status,
        text: t(`alarm.instructionsFilter.statuses.${status}`),
      })),
    );

    const autoInstructionStatuses = computed(() =>
      REMEDIATION_INSTRUCTION_FILTER_AUTO_STATUSES.map(status => ({
        value: status,
        text: t(`alarm.instructionsFilter.statuses.${status}`),
      })),
    );

    const instructionStatuses = computed(() =>
      props.modelValue.instruction_type === REMEDIATION_INSTRUCTION_TYPES.manual
        ? manualInstructionStatuses.value
        : autoInstructionStatuses.value,
    );

    /**
     * Returns the instruction name for a given instruction ID from the instructions array.
     *
     * @param {string|Object} instruction - The instruction ID or object to look up.
     * @returns {string} The name of the instruction, or an empty string if not found.
     */
    const getInstructionNameById = instruction =>
      getSelectionText(props.instructions, instruction, '_id', 'name');

    /**
     * Updates the filter with the selected instruction filter type. Clears dependent fields if filter type is unset.
     *
     * @param {string|null} instructionFilterType - The selected instruction filter type, or null to clear.
     */
    const updateInstructionFilterType = instructionFilterType =>
      updateModel(
        isNil(instructionFilterType)
          ? {}
          : {
            ...omit(props.modelValue, ['instruction_type', 'instruction_statuses', 'instruction_ids']),
            instruction_filter_type: instructionFilterType,
          },
      );

    /**
     * Updates the filter with the selected instruction type. Clears dependent fields if type is unset.
     *
     * @param {string|null} instructionType - The selected instruction type, or null to clear.
     */
    const updateInstructionType = instructionType =>
      updateModel(
        isNil(instructionType)
          ? omit(props.modelValue, ['instruction_type', 'instruction_statuses', 'instruction_ids'])
          : {
            ...omit(props.modelValue, ['instruction_statuses', 'instruction_ids']),
            instruction_type: instructionType,
          },
      );

    const updateInstructionStatuses = instructionStatusesValue =>
      updateModel({
        ...props.modelValue,
        instruction_statuses: instructionStatusesValue ?? [],
      });

    const updateInstructionIds = instructionIds =>
      updateModel({
        ...props.modelValue,
        instruction_ids: instructionIds ?? [],
      });

    return {
      fieldsMap,
      instructionFilterTypes,
      instructionTypes,
      instructionStatuses,
      getInstructionNameById,
      updateInstructionFilterType,
      updateInstructionType,
      updateInstructionStatuses,
      updateInstructionIds,
    };
  },
};
</script>

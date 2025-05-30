<template>
  <v-list-item-title>
    <c-select-field
      v-if="field === fieldsMap.instructionFilterType"
      :value="filter.instruction_filter_type"
      :items="instructionFilterTypes"
      :label="$t('alarm.instructionsFilter.filter')"
      clearable
      @input="updateInstructionFilterType"
    />
    <c-select-field
      v-else-if="field === fieldsMap.instructionType"
      :value="filter.instruction_type"
      :items="instructionTypes"
      :label="$t('alarm.instructionsFilter.type')"
      clearable
      @input="updateInstructionType"
    />
    <c-select-field
      v-else-if="field === fieldsMap.instructionStatuses"
      v-field="filter.instruction_statuses"
      :items="instructionStatuses"
      :label="$t('alarm.instructionsFilter.status')"
      multiple
      chips
      clearable
      deletable-chips
    />
    <c-select-field
      v-else-if="field === fieldsMap.instructionIds"
      v-field="filter.instruction_ids"
      :items="instructions"
      :label="$t('alarm.instructionsFilter.name')"
      :loading="instructionsPending"
      :return-object="false"
      item-text="name"
      item-value="_id"
      multiple
      chips
      clearable
      combobox
      deletable-chips
    >
      <template #selection="{ item, parent }">
        <v-chip close @click:close="parent.onChipInput(item)">
          {{ getInstructionNameById(item) }}
        </v-chip>
      </template>
    </c-select-field>
  </v-list-item-title>
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
    const fieldsMap = REMEDIATION_INSTRUCTION_FILTER_FIELDS;

    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const instructionFilterTypes = computed(() => Object.values(REMEDIATION_INSTRUCTION_FILTER_TYPES).map(type => ({
      value: type,
      text: t(`alarm.instructionsFilter.filters.${type}`),
    })));

    const instructionTypes = computed(() => [
      REMEDIATION_INSTRUCTION_TYPES.manual,
      REMEDIATION_INSTRUCTION_TYPES.auto,
    ].map(type => ({
      value: type,
      text: t(`remediation.instruction.types.${type}`),
    })));

    const manualInstructionStatuses = computed(() => REMEDIATION_INSTRUCTION_FILTER_MANUAL_STATUSES.map(status => ({
      value: status,
      text: t(`alarm.instructionsFilter.statuses.${status}`),
    })));

    const autoInstructionStatuses = computed(() => REMEDIATION_INSTRUCTION_FILTER_AUTO_STATUSES.map(status => ({
      value: status,
      text: t(`alarm.instructionsFilter.statuses.${status}`),
    })));

    const instructionStatuses = computed(() => (
      props.filter.instruction_type === REMEDIATION_INSTRUCTION_TYPES.manual
        ? manualInstructionStatuses.value
        : autoInstructionStatuses.value
    ));

    const hasInstructionFilterType = computed(() => (
      props.filter.instruction_filter_type === REMEDIATION_INSTRUCTION_FILTER_TYPES.hasInstructions
    ));

    const hasInstructionType = computed(() => !isNil(props.filter.instruction_type));

    /**
     * Returns the instruction name for a given instruction ID from the instructions array.
     *
     * @param {string} instructionId - The ID of the instruction to look up.
     * @returns {string} The name of the instruction, or an empty string if not found.
     */
    const getInstructionNameById = instructionId => getSelectionText(props.instructions, instructionId, '_id', 'name');

    /**
     * Updates the filter with the selected instruction filter type. Clears dependent fields if flter type is unset.
     *
     * @param {string|null} instructionFilterType - The selected instruction filter type, or null to clear.
     */
    const updateInstructionFilterType = instructionFilterType => updateModel(
      isNil(instructionFilterType)
        ? {}
        : {
          ...omit(props.filter, ['instruction_type', 'instruction_statuses', 'instruction_ids']),
          instruction_filter_type: instructionFilterType,
        },
    );

    /**
     * Updates the filter with the selected instruction type. Clears dependent fields if type is unset.
     *
     * @param {string|null} instructionType - The selected instruction type, or null to clear.
     */
    const updateInstructionType = instructionType => updateModel(
      isNil(instructionType)
        ? omit(props.filter, ['instruction_type', 'instruction_statuses', 'instruction_ids'])
        : { ...omit(props.filter, ['instruction_statuses', 'instruction_ids']), instruction_type: instructionType },
    );

    return {
      fieldsMap,
      instructionFilterTypes,
      instructionTypes,
      instructionStatuses,
      hasInstructionFilterType,
      hasInstructionType,
      getInstructionNameById,
      updateInstructionFilterType,
      updateInstructionType,
    };
  },
};
</script>

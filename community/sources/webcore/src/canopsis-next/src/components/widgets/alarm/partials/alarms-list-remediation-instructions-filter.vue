<template>
  <c-select-field
    :value="selectedItems"
    :loading="instructionsPending"
    :label="$t('alarm.instructionsFilter.filter')"
    :menu-props="menuProps"
    :items="fields"
    :hide-input="selectedAllFields"
    multiple
    chips
    clearable
    hide-details
    @input="updateFilter"
  >
    <template #selection="{ item }">
      <alarms-list-remediation-instructions-filter-chip
        v-field="filter"
        :field="item"
        :instructions="instructions"
        :instructions-pending="instructionsPending"
      />
    </template>
    <template #item="{ item }">
      <v-list-item>
        <alarms-list-remediation-instructions-filter-fields
          v-field="filter"
          :field="item"
          :instructions="instructions"
          :instructions-pending="instructionsPending"
        />
      </v-list-item>
    </template>
  </c-select-field>
</template>

<script>
import { computed, toRef } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import {
  useAlarmsListRemediationInstructionsFilterFetch,
  useAlarmsListRemediationInstructionsFilterFields,
} from './hooks/alarms-list-remediation-instructions-filter';
import AlarmsListRemediationInstructionsFilterChip from './alarms-list-remediation-instructions-filter-chip.vue';
import AlarmsListRemediationInstructionsFilterFields from './alarms-list-remediation-instructions-filter-fields.vue';

export default {
  components: {
    AlarmsListRemediationInstructionsFilterChip,
    AlarmsListRemediationInstructionsFilterFields,
  },
  model: {
    prop: 'filter',
    event: 'input',
  },
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const { instructions, instructionsPending } = useAlarmsListRemediationInstructionsFilterFetch(toRef(props, 'filter'));
    const { fields, selectedItems, selectedAllFields } = useAlarmsListRemediationInstructionsFilterFields(toRef(props, 'filter'));

    /**
     * Handles filter updates from the select field.
     * If the provided value is empty or falsy, resets the filter by emitting an empty object.
     *
     * @param {Array|undefined} value - The new value from the select field. If empty, filter is reset.
     */
    const updateFilter = (value) => {
      if (value?.length) {
        return;
      }

      updateModel({});
    };

    const menuProps = computed(() => ({
      maxHeight: '350px',
    }));

    return {
      instructions,
      instructionsPending,
      selectedItems,
      selectedAllFields,
      updateFilter,
      fields,
      menuProps,
    };
  },
};
</script>

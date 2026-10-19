<template>
  <c-select
    :model-value="selectedItems"
    :label="$t('alarm.instructionsFilter.filter')"
    :menu-props="menuProps"
    :items="fields"
    class="alarms-list-remediation-instructions-filter"
    multiple
    chips
    clearable
    hide-no-data
    hide-details
    @update:model-value="updateFilter"
  >
    <template #item="{ item }">
      <div class="alarms-list-remediation-instructions-filter__field px-4 py-2" @click.stop>
        <alarms-list-remediation-instructions-filter-fields
          :model-value="modelValue"
          :field="item"
          :instructions="instructions"
          :instructions-pending="instructionsPending"
          :has-more-instructions="hasMoreInstructions"
          @fetch:instructions="fetchInstructionsList"
          @fetch:more-instructions="fetchMoreInstructions"
          @update:instructions-search="updateInstructionsSearch"
          @update:model-value="updateModel"
        />
      </div>
    </template>
    <template #chip="{ item }">
      <alarms-list-remediation-instructions-filter-chip
        :model-value="modelValue"
        :field="item"
        :instructions="instructions"
        :instructions-pending="initialInstructionsPending"
        @update:model-value="updateModel"
      />
    </template>
  </c-select>
</template>

<script>
import { computed, toRef } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

import AlarmsListRemediationInstructionsFilterChip from './alarms-list-remediation-instructions-filter-chip.vue';
import AlarmsListRemediationInstructionsFilterFields from './alarms-list-remediation-instructions-filter-fields.vue';
import {
  useAlarmsListRemediationInstructionsFilterFetch,
  useAlarmsListRemediationInstructionsFilterFields,
} from './hooks/alarms-list-remediation-instructions-filter';

export default {
  components: {
    AlarmsListRemediationInstructionsFilterChip,
    AlarmsListRemediationInstructionsFilterFields,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);
    const {
      instructions,
      instructionsPending,
      initialInstructionsPending,
      hasMoreInstructions,
      fetchInstructionsList,
      fetchMoreInstructions,
      updateInstructionsSearch,
    } =
      useAlarmsListRemediationInstructionsFilterFetch(toRef(props, 'modelValue'));
    const { fields, selectedItems, selectedAllFields } =
      useAlarmsListRemediationInstructionsFilterFields(toRef(props, 'modelValue'));

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
      initialInstructionsPending,
      hasMoreInstructions,
      selectedItems,
      selectedAllFields,
      fetchInstructionsList,
      fetchMoreInstructions,
      updateInstructionsSearch,
      updateFilter,
      updateModel,
      fields,
      menuProps,
    };
  },
};
</script>

<style lang="scss" scoped>
.alarms-list-remediation-instructions-filter__field {
  min-width: 320px;
}
</style>

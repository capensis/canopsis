<template>
  <v-layout class="gap-4">
    <v-flex md3 xs6>
      <c-select-field
        :value="options.last_run_status"
        :label="$t('jobs.filterByStatus')"
        :items="statusItems"
        item-value="value"
        item-text="text"
        clearable
        hide-details
        @input="updateLastRunStatus"
      />
    </v-flex>
    <v-flex md3 xs6>
      <c-select-field
        :value="options.status"
        :label="$t('jobs.filterByActiveState')"
        :items="activeStateItems"
        item-value="value"
        item-text="text"
        clearable
        hide-details
        @input="updateStatus"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { JOB_STATUS, JOB_LAST_RUN_STATUS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    options: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const statusItems = computed(() => Object.values(JOB_LAST_RUN_STATUS).map(value => ({
      value,
      text: t(`jobs.lastRunStatus.${value}`),
    })));

    const activeStateItems = computed(() => Object.values(JOB_STATUS).map(value => ({
      value,
      text: t(`jobs.status.${value}`),
    })));

    /**
     * Emits updated options to parent.
     *
     * @param {Object} newOptions - Updated table options
     */
    const updateOptions = newOptions => emit('update:options', newOptions);

    /**
     * Updates status filter and resets pagination.
     *
     * @param {number|undefined} status - Run status value
     */
    const updateLastRunStatus = lastRunStatus => updateOptions({
      ...props.options,

      last_run_status: lastRunStatus ?? undefined,
      page: 1,
    });

    /**
     * Updates active state filter and resets pagination.
     *
     * @param {number|undefined} activeState - Job state value
     */
    const updateStatus = status => updateOptions({
      ...props.options,

      status: status ?? undefined,
      page: 1,
    });

    return {
      statusItems,
      activeStateItems,
      updateLastRunStatus,
      updateStatus,
    };
  },
};
</script>

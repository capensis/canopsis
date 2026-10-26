<template>
  <v-row>
    <v-col cols="12" sm="6" md="3">
      <c-select
        :model-value="options.status"
        :label="t('jobs.filterByActiveState')"
        :items="activeStateItems"
        item-title="text"
        item-value="value"
        clearable
        hide-details
        @update:model-value="updateStatus"
      />
    </v-col>
    <v-col cols="12" sm="6" md="3">
      <c-select
        :model-value="options.last_run_status"
        :label="t('jobs.filterByLastStatus')"
        :items="statusItems"
        item-title="text"
        item-value="value"
        clearable
        hide-details
        @update:model-value="updateLastRunStatus"
      />
    </v-col>
  </v-row>
</template>

<script>
import { computed } from 'vue';

import { JOB_LAST_RUN_STATUS, JOB_STATUS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    options: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const statusItems = computed(() => Object.values(JOB_LAST_RUN_STATUS).map(value => ({
        value,
        text: t(`jobs.lastRunStatus.${value}`),
      })),
    );
    const activeStateItems = computed(() => Object.values(JOB_STATUS).map(value => ({
        value,
        text: t(`jobs.status.${value}`),
      })),
    );
    const updateOptions = newOptions => emit('update:options', newOptions);
    const updateLastRunStatus = lastRunStatus => updateOptions({
        ...props.options,
        last_run_status: lastRunStatus ?? undefined,
        page: 1,
      });
    const updateStatus = status => updateOptions({
        ...props.options,
        status: status ?? undefined,
        page: 1,
      });

    return {
      t,
      statusItems,
      activeStateItems,
      updateLastRunStatus,
      updateStatus,
    };
  },
};
</script>

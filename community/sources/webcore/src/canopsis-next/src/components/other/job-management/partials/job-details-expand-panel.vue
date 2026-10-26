<template>
  <div>
    <v-tabs
      v-model="activeTab"
      color="primary"
      class="secondary lighten-1"
      align-tabs="center"
    >
      <v-tab :value="TABS.data">
        {{ t('common.data') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <c-window-expand-item :value="TABS.data" cols="8">
        <v-row density="compact">
          <v-col v-if="!isTicketStatusJob || item.fail_reason" cols="12">
            <strong>{{ statusPrefix }}</strong>
            <template v-if="item.fail_reason">: {{ item.fail_reason }}</template>
          </v-col>

          <v-col v-if="item.raw_request" cols="12">
            <job-http-data-collapse-panel
              :title-prefix="t('jobs.data.request')"
              :raw="item.raw_request"
              :expanded="!isTicketStatusJob && isRunning"
            />
          </v-col>

          <v-col v-if="(isTicketStatusJob || !isRunning) && item.raw_response" cols="12">
            <job-http-data-collapse-panel
              :title-prefix="t('jobs.data.response')"
              :raw="item.raw_response"
              expanded
              tree
              copyable
            />
          </v-col>
        </v-row>
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { computed, ref } from 'vue';

import { JOB_EXECUTION_STATUS, JOBS_TABS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import JobHttpDataCollapsePanel from './job-http-data-collapse-panel.vue';

const TABS = {
  data: 'data',
};

export default {
  components: { JobHttpDataCollapsePanel },

  props: {
    item: {
      type: Object,
      required: true,
    },
    jobType: {
      type: String,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const activeTab = ref(TABS.data);
    const isTicketStatusJob = computed(() => props.jobType === JOBS_TABS.ticketStatus);
    const isRunning = computed(() => props.item.status === JOB_EXECUTION_STATUS.running);
    const statusPrefix = computed(() => isTicketStatusJob.value
      ? t('jobs.data.webhookFailedPrefix')
      : t(`jobs.data.statusPrefix.${props.jobType}.${props.item.status}`));

    return {
      t,
      TABS,
      activeTab,
      isTicketStatusJob,
      isRunning,
      statusPrefix,
    };
  },
};
</script>

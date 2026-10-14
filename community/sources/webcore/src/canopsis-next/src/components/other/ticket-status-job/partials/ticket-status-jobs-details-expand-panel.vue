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
        <c-layout gap="3" column>
          <span v-if="item.fail_reason">
            <strong>{{ t('jobs.data.webhookFailedPrefix') }} : </strong>
            {{ item.fail_reason }}
          </span>
          <ticket-status-jobs-http-data-collapse-panel
            v-if="item.raw_request"
            :title-prefix="t('jobs.data.request')"
            :raw="item.raw_request"
          />
          <ticket-status-jobs-http-data-collapse-panel
            v-if="item.raw_response"
            :title-prefix="t('jobs.data.response')"
            :raw="item.raw_response"
            expanded
          />
        </c-layout>
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobsHttpDataCollapsePanel from './ticket-status-jobs-http-data-collapse-panel.vue';

const TABS = {
  data: 'data',
};

export default {
  components: {
    TicketStatusJobsHttpDataCollapsePanel,
  },

  props: {
    item: {
      type: Object,
      required: true,
    },
  },

  setup() {
    const { t } = useI18n();
    const activeTab = ref(TABS.data);

    return {
      t,
      TABS,
      activeTab,
    };
  },
};
</script>

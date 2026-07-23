<template>
  <div>
    <v-tabs v-model="activeTab" color="primary" class="secondary lighten-1" align-tabs="center">
      <v-tab :value="TABS.data">
        {{ t('common.data') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item :value="TABS.data">
        <div class="pa-3">
          <div class="ticket-status-jobs-details-expand-panel__content">
            <v-card>
              <v-card-text>
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
              </v-card-text>
            </v-card>
          </div>
        </div>
      </v-window-item>
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

<style lang="scss" scoped>
.ticket-status-jobs-details-expand-panel__content {
  max-width: 66.666%;
  margin: 0 auto;
}
</style>

<template>
  <view-tab-widgets
    v-if="activeTab"
    :tab="activeTab"
    kiosk
    visible
  />
</template>

<script>
import Observer from '@/services/observer';

import { toSeconds } from '@/helpers/date/duration';

import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';

import ViewTabWidgets from '@/components/other/view/view-tab-widgets.vue';

export default {
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  components: {
    ViewTabWidgets,
  },
  mixins: [
    queryMixin,
    activeViewMixin,
  ],
  props: {
    id: {
      type: String,
      required: true,
    },
    tabId: {
      type: String,
      required: true,
    },
  },
  computed: {
    activeTab() {
      const { tabId } = this.$route.params;

      return this.view?.tabs?.find(tab => tab._id === tabId);
    },

    periodicRefreshEnabled() {
      return this.view?.periodic_refresh?.enabled;
    },

    periodicRefreshSeconds() {
      const { value, unit } = this.view?.periodic_refresh ?? {};

      return toSeconds(value, unit);
    },
  },
  watch: {
    periodicRefreshEnabled(enabled) {
      if (enabled) {
        this.startPeriodicRefresh();
      } else {
        this.stopPeriodicRefresh();
      }
    },

    periodicRefreshSeconds() {
      this.stopPeriodicRefresh();

      if (this.periodicRefreshEnabled) {
        this.startPeriodicRefresh();
      }
    },
  },

  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },

  created() {
    this.clearActiveView();

    this.$periodicRefresh.register(this.refreshView);
  },

  async mounted() {
    await this.fetchActiveView({ id: this.id });

    if (!this.activeTab) {
      this.$router.replace({
        params: {
          id: this.id,
          tabId: this.view.tabs[0]._id,
        },
      });
    }
  },

  beforeDestroy() {
    this.$periodicRefresh.unregister(this.refreshView);
  },

  methods: {
    async refreshView() {
      await this.fetchActiveView({ id: this.id });

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    refresh() {
      return this.$periodicRefresh.notify();
    },

    startPeriodicRefresh() {
      this.periodicRefreshTimer = setInterval(this.refresh, this.periodicRefreshSeconds * 1000);
    },

    stopPeriodicRefresh() {
      clearInterval(this.periodicRefreshTimer);
    },
  },
};
</script>

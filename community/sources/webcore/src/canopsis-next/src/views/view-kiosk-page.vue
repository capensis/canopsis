<template>
  <div class="view-kiosk-page">
    <c-progress-overlay :pending="activeViewPending" />

    <view-tab-widgets v-if="activeTab" :tab="activeTab" visible />
  </div>
</template>

<script>
import { ROUTES_NAMES, VIEW_SCREEN_MODES } from '@/constants';

import Observer from '@/services/observer';

import { toSeconds } from '@/helpers/date/duration';

import { activeViewMixin } from '@/mixins/active-view';
import { queryMixin } from '@/mixins/query';

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

  mixins: [activeViewMixin, queryMixin],

  props: {
    id: {
      type: [String, Number],
      required: true,
    },
    tabId: {
      type: [String, Number],
      required: true,
    },
  },

  data() {
    return {
      periodicRefreshTimer: null,
    };
  },

  computed: {
    activeTab() {
      if (!this.view?.tabs?.length) {
        return null;
      }

      return this.view.tabs.find(tab => tab._id === this.tabId) ?? null;
    },

    periodicRefreshEnabled() {
      return this.view?.periodic_refresh?.enabled ?? false;
    },

    periodicRefreshSeconds() {
      const { value, unit } = this.view?.periodic_refresh ?? {};

      return value ? toSeconds(value, unit) : 0;
    },
  },

  watch: {
    async id() {
      this.clearActiveView();
      this.setActiveViewScreenMode(VIEW_SCREEN_MODES.kiosk);

      await this.initializeView();
    },

    tabId: 'ensureActiveTabRoute',

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
    this.setActiveViewScreenMode(VIEW_SCREEN_MODES.kiosk);
    this.$periodicRefresh.register(this.refreshView);
  },

  async mounted() {
    await this.initializeView();
  },

  beforeUnmount() {
    this.stopPeriodicRefresh();
    this.$periodicRefresh.unregisterAll();
    this.clearActiveView();
  },

  methods: {
    async initializeView() {
      try {
        await this.fetchActiveView({ id: this.id });
        await this.ensureActiveTabRoute();

        if (this.periodicRefreshEnabled) {
          this.startPeriodicRefresh();
        }
      } catch (err) {
        console.error(err);

        await this.$router.replace({ name: ROUTES_NAMES.home });
      }
    },

    async ensureActiveTabRoute() {
      if (!this.view?.tabs?.length) {
        return;
      }

      if (!this.activeTab) {
        await this.$router.replace({
          name: ROUTES_NAMES.viewKiosk,
          params: {
            id: this.id,
            tabId: this.view.tabs[0]._id,
          },
        });
      }
    },

    async refreshView() {
      await this.fetchActiveView({ id: this.id });
      await this.ensureActiveTabRoute();

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    refresh() {
      return this.$periodicRefresh.notify();
    },

    startPeriodicRefresh() {
      if (!this.periodicRefreshSeconds) {
        return;
      }

      this.stopPeriodicRefresh();

      this.periodicRefreshTimer = setInterval(this.refresh, this.periodicRefreshSeconds * 1000);
    },

    stopPeriodicRefresh() {
      clearInterval(this.periodicRefreshTimer);
      this.periodicRefreshTimer = null;
    },
  },
};
</script>

<style lang="scss" scoped>
.view-kiosk-page {
  min-height: 100%;
}
</style>

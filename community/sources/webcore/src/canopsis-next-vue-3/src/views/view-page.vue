<template>
  <div class="view-wrapper position-relative">
    <c-progress-overlay :pending="activeViewPending" />

    <v-fade-transition>
      <view-tabs-wrapper
        v-if="isViewTabsReady"
        :updatable="hasUpdateAccess"
      />
    </v-fade-transition>

    <view-fab-btns
      v-if="isViewTabsReady"
      :active-tab="activeTab"
      :updatable="hasUpdateAccess"
      :sharable="!view.is_private && hasCreateAnyShareTokenAccess"
    />
  </div>
</template>

<script>
import { ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import Observer from '@/services/observer';

import { authMixin } from '@/mixins/auth';
import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';

import ViewFabBtns from '@/components/other/view/partials/view-fab-btns.vue';
import ViewTabsWrapper from '@/components/other/view/view-tabs-wrapper.vue';

export default {
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  components: {
    ViewFabBtns,
    ViewTabsWrapper,
  },
  mixins: [
    authMixin,
    activeViewMixin,
    viewRouterMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  computed: {
    hasUpdateAccess() {
      return this.view?.is_private || this.checkUpdateAccess(this.id);
    },

    hasCreateAnyShareTokenAccess() {
      return this.checkCreateAccess(USER_PERMISSIONS.technical.shareToken);
    },

    activeTab() {
      const { tabId } = this.$route.query;

      if (!this.view?.tabs?.length) {
        return null;
      }

      if (!tabId) {
        return this.view.tabs[0];
      }

      return this.view.tabs.find(tab => tab._id === tabId) ?? null;
    },

    isViewTabsReady() {
      return this.view?.tabs && this.$route.query.tabId && this.activeTab;
    },
  },
  watch: {
    async id() {
      this.clearActiveView();

      await this.initializeView();
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
    await this.initializeView();
  },
  beforeUnmount() {
    this.$periodicRefresh.unregisterAll();
  },
  methods: {
    async initializeView() {
      try {
        await this.fetchActiveView({ id: this.id });
        await this.ensureActiveTabRoute();
      } catch {
        await this.$router.replace({ name: ROUTES_NAMES.home });
      }
    },

    async ensureActiveTabRoute() {
      const { tabId } = this.$route.query;

      if (!this.view?.tabs?.length) {
        if (tabId) {
          await this.redirectToViewRoot();
        }

        return;
      }

      if (!tabId || !this.activeTab) {
        await this.redirectToFirstTab();
      }
    },

    async refreshView() {
      await this.fetchActiveView({ id: this.id });
      await this.ensureActiveTabRoute();
    },
  },
};
</script>

<style lang="scss" scoped>
.view-wrapper {
  padding-bottom: 70px;
}
</style>

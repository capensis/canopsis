<template lang="pug">
  div.view-wrapper
    v-fade-transition
      view-tabs-wrapper(
        v-if="isViewTabsReady",
        :editing="editing",
        :updatable="hasUpdateAccess"
      )
    v-fade-transition
      view-fab-btns(v-if="view", :active-tab="activeTab", :updatable="hasUpdateAccess")
</template>

<script>
import Observer from '@/services/observer';

import ViewTabsWrapper from '@/components/other/view/view-tabs-wrapper.vue';
import ViewFabBtns from '@/components/other/view/buttons/view-fab-btns.vue';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';
import { viewRouterMixin } from '@/mixins/view/router';

export default {
  provide() {
    return {
      $periodicRefresh: this.$periodicRefresh,
    };
  },
  components: {
    ViewTabsWrapper,
    ViewFabBtns,
  },
  mixins: [
    authMixin,
    queryMixin,
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
      return this.checkUpdateAccess(this.id);
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
      return this.view?.tabs && this.$route.query.tabId;
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
    const { tabId } = this.$route.query;

    await this.fetchActiveView({ id: this.id });

    if (!tabId) {
      await this.redirectToFirstTab();
    } else if (!this.activeTab) {
      await this.redirectToViewRoot();
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
  },
};
</script>

<style lang="scss" scoped>
  .view-wrapper {
    padding-bottom: 70px;
  }
</style>

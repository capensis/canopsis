<template lang="pug">
  div.view-wrapper
    v-fade-transition
      view-tabs-wrapper(
        v-if="isViewTabsReady",
        :editing="editing",
        :updatable="hasUpdateAccess"
      )
    view-fab-btns(:active-tab="activeTab", :updatable="hasUpdateAccess")
</template>

<script>
import Observer from '@/services/observer';

import ViewTabsWrapper from '@/components/other/view/view-tabs-wrapper.vue';
import ViewFabBtns from '@/components/other/view/buttons/view-fab-btns.vue';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';
import { activeViewMixin } from '@/mixins/active-view';

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

      if (this.view?.tabs?.length) {
        if (!tabId) {
          return this.view.tabs[0];
        }

        return this.view.tabs.find(tab => tab._id === tabId) ?? null;
      }

      return null;
    },

    isViewTabsReady() {
      return this.view?.tabs && this.$route.query.tabId;
    },
  },

  beforeCreate() {
    this.$periodicRefresh = new Observer();
  },

  created() {
    this.registerViewOnceWatcher();
    this.$periodicRefresh.register(this.refreshView);
  },

  mounted() {
    this.fetchActiveView({ id: this.id });
  },

  beforeDestroy() {
    this.$periodicRefresh.unregister(this.refreshView);
    this.clearActiveView();
  },

  methods: {
    async refreshView() {
      await this.fetchActiveView({ id: this.id });

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    registerViewOnceWatcher() {
      const unwatch = this.$watch('view', (view) => {
        if (view) {
          const { tabId } = this.$route.query;

          if (!tabId && view.tabs && view.tabs.length) {
            this.$router.replace({ query: { tabId: view.tabs[0]._id } });
          }

          unwatch();
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .refresh-btn {
    text-decoration: none;
    text-transform: none;
  }

  .view-wrapper {
    padding-bottom: 70px;
  }
</style>

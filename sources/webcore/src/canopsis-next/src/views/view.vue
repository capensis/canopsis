<template lang="pug">
  div(:data-test="`view-page-${id}`")
    v-fade-transition
      view-tabs-wrapper(
        v-if="isViewTabsReady",
        :view="view",
        :isEditingMode="isEditingMode",
        :hasUpdateAccess="hasUpdateAccess",
        :updateViewMethod="data => updateView({ id, data })"
      )
    .fab
      v-layout(data-test="controlViewLayout", column)
        v-tooltip(left)
          v-btn(
            slot="activator",
            :input-value="isPeriodicRefreshEnabled",
            color="secondary",
            fab,
            dark,
            @click.stop="refreshViewWithProgress"
          )
            v-icon(v-if="!isPeriodicRefreshEnabled") refresh
            v-progress-circular.periodic-refresh-progress(
              v-else,
              :rotate="270",
              :size="30",
              :width="2",
              :value="periodicBehaviorProgressValue",
              color="white",
              button
            )
              span {{ periodicRefreshProgress }}
          span {{ $t('common.refresh') }}
        v-speed-dial(
          v-if="hasUpdateAccess",
          v-model="isVSpeedDialOpen",
          direction="left",
          transition="slide-y-reverse-transition"
        )
          v-btn(
            data-test="menuViewButton",
            slot="activator",
            :input-value="isVSpeedDialOpen",
            color="primary",
            dark,
            fab
          )
            v-icon menu
            v-icon close
          v-tooltip(top)
            v-btn(
              slot="activator",
              v-model="isFullScreenMode",
              fab,
              dark,
              small,
              @click="toggleFullScreenMode"
            )
              v-icon fullscreen
              v-icon fullscreen_exit
            span alt + enter / command + enter
          v-tooltip(v-if="hasUpdateAccess", top)
            v-btn(
              data-test="editViewButton",
              slot="activator",
              fab,
              dark,
              small,
              @click.stop="toggleViewEditingMode",
              v-model="isEditingMode"
            )
              v-icon edit
              v-icon done
            span {{ $t('common.toggleEditView') }}  (ctrl + e / command + e)
          v-tooltip(top)
            v-btn(
              data-test="addWidgetButton",
              v-if="hasUpdateAccess",
              slot="activator",
              fab,
              dark,
              small,
              color="indigo",
              @click.stop="showCreateWidgetModal"
            )
              v-icon add
            span {{ $t('common.addWidget') }}
          v-tooltip(top)
            v-btn(
              data-test="addTabButton",
              v-if="hasUpdateAccess",
              slot="activator",
              fab,
              dark,
              small,
              color="green",
              @click.stop="showCreateTabModal"
            )
              v-icon add
            span {{ $t('common.addTab') }}
        v-tooltip(v-else, left)
          v-btn(
            slot="activator",
            v-model="isFullScreenMode",
            fab,
            dark,
            @click="toggleFullScreenMode"
          )
            v-icon fullscreen
            v-icon fullscreen_exit
          div {{ $t('view.fullScreen') }}
            .font-italic.caption.ml-1 ({{ $t('view.fullScreenShortcut') }})
</template>

<script>
import { get } from 'lodash';
import { MODALS, USERS_RIGHTS_MASKS } from '@/constants';
import { generateViewTab } from '@/helpers/entities';

import ViewTabRows from '@/components/other/view/view-tab-rows.vue';
import ViewTabsWrapper from '@/components/other/view/view-tabs-wrapper.vue';

import authMixin from '@/mixins/auth';
import queryMixin from '@/mixins/query';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    ViewTabRows,
    ViewTabsWrapper,
  },
  mixins: [
    authMixin,
    queryMixin,
    entitiesViewMixin,
  ],
  props: {
    id: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      isEditingMode: false,
      isFullScreenMode: false,
      isVSpeedDialOpen: false,
      periodicRefreshInterval: null,
      periodicRefreshProgress: null,
    };
  },
  computed: {
    periodicBehaviorProgressValue() {
      return this.periodicRefreshProgress / (this.periodicRefreshValue / 100);
    },

    isPeriodicRefreshEnabled() {
      return get(this.view, 'periodicRefresh.enabled', false);
    },

    periodicRefreshValue() {
      return get(this.view, 'periodicRefresh.value', 0);
    },

    hasUpdateAccess() {
      return this.checkUpdateAccess(this.id, USERS_RIGHTS_MASKS.update);
    },

    activeTab() {
      const { tabId } = this.$route.query;

      if (this.view.tabs && this.view.tabs.length) {
        if (!tabId) {
          return this.view.tabs[0];
        }

        return this.view.tabs.find(tab => tab._id === tabId) || null;
      }

      return null;
    },

    isViewTabsReady() {
      return this.view && this.$route.query.tabId;
    },
  },

  watch: {
    isPeriodicRefreshEnabled(value, oldValue) {
      if (value && (!oldValue || !this.periodicRefreshInterval)) {
        this.startPeriodicRefreshInterval();
      } else if (oldValue && !value) {
        this.stopPeriodicRefreshInterval();
      }
    },
  },

  created() {
    document.addEventListener('keydown', this.keyDownListener);
    this.registerViewOnceWatcher();
  },

  mounted() {
    this.fetchView({ id: this.id });

    if (this.isPeriodicRefreshEnabled) {
      this.startPeriodicRefreshInterval();
    }
  },

  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);

    this.stopPeriodicRefreshInterval();
  },

  methods: {
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

    keyDownListener(event) {
      if (event.key === 'Enter' && event.altKey) {
        this.toggleFullScreenMode();
        event.preventDefault();
      } else if (event.key === 'e' && event.ctrlKey) {
        this.toggleViewEditingMode();
        event.preventDefault();
      }
    },

    toggleFullScreenMode() {
      if (this.activeTab) {
        const element = document.getElementById(`view-tab-${this.activeTab._id}`);

        if (element) {
          this.$fullscreen.toggle(element, {
            fullscreenClass: 'full-screen',
            background: 'white',
            callback: value => this.isFullScreenMode = value,
          });
        }
      } else {
        this.$popups.warning({ text: this.$t('view.errors.emptyTabs') });
      }
    },

    async refreshView() {
      await this.fetchView({ id: this.id });

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    async refreshViewWithProgress() {
      this.stopPeriodicRefreshInterval();

      await this.refreshView();

      this.startPeriodicRefreshInterval();
    },

    showCreateWidgetModal() {
      if (this.activeTab) {
        this.$modals.show({
          name: MODALS.createWidget,
          config: {
            tabId: this.activeTab._id,
          },
        });
      } else {
        this.$popups.warning({ text: this.$t('view.errors.emptyTabs') });
      }
    },

    showCreateTabModal() {
      this.$modals.show({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.create.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            validationRules: 'required',
          },
          action: (title) => {
            const oldTabs = this.view.tabs || [];
            const newTab = { ...generateViewTab(), title };
            const view = {
              ...this.view,
              tabs: [...oldTabs, newTab],
            };

            return this.updateView({ id: this.id, data: view });
          },
        },
      });
    },

    toggleViewEditingMode() {
      this.isEditingMode = !this.isEditingMode;
    },

    resetRefreshInterval() {
      this.periodicRefreshProgress = this.periodicRefreshValue;
    },

    refreshTick() {
      if (this.periodicRefreshProgress <= 0) {
        this.refreshViewWithProgress();
      } else {
        this.periodicRefreshProgress -= 1;
      }
    },

    startPeriodicRefreshInterval() {
      this.resetRefreshInterval();

      if (this.periodicRefreshInterval) {
        this.stopPeriodicRefreshInterval();
      }

      this.periodicRefreshInterval = setInterval(this.refreshTick, 1000);
    },

    stopPeriodicRefreshInterval() {
      clearInterval(this.periodicRefreshInterval);
    },
  },
};
</script>

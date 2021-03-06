<template lang="pug">
  div.view-wrapper(:data-test="`view-page-${id}`")
    v-fade-transition
      view-tabs-wrapper(
        v-if="isViewTabsReady",
        :view="view",
        :isEditingMode="isEditingMode",
        :hasUpdateAccess="hasUpdateAccess",
        :updateViewMethod="updateViewMethod",
        @update:widgetsFields="updateWidgetsFieldsForUpdateById"
      )
    .fab
      v-layout(data-test="controlViewLayout", row)
        v-tooltip(top)
          v-btn(
            slot="activator",
            :input-value="isPeriodicRefreshEnabled",
            color="secondary",
            fab,
            dark,
            @click.stop="refreshHandler"
          )
            v-icon(v-if="!isPeriodicRefreshEnabled") refresh
            v-progress-circular.periodic-refresh-progress(
              v-else,
              :rotate="270",
              :size="30",
              :width="2",
              :value="periodicRefreshProgressValue",
              color="white",
              button
            )
              span.refresh-btn {{ periodicRefreshProgress | maxDurationByUnit }}
          span {{ tooltipContent }}
        v-speed-dial(
          v-if="hasUpdateAccess",
          v-model="isVSpeedDialOpen",
          direction="top",
          transition="slide-y-reverse-transition"
        )
          v-btn(
            slot="activator",
            :input-value="isVSpeedDialOpen",
            data-test="menuViewButton",
            color="primary",
            dark,
            fab
          )
            v-icon menu
            v-icon close
          v-tooltip(left)
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
          v-tooltip(v-if="hasUpdateAccess", left)
            v-btn(
              slot="activator",
              :input-value="isEditingMode",
              data-test="editViewButton",
              fab,
              dark,
              small,
              @click.stop="toggleViewEditingMode"
            )
              v-icon edit
              v-icon done
            div
              div {{ $t('common.toggleEditView') }}  (ctrl + e / command + e)
              div.font-italic {{ $t('common.toggleEditViewSubtitle') }}
          v-tooltip(left)
            v-btn(
              slot="activator",
              v-if="hasUpdateAccess",
              data-test="addWidgetButton",
              fab,
              dark,
              small,
              color="indigo",
              @click.stop="showCreateWidgetModal"
            )
              v-icon add
            span {{ $t('common.addWidget') }}
          v-tooltip(left)
            v-btn(
              slot="activator",
              v-if="hasUpdateAccess",
              data-test="addTabButton",
              fab,
              dark,
              small,
              color="green",
              @click.stop="showCreateTabModal"
            )
              v-icon add
            span {{ $t('common.addTab') }}
        v-tooltip(v-else, top)
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
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';
import { generateViewTab } from '@/helpers/entities';
import { setSeveralFields } from '@/helpers/immutable';

import ViewTabsWrapper from '@/components/other/view/view-tabs-wrapper.vue';

import authMixin from '@/mixins/auth';
import queryMixin from '@/mixins/query';
import entitiesViewMixin from '@/mixins/entities/view';
import periodicRefreshMixin from '@/mixins/view/periodic-refresh';

export default {
  components: {
    ViewTabsWrapper,
  },
  mixins: [
    authMixin,
    queryMixin,
    entitiesViewMixin,
    periodicRefreshMixin,
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
      widgetsFieldsForUpdateById: {},
    };
  },
  computed: {
    tooltipContent() {
      return this.isPeriodicRefreshEnabled ? this.periodicRefreshProgressFormatted : this.$t('common.refresh');
    },

    hasUpdateAccess() {
      return this.checkUpdateAccess(this.id);
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

  created() {
    document.addEventListener('keydown', this.keyDownListener);
    this.registerViewOnceWatcher();
    this.$periodicRefresh.subscribe(this.refreshView);
  },

  mounted() {
    this.fetchView({ id: this.id });
  },

  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);
    this.$periodicRefresh.unsubscribe(this.refreshView);
  },

  methods: {
    updateViewMethod(data) {
      return this.updateView({ id: this.id, data });
    },

    async refreshView() {
      await this.fetchView({ id: this.id });

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

    keyDownListener(event) {
      if (event.key === 'Enter' && event.altKey) {
        this.toggleFullScreenMode();
        event.preventDefault();
      } else if (event.key === 'e' && event.ctrlKey && this.hasUpdateAccess) {
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

    updateWidgetsFieldsForUpdateById(widgetsFieldsForUpdateById) {
      this.widgetsFieldsForUpdateById = {
        ...this.widgetsFieldsForUpdateById,
        ...widgetsFieldsForUpdateById,
      };
    },

    updateTabs() {
      const view = {
        ...this.view,

        tabs: this.view.tabs.map(tab => ({
          ...tab,

          widgets: tab.widgets.map((widget) => {
            const fields = this.widgetsFieldsForUpdateById[widget._id];

            if (fields) {
              return setSeveralFields(widget, fields);
            }

            return widget;
          }),
        })),
      };

      return this.updateView({ id: this.id, data: view });
    },

    async toggleViewEditingMode() {
      if (this.isEditingMode && !isEmpty(this.widgetsFieldsForUpdateById)) {
        await this.updateTabs();
      }

      this.isEditingMode = !this.isEditingMode;
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

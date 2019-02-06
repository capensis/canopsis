<template lang="pug">
  v-tabs.view-tabs(
  ref="tabs",
  :key="vTabsKey",
  :value="value",
  :class="{ hidden: this.tabs.length < 2 && !isEditingMode, 'tabs-editing': isEditingMode }",
  :hide-slider="isTabsChanged",
  color="secondary lighten-2",
  slider-color="primary",
  dark,
  @change="$emit('input', $event)"
  )
    draggable.d-flex(
    :value="tabs",
    :options="draggableOptions",
    @end="onDragEnd",
    @input="$emit('update:tabs', $event)"
    )
      v-tab.draggable-item(v-if="tabs.length", v-for="tab in tabs", :key="tab._id", :disabled="isTabsChanged", ripple)
        span {{ tab.title }}
        v-btn(
        v-show="hasUpdateAccess && isEditingMode",
        small,
        flat,
        icon,
        @click.stop="showUpdateTabModal(tab)"
        )
          v-icon(small) edit
        v-btn(
        v-show="hasUpdateAccess && isEditingMode",
        small,
        flat,
        icon,
        @click.stop="showDuplicateTabModal(tab)"
        )
          v-icon(small) file_copy
        v-btn(
        v-show="hasUpdateAccess && isEditingMode",
        small,
        flat,
        icon,
        @click.stop="showDeleteTabModal(tab)"
        )
          v-icon(small) delete
    v-tabs-items(v-if="$scopedSlots.default", active-class="active-view-tab")
      v-tab-item(v-for="tab in tabs", :key="tab._id", lazy)
        slot(
        :tab="tab",
        :isEditingMode="isEditingMode",
        :hasUpdateAccess="hasUpdateAccess",
        :updateTabMethod="updateTab"
        )
</template>

<script>
import { omit } from 'lodash';
import Draggable from 'vuedraggable';
import { createNamespacedHelpers } from 'vuex';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS } from '@/constants';

import {
  generateViewTab,
  generateViewRow,
  generateWidgetByType,
  generateUserPreferenceByWidgetAndUser,
} from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

const { mapActions: userPreferenceMapActions } = createNamespacedHelpers('userPreference');

export default {
  components: { Draggable },
  mixins: [
    authMixin,
    modalMixin,
    vuetifyTabsMixin,
  ],
  props: {
    view: {
      type: Object,
      required: true,
    },
    tabs: {
      type: Array,
      required: true,
    },
    value: {
      type: Number,
      default: null,
    },
    hasUpdateAccess: {
      type: Boolean,
      default: false,
    },
    isTabsChanged: {
      type: Boolean,
      default: false,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateViewMethod: {
      type: Function,
      default: () => {},
    },
  },
  computed: {
    vTabsKey() {
      return this.view.tabs.map(tab => tab._id).join('-');
    },
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.isEditingMode,
      };
    },
  },
  watch: {
    isEditingMode() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
    tabs: {
      immediate: true,
      handler() {
        this.onUpdateTabs();
      },
    },
  },
  methods: {
    ...userPreferenceMapActions({
      createUserPreference: 'create',
      fetchUserPreferenceByWidgetIdWithoutStore: 'fetchItemByWidgetIdWithoutStore',
    }),

    showUpdateTabModal(tab) {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.edit.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            value: tab.title,
            validationRules: 'required',
          },
          action: (title) => {
            const newTab = { ...tab, title };

            return this.updateTab(newTab);
          },
        },
      });
    },

    showDuplicateTabModal(tab) {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.duplicate.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            validationRules: 'required',
          },
          action: title => this.duplicateTabAction(tab, title),
        },
      });
    },

    showDeleteTabModal(tab) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.deleteTab(tab._id),
        },
      });
    },

    async duplicateTabAction(tab, title) {
      const widgetsIdsMap = {};
      const newTab = {
        ...generateViewTab(),

        title,
        rows: tab.rows.map(row => ({
          ...generateViewRow(),

          title: row.title,
          widgets: row.widgets.map((widget) => {
            const newWidget = generateWidgetByType(widget.type);

            widgetsIdsMap[widget._id] = newWidget._id;

            return {
              ...newWidget,
              ...omit(widget, ['_id']),
            };
          }),
        })),
      };

      await this.copyUserPreferencesForWidgets(widgetsIdsMap);

      return this.addTab(newTab);
    },

    async copyUserPreferencesForWidgets(widgetsIdsMap) {
      const oldWidgetsIds = Object.keys(widgetsIdsMap);
      const userPreferences = await Promise.all(oldWidgetsIds.map(widgetId =>
        this.fetchUserPreferenceByWidgetIdWithoutStore({ widgetId })));


      return Promise.all(userPreferences.map((userPreference) => {
        if (!userPreference) {
          return Promise.resolve();
        }

        const newWidgetId = widgetsIdsMap[userPreference.widget_id];
        const newUserPreference = generateUserPreferenceByWidgetAndUser({ _id: newWidgetId }, this.currentUser);

        return this.createUserPreference({
          userPreference: {
            ...newUserPreference,
            ...omit(userPreference, ['_id', 'widget_id']),
          },
        });
      }));
    },

    updateTab(tab) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.map((viewTab) => {
          if (viewTab._id === tab._id) {
            return tab;
          }

          return viewTab;
        }),
      };

      return this.updateViewMethod(view);
    },

    addTab(tab) {
      const view = {
        ...this.view,
        tabs: [...this.view.tabs, tab],
      };

      return this.updateViewMethod(view);
    },

    deleteTab(tabId) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.filter(viewTab => viewTab._id !== tabId),
      };

      return this.updateViewMethod(view);
    },

    onUpdateTabs() {
      this.$nextTick(() => {
        this.callTabsOnResizeMethod();
        this.callTabsUpdateTabsMethod();
      });
    },

    onDragEnd() {
      this.onUpdateTabs();
    },
  },
};
</script>

<style lang="scss" scoped>
  .view-tabs.hidden {
    & /deep/ .v-tabs__bar {
      display: none;
    }
  }

  .draggable-item {
    position: relative;
    transform: translateZ(0);

    .tabs-editing & {
      cursor: move;

      & /deep/ .v-tabs__item {
        cursor: move;
      }
    }

    & /deep/ .v-tabs__item--disabled {
      color: #fff;
      opacity: 1;

      button {
        color: rgba(255,255,255,0.3) !important;
        box-shadow: none !important;
        pointer-events: none;
      }
    }
  }
</style>

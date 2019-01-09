<template lang="pug">
  div
    view-tabs(
    v-if="!viewPending || !view",
    v-model="activeTabIndex",
    :view="view",
    :isEditingMode="isEditingMode",
    :hasUpdateAccess="hasUpdateAccess",
    :updateViewMethod="data => updateView({ id, data })"
    )
      view-tab-rows(
      slot-scope="props",
      v-bind="props",
      )
    .fab
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, color="secondary", @click.stop="refreshView")
          v-icon refresh
        span {{ $t('common.refresh') }}
      v-speed-dial(
      v-model="isVSpeedDialOpen",
      direction="left",
      transition="slide-y-reverse-transition"
      )
        v-btn(slot="activator", :input-value="isVSpeedDialOpen", color="primary", dark, fab)
          v-icon menu
          v-icon close
        v-tooltip(top)
          v-btn(
          slot="activator",
          v-model="isFullScreenMode"
          fab,
          dark,
          small,
          @click="toggleFullScreenMode",
          )
            v-icon fullscreen
            v-icon fullscreen_exit
          span alt + enter / command + enter
        v-tooltip(v-if="hasUpdateAccess", top)
          v-btn(slot="activator", fab, dark, small, @click.stop="toggleViewEditingMode", v-model="isEditingMode")
            v-icon edit
            v-icon done
          span {{ $t('common.toggleEditView') }}  (ctrl + e / command + e)
        v-tooltip(top)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateWidgetModal")
            v-icon add
          span {{ $t('common.addWidget') }}
        v-tooltip(top)
          v-btn(slot="activator", fab, dark, small, color="green", @click.stop="showCreateTabModal")
            v-icon add
          span {{ $t('common.addTab') }}
</template>

<script>
import isNull from 'lodash/isNull';

import { MODALS, USERS_RIGHTS_MASKS } from '@/constants';
import { generateViewTab } from '@/helpers/entities';

import ViewTabs from '@/components/other/view/view-tabs.vue';
import ViewTabRows from '@/components/other/view/view-tab-rows.vue';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import popupMixin from '@/mixins/popup';
import queryMixin from '@/mixins/query';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    ViewTabs,
    ViewTabRows,
  },
  mixins: [
    authMixin,
    modalMixin,
    popupMixin,
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
      isVSpeedDialOpen: false,
      activeTabIndex: null,
      isEditingMode: false,
      isFullScreenMode: false,
    };
  },
  computed: {
    hasUpdateAccess() {
      return this.checkUpdateAccess(this.id, USERS_RIGHTS_MASKS.update);
    },

    activeTab() {
      if (this.view.tabs && this.view.tabs.length) {
        if (isNull(this.activeTabIndex)) {
          return this.view.tabs[0];
        }

        return this.view.tabs[this.activeTabIndex];
      }

      return null;
    },
  },
  created() {
    document.addEventListener('keydown', this.keyDownListener);
    this.fetchView({ id: this.id });
  },
  beforeDestroy() {
    this.$fullscreen.exit();
    document.removeEventListener('keydown', this.keyDownListener);
  },
  methods: {
    keyDownListener(event) {
      if (event.keyCode === 13 && event.altKey) { // alt + enter
        this.toggleFullScreenMode();
        event.preventDefault();
      } else if (event.keyCode === 69 && event.ctrlKey) { // ctrl + e
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
        this.addWarningPopup({ text: this.$t('view.errors.emptyTabs') });
      }
    },

    async refreshView() {
      await this.fetchView({ id: this.id });

      if (this.activeTab) {
        this.forceUpdateQuery({ id: this.activeTab._id });
      }
    },

    showCreateWidgetModal() {
      if (this.activeTab) {
        this.showModal({
          name: MODALS.createWidget,
          config: {
            tabId: this.activeTab._id,
          },
        });
      } else {
        this.addWarningPopup({ text: this.$t('view.errors.emptyTabs') });
      }
    },

    showCreateTabModal() {
      this.showModal({
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
  },
};
</script>

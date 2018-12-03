<template lang="pug">
  div#view-page
    view-tabs(
    v-model="activeTabIndex",
    :view="view",
    :isEditingMode="isEditingMode",
    :hasUpdateAccess="hasUpdateAccess",
    :updateViewMethod="view => updateView({ id, data: view })"
    )
      view-widgets(
      slot-scope="{ tab, updateTabMethod }",
      :tab="tab",
      :isEditingMode="isEditingMode",
      :hasUpdateAccess="hasUpdateAccess",
      :updateTabMethod="updateTabMethod"
      )
    .fab
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, color="secondary", @click.stop="refreshView")
          v-icon refresh
        span {{ $t('common.refresh') }}
      v-speed-dial(
      v-model="fab",
      direction="left",
      transition="slide-y-reverse-transition"
      )
        v-btn(slot="activator", color="primary", dark, fab, v-model="fab")
          v-icon menu
          v-icon close
        v-tooltip(top)
          v-btn(
          slot="activator",
          v-model="isFullScreenMode"
          fab,
          dark,
          small,
          @click="toggleFullScreen",
          )
            v-icon fullscreen
            v-icon fullscreen_exit
          span alt + enter / command + enter
        v-tooltip(v-if="hasUpdateAccess", top)
          v-btn(slot="activator", fab, dark, small, @click.stop="toggleViewEditingMode", v-model="isEditingMode")
            v-icon edit
            v-icon done
          span {{ $t('common.toggleEditView') }}
        v-tooltip(top)
          v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateWidgetModal")
            v-icon add
          span {{ $t('common.addWidget') }}
        v-tooltip(top)
          v-btn(slot="activator", fab, dark, small, color="green", @click.stop="showCreateTabModal")
            v-icon add
          span Add tab
</template>

<script>
import get from 'lodash/get';

import { MODALS, USERS_RIGHTS_MASKS } from '@/constants';
import uid from '@/helpers/uid';
import { generateViewTab } from '@/helpers/entities';

import ViewTabs from '@/components/other/view/view-tabs.vue';
import ViewWidgets from '@/components/other/view/view-widgets.vue';

import authMixin from '@/mixins/auth';
import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import entitiesViewMixin from '@/mixins/entities/view';

export default {
  components: {
    ViewTabs,
    ViewWidgets,
  },
  mixins: [
    authMixin,
    popupMixin,
    modalMixin,
    sideBarMixin,
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
      fab: false,
      activeTabIndex: null,
      isEditingMode: false,
      isFullScreenMode: false,
    };
  },
  computed: {
    hasUpdateAccess() {
      return this.checkUpdateAccess(this.id, USERS_RIGHTS_MASKS.update);
    },
    getWidgetFlexClass() {
      return widget => [
        `xs${widget.size.sm}`,
        `md${widget.size.md}`,
        `lg${widget.size.lg}`,
      ];
    },
    rows() {
      return get(this.view, 'rows', []);
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
    keyDownListener({ keyCode, altKey }) {
      if (keyCode === 13 && altKey) {
        this.toggleFullScreen();
      }
    },

    toggleFullScreen() {
      const element = document.getElementById('view');

      if (element) {
        this.$fullscreen.toggle(element, {
          fullscreenClass: 'full-screen',
          background: 'white',
          callback: value => this.isFullScreenMode = value,
        });
      }
    },

    async refreshView() {
      await this.fetchView({ id: this.id });

      // TODO: fix it

      this.widgetKeyPrefix = uid();
    },

    showCreateWidgetModal() {
      this.showModal({
        name: MODALS.createWidget,
        config: {
          tabId: this.view.tabs[this.activeTabIndex]._id,
        },
      });
    },

    showCreateTabModal() {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: 'Create tab',
          field: {
            name: 'text',
            label: 'Title',
            validationRules: 'required',
          },
          action: (title) => {
            const newTab = { ...generateViewTab(), title };
            const view = {
              ...this.view,
              tabs: [...this.view.tabs, newTab],
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

<style lang="scss">
  #view-page {
    .full-screen {
      .hide-on-full-screen {
        display: none;
      }
    }
  }
</style>

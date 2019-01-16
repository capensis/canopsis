<template lang="pug">
  v-tabs.view-tabs(
  ref="tabs",
  :value="value",
  :class="{ hidden: this.tabs.length < 2 }",
  color="secondary lighten-2",
  slider-color="primary",
  dark,
  @change="$emit('input', $event)"
  )
    draggable.d-flex(v-model="tabs", :options="draggableOptions", @end="onUpdateTabs")
      v-tab.draggable-item(v-if="tabs.length", v-for="tab in tabs", :key="`tab-${tab._id}`", ripple)
        span {{ tab.title }}
        v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showUpdateTabModal(tab)")
          v-icon(small) edit
        v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showDeleteTabModal(tab)")
          v-icon(small) delete
    v-tabs-items(ref="tabItems", active-class="active-view-tab")
      v-tab-item(v-for="tab in tabs", :key="`tab-item-${tab._id}`", lazy)
        slot(
        :tab="tab",
        :isEditingMode="isEditingMode",
        :hasUpdateAccess="hasUpdateAccess",
        :updateTabMethod="updateTab"
        )
</template>

<script>
import Draggable from 'vuedraggable';
import isEqual from 'lodash/isEqual';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

export default {
  components: { Draggable },
  mixins: [modalMixin, vuetifyTabsMixin],
  props: {
    view: {
      type: Object,
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
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    updateViewMethod: {
      type: Function,
      required: true,
    },
  },
  data() {
    return {
      tabs: [],
    };
  },
  computed: {
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        disabled: !this.isEditingMode,
      };
    },
  },
  watch: {
    isEditingMode(value) {
      this.$nextTick(this.callTabsOnResizeMethod);

      if (!value) {
        this.updateViewMethod({
          ...this.view,

          tabs: this.tabs,
        });
      }
    },
    'view.tabs': {
      immediate: true,
      handler(tabs, prevTabs) {
        if (!isEqual(tabs, prevTabs)) {
          this.tabs = [...tabs];
        }

        this.onUpdateTabs();
      },
    },
  },
  methods: {
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

    showDeleteTabModal(tab) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const view = {
              ...this.view,
              tabs: this.view.tabs.filter(viewTab => viewTab._id !== tab._id),
            };

            await this.updateViewMethod(view);
          },
        },
      });
    },

    updateTab(newTab) {
      const view = {
        ...this.view,
        tabs: this.view.tabs.map((viewTab) => {
          if (viewTab._id === newTab._id) {
            return newTab;
          }

          return viewTab;
        }),
      };

      return this.updateViewMethod(view);
    },

    onUpdateTabs() {
      this.$nextTick(() => {
        this.callTabsOnResizeMethod();
        this.callTabsUpdateTabsMethod();
      });
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
  }
</style>

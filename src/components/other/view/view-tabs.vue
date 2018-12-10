<template lang="pug">
  v-tabs.view-tabs(
  ref="tabs",
  :value="value",
  :class="{ hidden: this.tabs.length < 2 }",
  color="secondary lighten-2",
  slider-color="primary",
  dark,
  @input="$emit('input', $event)"
  )
    v-tab(v-if="tabs.length", v-for="tab in tabs", :key="`tab-${tab._id}`", ripple)
      span {{ tab.title }}
      v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showUpdateTabModal(tab)")
        v-icon(small) edit
      v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showDeleteTabModal(tab)")
        v-icon(small) delete
    v-tabs-items(ref="tabItems")
      v-tab-item(v-for="tab in tabs", :key="`tab-item-${tab._id}`", :value="tab._id", lazy)
        slot(
        :tab="tab",
        :isEditingMode="isEditingMode",
        :hasUpdateAccess="hasUpdateAccess",
        :updateTabMethod="updateTab"
        )
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

export default {
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
  computed: {
    tabs() {
      return this.view.tabs || [];
    },
  },
  watch: {
    isEditingMode() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
    'view.tabs': {
      handler() {
        this.$nextTick(() => {
          this.callTabsOnResizeMethod();
          this.callTabsUpdateTabsMethod();
        });
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
  },
};
</script>

<style lang="scss" scoped>
  .view-tabs.hidden {
    & /deep/ .v-tabs__bar {
      display: none;
    }
  }
</style>

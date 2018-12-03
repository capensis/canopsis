<template lang="pug">
  v-tabs(v-model="activeTab", color="secondary lighten-2", slider-color="primary", dark)
    v-tab(v-for="tab in view.tabs", :key="`tab-${tab._id}`", ripple)
      span {{ tab.title }}
      v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showUpdateTabModal(tab)")
        v-icon(small) edit
      v-btn(v-show="hasUpdateAccess && isEditingMode", small, flat, icon, @click.stop="showDeleteTabModal(tab)")
        v-icon(small) delete
    slot(v-for="tab in tabs", :tab="tab", @updateTab="updateTab")
</template>

<script>
import get from 'lodash/get';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';

export default {
  mixins: [modalMixin],
  props: {
    view: {
      type: Object,
      required: true,
    },
  },
  computed: {
    tabs() {
      return get(this.view, 'tabs', []);
    },
  },
  methods: {
    showUpdateTabModal(tab) {
      this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: 'Update tab',
          field: {
            name: 'text',
            label: 'Title',
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
          action: () => {
            const view = {
              ...this.view,
              tabs: this.view.tabs.filter(viewTab => viewTab._id !== tab._id),
            };

            return this.$emit('updateView', view);
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

      return this.$emit('updateView', view);
    },
  },
};
</script>

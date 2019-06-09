<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.view.select.title') }}
    v-card-text.view-select
      div Selected tab: {{ selected }}
      v-radio-group(v-model="selected")
        v-list.py-0
          v-list-group(v-for="group in availableGroups", :key="group._id")
            v-list-tile(slot="activator")
              v-list-tile-title {{ group.name }}
            v-list
              v-list-group(v-for="view in group.views", :key="view._id")
                v-list-tile(slot="activator")
                  v-list-tile-title.pl-2 {{ view.title }}
                v-list-tile(v-for="tab in view.tabs", :key="tab._id", @click="selectTab(view._id, tab._id)")
                  v-list-tile-action
                    v-radio(:value="{ viewId: view._id, tabId: tab._id }")
                  v-list-tile-title.pl-4 {{ tab.title }}
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, SIDE_BARS_BY_WIDGET_TYPES } from '@/constants';

import sideBarMixin from '@/mixins/side-bar/side-bar';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import rightsTechnicalViewMixin from '@/mixins/rights/technical/view';

const { mapGetters } = createNamespacedHelpers('entities');

export default {
  name: MODALS.cloneWidget,
  mixins: [
    sideBarMixin,
    modalInnerMixin,
    entitiesViewsGroupsMixin,
    rightsTechnicalViewMixin,
  ],
  data() {
    return {
      selected: null,
    };
  },
  computed: {
    ...mapGetters({
      getItem: 'getItem',
    }),

    selectedView() {
      return this.selected && this.getItem('view', this.selected.viewId);
    },

    availableGroups() {
      return this.groups.reduce((acc, group) => {
        const views = group.views.filter(view => this.checkReadAccess(view._id));

        if (views.length) {
          acc.push({ ...group, views });
        }

        return acc;
      }, []);
    },

    getViewLink() {
      return (view = {}) => {
        const link = {
          name: 'view',
          params: { id: view._id },
        };

        if (view.tabs && view.tabs.length) {
          link.query = { tabId: view.tabs[0]._id };
        }

        return link;
      };
    },
  },
  methods: {
    selectTab(tabId) {
      this.tabId = tabId;
    },

    submit() {
      /* this.$router.push({
        name: 'view',
        params: {
          id: this.selected.viewId,
        },
      }); */

      this.hideModal();

      this.showSideBar({
        name: SIDE_BARS_BY_WIDGET_TYPES[this.selectedView.type],
        config: {
          tabId: this.selected.tabId,
          widget: this.config.widget,
        },
      });
    },
  },
};
</script>

<style lang="scss">
  .view-select .v-input__control {
    width: 100%;
  }
</style>


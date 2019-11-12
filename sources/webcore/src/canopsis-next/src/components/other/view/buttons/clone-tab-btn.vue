<template lang="pug">
  v-btn(
    data-test="copyTab",
    small,
    flat,
    icon,
    @click.prevent="showSelectViewModal(tab)"
  )
    v-icon(small) file_copy
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { generateCopyOfViewTab, getViewsTabsWidgetsIdsMappings } from '@/helpers/entities';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

const { mapGetters: viewMapGetters, mapActions: viewMapActions } = createNamespacedHelpers('view');

export default {
  mixins: [
    authMixin,
    modalMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    tab: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...viewMapGetters({
      getViewById: 'getItemById',
    }),
  },
  methods: {
    ...viewMapActions({
      updateView: 'update',
    }),

    showSelectViewModal(tab) {
      this.showModal({
        name: MODALS.selectView,
        config: {
          action: viewId => this.showCloneTabModalWithPromise(tab, viewId),
        },
      });
    },

    showCloneTabModalWithPromise(tab, viewId) {
      return new Promise(resolve => this.showModal({
        name: MODALS.textFieldEditor,
        config: {
          title: this.$t('modals.viewTab.duplicate.title'),
          field: {
            name: 'text',
            label: this.$t('modals.viewTab.fields.title'),
            validationRules: 'required',
          },
          action: async (title) => {
            await this.cloneTabAction(tab, viewId, title);

            resolve();
          },
        },
      }));
    },

    async cloneTabAction(tab, viewId, title) {
      const newTab = {
        ...generateCopyOfViewTab(tab),

        title,
      };

      const widgetsIdsMappings = getViewsTabsWidgetsIdsMappings(tab, newTab);

      await this.copyUserPreferencesByWidgetsIdsMappings(widgetsIdsMappings);
      await this.addTabIntoViewById(newTab, viewId);

      this.$router.push({
        name: 'view',
        params: {
          id: viewId,
        },
        query: {
          tabId: newTab._id,
        },
      });
    },

    addTabIntoViewById(tab, viewId) {
      const view = this.getViewById(viewId);

      if (!view) {
        throw new Error('View was not found');
      }

      const data = {
        ...view,

        tabs: [...view.tabs, tab],
      };

      return this.updateView({
        data,

        id: viewId,
      });
    },
  },
};
</script>

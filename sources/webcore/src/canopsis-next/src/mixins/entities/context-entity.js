import { createNamespacedHelpers } from 'vuex';

import popupMixin from '@/mixins/popup';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

/**
 * @mixin Helpers' for context store
 */
export default {
  mixins: [popupMixin],
  computed: {
    ...mapGetters({
      getContextEntitiesListByWidgetId: 'getListByWidgetId',
      getContextEntitiesMetaByWidgetId: 'getMetaByWidgetId',
      getContextEntitiesPendingByWidgetId: 'getPendingByWidgetId',
      getContextEntitiesFetchingParamsByWidgetId: 'getFetchingParamsByWidgetId',
    }),

    contextEntities() {
      return this.getContextEntitiesListByWidgetId(this.widget._id);
    },
    contextEntitiesMeta() {
      return this.getContextEntitiesMetaByWidgetId(this.widget._id);
    },
    contextEntitiesPending() {
      return this.getContextEntitiesPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchContextEntitiesList: 'fetchList',
      removeContextEntity: 'remove',
      updateContextEntity: 'update',
      createContextEntity: 'create',
      refreshContextEntitiesLists: 'refreshLists',
    }),

    async updateContextEntityWithPopup(data) {
      await this.updateContextEntity({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.edit') });
    },

    async createContextEntityWithPopup(data) {
      await this.createContextEntity({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.create') });
    },

    async duplicateContextEntityWithPopup(data) {
      await this.createContextEntity({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.duplicate') });
    },
  },
};

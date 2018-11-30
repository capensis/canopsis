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
      update: 'update',
      create: 'create',
      refreshContextEntitiesLists: 'refreshLists',
    }),

    async updateContextEntity(data) {
      await this.update({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.edit') });
    },

    async createContextEntity(data) {
      await this.create({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.create') });
    },

    async duplicateContextEntity(data) {
      await this.create({ data });
      this.addSuccessPopup({ text: this.$t('modals.createEntity.success.duplicate') });
    },
  },
};

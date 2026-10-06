import { mapActions, mapStores } from 'pinia';

import { useEntityStore } from '@/stores/modules/entity';

import { entityPauseMixin } from '@/mixins/entities/entity-pause';

export const entitiesContextEntityMixin = {
  mixins: [entityPauseMixin],

  computed: {
    ...mapStores(useEntityStore),

    getContextEntitiesListByWidgetId() {
      return this.entityStore.getListByWidgetId;
    },
    getContextEntitiesMetaByWidgetId() {
      return this.entityStore.getMetaByWidgetId;
    },
    getContextEntitiesPendingByWidgetId() {
      return this.entityStore.getPendingByWidgetId;
    },
    getContextEntitiesFetchingParamsByWidgetId() {
      return this.entityStore.getFetchingParamsByWidgetId;
    },

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
    ...mapActions(useEntityStore, {
      fetchContextEntitiesList: 'fetchList',
      removeContextEntity: 'remove',
      updateContextEntity: 'update',
      createContextEntity: 'create',
      fetchContextEntitiesListWithPreviousParams: 'fetchListWithPreviousParams',
      fetchContextEntitiesWithoutStore: 'fetchListWithoutStore',
      fetchBasicContextEntityWithoutStore: 'fetchBasicEntityWithoutStore',
      fetchContextEntityContextGraphWithoutStore: 'fetchContextGraphWithoutStore',
      createContextExport: 'createContextExport',
      fetchContextExport: 'fetchContextExport',
      archiveDisabledEntitiesData: 'archiveDisabledEntitiesData',
      archiveUnlinkedEntitiesData: 'archiveUnlinkedEntitiesData',
      cleanArchivedEntitiesData: 'cleanArchivedEntitiesData',
      bulkEnableEntities: 'bulkEnable',
      bulkDisableEntities: 'bulkDisable',
    }),

    pauseContextEntity(...args) {
      return this.pauseEntity(...args);
    },

    resumeContextEntity(...args) {
      return this.resumeEntity(...args);
    },

    bulkPauseContextEntities(...args) {
      return this.bulkPauseEntities(...args);
    },

    bulkResumeContextEntities(...args) {
      return this.bulkResumeEntities(...args);
    },

    async updateContextEntityWithPopup({ data, id }) {
      await this.updateContextEntity({ data, id });
      this.$popups.success({ text: this.$t('modals.createEntity.success.edit') });
    },
  },
};

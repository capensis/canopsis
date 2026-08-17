<template>
  <div class="entities-list-widget">
    <entities-list-table-with-pagination
      :widget="widget"
      :entities="contextEntities"
      :pending="contextEntitiesPending"
      :meta="contextEntitiesMeta"
      :query="query"
      :columns="widget.parameters.widgetColumns"
      :selectable="!hideMassSelection"
      :hide-actions="hideActions"
      @update:query="query = $event"
    >
      <template #toolbar>
        <c-layout class="entities-list-widget__toolbar gap-4 py-4 pt-0" align="end" wrap>
          <div class="entities-list-widget__toolbar-cell entities-list-widget__toolbar-cell--advanced-search">
            <c-entity-advanced-search
              :searches="searches"
              @submit="updateSearch"
              @reset="resetSearch"
              @toggle-pin:search="togglePinSearch"
              @remove:search="removeSearch"
            />
          </div>
          <div v-if="hasAccessToCategory" class="entities-list-widget__toolbar-cell">
            <c-entity-category-field
              :category="query.category"
              hide-details
              @update:category="updateCategory"
            />
          </div>
          <c-layout
            v-if="hasAccessToUserFilter"
            class="entities-list-widget__toolbar-cell entities-list-widget__toolbar-cell--filter"
            align="end"
          >
            <filter-selector
              :model-value="query.filter"
              :locked-value="query.lockedFilter"
              :label="$t('settings.selectAFilter')"
              :filters="userPreferenceFilters"
              :locked-filters="widget.filters"
              hide-details
              @update:model-value="updateSelectedFilter"
            />
            <filters-list-btn
              :widget-id="widget._id"
              addable
              editable
              private
              with-alarm
              with-entity
              with-pbehavior
              entity-counters-type
            />
          </c-layout>
          <div class="entities-list-widget__toolbar-cell">
            <c-enabled-field
              :model-value="query.no_events"
              :label="$t('context.noEventsFilter')"
              class="mt-0"
              color="primary"
              hide-details
              inline
              no-margin
              @update:model-value="updateNoEvents"
            />
          </div>
          <template v-if="!hideActions">
            <div v-if="hasAccessToCreateEntity">
              <context-fab />
            </div>
            <div v-if="hasAccessToExportAsCsv">
              <c-action-btn
                :loading="downloading"
                :tooltip="$t('settings.exportAsCsv')"
                icon="cloud_download"
                @click="exportContextList"
              />
            </div>
          </template>
        </c-layout>
      </template>
    </entities-list-table-with-pagination>
  </div>
</template>

<script>
import { isObject } from 'lodash';

import { USER_PERMISSIONS } from '@/constants';

import { getContextExportDownloadFileUrl } from '@/helpers/entities/entity/url';

import { authMixin } from '@/mixins/auth';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';
import { permissionsWidgetsContextCategory } from '@/mixins/permissions/widgets/context/category';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';
import { widgetAdvancedSearchSavedItemsMixin } from '@/mixins/widget/advanced-search/saved-items';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

import ContextFab from './actions/context-fab.vue';
import EntitiesListTableWithPagination from './partials/entities-list-table-with-pagination.vue';

/**
 * Context explorer widget runtime foundation.
 */
export default {
  components: {
    ContextFab,
    EntitiesListTableWithPagination,
    FilterSelector,
    FiltersListBtn,
  },

  mixins: [
    authMixin,
    widgetAdvancedSearchSavedItemsMixin,
    widgetFetchQueryMixin,
    entitiesContextEntityMixin,
    permissionsWidgetsContextFilters,
    permissionsWidgetsContextCategory,
    exportMixinCreator({
      createExport: 'createContextExport',
      fetchExport: 'fetchContextExport',
    }),
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hideMassSelection: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      downloading: false,
    };
  },

  computed: {
    hasAccessToCreateEntity() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.createEntity);
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USER_PERMISSIONS.business.context.actions.exportAsCsv);
    },

    userPreferenceFilters() {
      return this.userPreference.filters ?? [];
    },
  },

  methods: {
    updateNoEvents(noEvents) {
      this.updateContentInUserPreference({
        noEvents,
      });

      this.query = {
        ...this.query,

        page: 1,
        no_events: noEvents,
      };
    },

    updateCategory(category) {
      const categoryId = category?._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,

        page: 1,
        category: categoryId,
      };
    },

    updateSelectedFilter(filter) {
      this.updateContentInUserPreference({ mainFilter: filter });

      this.query = {
        ...this.query,
        page: 1,
        filter,
      };
    },

    fetchList() {
      if (this.widget.parameters.widgetColumns.length) {
        const params = this.getQuery();

        params.with_flags = true;

        this.fetchContextEntitiesList({
          widgetId: this.widget._id,
          params,
        });
      }
    },

    getExportQuery() {
      const query = this.getQuery();
      const {
        widgetExportColumns,
        widgetColumns,
        exportCsvSeparator,
        exportCsvDatetimeFormat,
      } = this.widget.parameters;
      const columns = widgetExportColumns?.length ? widgetExportColumns : widgetColumns;

      return {
        fields: columns.map(({ value, text }) => ({ name: value, label: text })),
        search: query.search,
        category: query.category,
        filters: query.filters,
        separator: exportCsvSeparator,
        /**
         * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
         */
        time_format: isObject(exportCsvDatetimeFormat) ? exportCsvDatetimeFormat.value : exportCsvDatetimeFormat,
      };
    },

    async exportContextList() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          widgetId: this.widget._id,
          data: this.getExportQuery(),
        });

        this.downloadFile(getContextExportDownloadFileUrl(fileData._id));
      } catch {
        this.$popups.error({ text: this.$t('context.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.entities-list-widget__toolbar-cell {
  min-width: 180px;
}

.entities-list-widget__toolbar-cell--advanced-search {
  flex: 1 1 360px;
  min-width: 320px;
}

.entities-list-widget__toolbar-cell--filter {
  min-width: 260px;
}
</style>

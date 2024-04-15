<template>
  <entities-list-table-with-pagination
    :widget="widget"
    :entities="contextEntities"
    :pending="contextEntitiesPending"
    :meta="contextEntitiesMeta"
    :query.sync="query"
    :columns="widget.parameters.widgetColumns"
    selectable
  >
    <template #toolbar="">
      <v-layout class="gap-4 pa-4 pt-0" wrap align-end>
        <v-flex>
          <c-advanced-search
            :columns="widget.parameters.widgetColumns"
            :saved-items="searches"
            combobox
            @submit="updateSearchInQuery"
            @add:item="addSearchIntoUserPreferences"
            @toggle-pin:item="togglePinSearchInUserPreferences"
            @remove:item="removeSearchFromUserPreferences"
          />
        </v-flex>
        <v-flex v-if="hasAccessToCategory">
          <c-entity-category-field
            :category="query.category"
            hide-details
            @input="updateCategory"
          />
        </v-flex>
        <v-flex>
          <v-layout
            v-if="hasAccessToUserFilter"
            align-end
          >
            <filter-selector
              :value="query.filter"
              :locked-value="query.lockedFilter"
              :label="$t('settings.selectAFilter')"
              :filters="userPreference.filters"
              :locked-filters="widget.filters"
              :disabled="!hasAccessToListFilters"
              hide-details
              @input="updateSelectedFilter"
            />
            <filters-list-btn
              v-if="hasAccessToAddFilter || hasAccessToEditFilter"
              :widget-id="widget._id"
              :addable="hasAccessToAddFilter"
              :editable="hasAccessToEditFilter"
              private
              with-alarm
              with-entity
              with-pbehavior
              entity-counters-type
            />
          </v-layout>
        </v-flex>
        <v-flex>
          <v-checkbox
            :input-value="query.no_events"
            :label="$t('context.noEventsFilter')"
            color="primary"
            hide-details
            @change="updateNoEvents"
          />
        </v-flex>
        <v-flex v-if="hasAccessToCreateEntity">
          <context-fab />
        </v-flex>
        <v-flex v-if="hasAccessToExportAsCsv">
          <c-action-btn
            :loading="downloading"
            :tooltip="$t('settings.exportAsCsv')"
            icon="cloud_download"
            @click="exportContextList"
          />
        </v-flex>
      </v-layout>
    </template>
  </entities-list-table-with-pagination>
</template>

<script>
import { isObject } from 'lodash';

import { USERS_PERMISSIONS } from '@/constants';

import { getContextExportDownloadFileUrl } from '@/helpers/entities/entity/url';

import { authMixin } from '@/mixins/auth';
import { widgetSearchMixin } from '@/mixins/widget/search';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';
import { permissionsWidgetsContextCategory } from '@/mixins/permissions/widgets/context/category';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

import ContextFab from './actions/context-fab.vue';
import EntitiesListTableWithPagination from './partials/entities-list-table-with-pagination.vue';

export default {
  components: {
    FilterSelector,
    FiltersListBtn,
    ContextFab,
    EntitiesListTableWithPagination,
  },
  mixins: [
    authMixin,
    widgetSearchMixin,
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
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
  },
  data() {
    return {
      downloading: false,
    };
  },
  computed: {
    hasAccessToCreateEntity() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.createEntity);
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.exportAsCsv);
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
      const categoryId = category && category._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,

        page: 1,
        category: categoryId,
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
      } catch (err) {
        this.$popups.error({ text: this.$t('context.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>

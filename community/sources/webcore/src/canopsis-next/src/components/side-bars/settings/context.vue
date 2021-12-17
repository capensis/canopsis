<template lang="pug">
  div
    v-list.pt-0(expand)
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="settings.widget.parameters.widgetColumns",
            :columns-label="$t('settings.columnName')"
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.widgetColumns",
            :label="$t('settings.columnNames')"
          )
          v-divider
          field-columns(
            v-model="settings.widget.parameters.serviceDependenciesColumns",
            :label="$t('settings.treeOfDependenciesColumnNames')",
            with-color-indicator
          )
          v-divider
          template(v-if="hasAccessToListFilters")
            field-filters(
              v-model="settings.widget.parameters.mainFilter",
              :entities-type="$constants.ENTITIES_TYPES.entity",
              :filters.sync="settings.widget.parameters.viewFilters",
              :condition.sync="settings.widget.parameters.mainFilterCondition",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
          field-context-entities-types-filter(v-model="settings.widget.parameters.selectedTypes")
          v-divider
          export-csv-form(v-model="settings.widget.parameters")
      v-divider
    v-btn.primary(data-test="submitContext", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { SIDE_BARS, USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { widgetSettingsMixin } from '@/mixins/widget/settings';

import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldContextEntitiesTypesFilter from './fields/context/context-entities-types-filter.vue';
import ExportCsvForm from './forms/export-csv.vue';

export default {
  name: SIDE_BARS.contextSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldFilters,
    FieldContextEntitiesTypesFilter,
    ExportCsvForm,
  },
  mixins: [
    authMixin,
    widgetSettingsMixin,
  ],
  data() {
    const { widget } = this.config;

    return {
      settings: {
        widget: cloneDeep(widget),
      },
    };
  },
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.addFilter);
    },
  },
  methods: {
    prepareWidgetQuery(newQuery, oldQuery) {
      return {
        searchFilter: oldQuery.searchFilter,

        ...newQuery,
      };
    },

    updateMainFilterUpdatedAt() {
      this.settings.widget.parameters.mainFilterUpdatedAt = Date.now();
    },
  },
};
</script>

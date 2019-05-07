<template lang="pug">
  div
    v-list.pt-0(expand)
      field-row-grid-size(
      :rowId.sync="settings.rowId",
      :size.sync="settings.widget.size",
      :availableRows="availableRows",
      @createRow="createRow"
      )
      v-divider
      field-title(v-model="settings.widget.title", :title="$t('common.title')")
      v-divider
      field-periodic-refresh(v-model="settings.widget.parameters.periodicRefresh")
      v-divider
      v-list-group
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
          v-model="settings.widget.parameters.sort",
          :columns="settings.widget.parameters.widgetColumns"
          )
          v-divider
          field-columns(v-model="settings.widget.parameters.widgetColumns", withHtml)
          v-divider
          field-default-elements-per-page(v-model="settings.widget_preferences.itemsPerPage")
          v-divider
          field-opened-resolved-filter(v-model="settings.widget.parameters.alarmsStateFilter")
          v-divider
          template(v-if="hasAccessToListFilters")
            field-filters(
            v-model="settings.widget.parameters.mainFilter",
            :filters.sync="settings.widget.parameters.viewFilters",
            :condition.sync="settings.widget.parameters.mainFilterCondition",
            :hasAccessToAddFilter="hasAccessToAddFilter",
            :hasAccessToEditFilter="hasAccessToEditFilter"
            )
            v-divider
          field-info-popup(
          v-model="settings.widget.parameters.infoPopups",
          :columns="settings.widget.parameters.widgetColumns",
          )
          v-divider
          field-text-editor(
          v-model="settings.widget.parameters.moreInfoTemplate",
          :title="$t('settings.moreInfosModal')"
          )
          v-divider
          field-switcher(
          v-model="settings.widget.parameters.isAckNoteRequired",
          :title="$t('settings.isAckNoteRequired')",
          )
          v-divider
          field-switcher(
          v-model="settings.widget.parameters.isMultiAckEnabled",
          :title="$t('settings.isMultiAckEnabled')",
          )
          v-divider
          field-switcher(
          v-model="settings.widget.parameters.isHtmlEnabledOnTimeLine",
          :title="$t('settings.isHtmlEnabledOnTimeLine')",
          )
      v-divider
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import { get, cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldSwitcher from './fields/common/switcher.vue';

/**
 * Component to regroup the alarms list settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FieldRowGridSize,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldInfoPopup,
    FieldTextEditor,
    FieldSwitcher,
  },
  mixins: [authMixin, widgetSettingsMixin, sideBarSettingsWidgetAlarmMixin],
  data() {
    const { widget, rowId } = this.config;

    return {
      settings: {
        rowId,
        widget: this.prepareAlarmWidgetSettings(cloneDeep(widget), true),
        widget_preferences: {
          itemsPerPage: PAGINATION_LIMIT,
        },
      },
    };
  },
  computed: {
    hasAccessToListFilters() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.editFilter);
    },

    hasAccessToAddFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.addFilter);
    },
  },
  mounted() {
    const { widget_preferences: widgetPreference } = this.userPreference;

    this.settings.widget_preferences = {
      itemsPerPage: get(widgetPreference, 'itemsPerPage', PAGINATION_LIMIT),
    };
  },
  methods: {
    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareAlarmWidgetSettings(widget);
    },
  },
};
</script>

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
      v-list-group(data-test="advancedSettings")
        v-list-tile(slot="activator") {{ $t('settings.advancedSettings') }}
        v-list.grey.lighten-4.px-2.py-0(expand)
          field-default-sort-column(
            v-model="settings.widget.parameters.sort",
            :columns="settings.widget.parameters.widgetColumns",
            :columnsLabel="$t('settings.columnName')"
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
              :entitiesType="$constants.ENTITIES_TYPES.alarm",
              :filters.sync="settings.widget.parameters.viewFilters",
              :condition.sync="settings.widget.parameters.mainFilterCondition",
              :hasAccessToAddFilter="hasAccessToAddFilter",
              :hasAccessToEditFilter="hasAccessToEditFilter",
              @input="updateMainFilterUpdatedAt"
            )
            v-divider
          field-live-reporting(v-model="settings.widget.parameters.liveReporting")
          v-divider
          field-info-popup(
            v-model="settings.widget.parameters.infoPopups",
            :columns="settings.widget.parameters.widgetColumns"
          )
          v-divider
          field-text-editor(
            data-test="widgetMoreInfoTemplate",
            v-model="settings.widget.parameters.moreInfoTemplate",
            :title="$t('settings.moreInfosModal')"
          )
          v-divider
          field-grid-range-size(
            v-model="settings.widget.parameters.expandGridRangeSize",
            :title="$t('settings.expandGridRangeSize')"
          )
          v-divider
          field-switcher(
            data-test="isHtmlEnabledOnTimeLine",
            v-model="settings.widget.parameters.isHtmlEnabledOnTimeLine",
            :title="$t('settings.isHtmlEnabledOnTimeLine')"
          )
          v-divider
          v-list-group(data-test="ackGroup")
            v-list-tile(slot="activator") Ack
            v-list.grey.lighten-4.px-2.py-0(expand)
            field-switcher(
              data-test="isAckNoteRequired",
              v-model="settings.widget.parameters.isAckNoteRequired",
              :title="$t('settings.isAckNoteRequired')"
            )
            v-divider
            field-switcher(
              data-test="isMultiAckEnabled",
              v-model="settings.widget.parameters.isMultiAckEnabled",
              :title="$t('settings.isMultiAckEnabled')"
            )
            v-divider
            field-fast-ack-output(v-model="settings.widget.parameters.fastAckOutput")
      v-divider
    copy-widget-id(:widgetId="settings.widget._id")
    v-btn.primary(data-test="submitAlarms", @click="submit") {{ $t('common.save') }}
</template>

<script>
import { get, cloneDeep } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SIDE_BARS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import widgetSettingsMixin from '@/mixins/widget/settings';
import sideBarSettingsWidgetAlarmMixin from '@/mixins/side-bar/settings/widgets/alarm';
import CopyWidgetId from '@/components/side-bars/settings/widgets/fields/common/copy-widget-id.vue';

import FieldRowGridSize from './fields/common/row-grid-size.vue';
import FieldTitle from './fields/common/title.vue';
import FieldDefaultSortColumn from './fields/common/default-sort-column.vue';
import FieldColumns from './fields/common/columns.vue';
import FieldLiveReporting from './fields/common/live-reporting.vue';
import FieldPeriodicRefresh from './fields/common/periodic-refresh.vue';
import FieldDefaultElementsPerPage from './fields/common/default-elements-per-page.vue';
import FieldOpenedResolvedFilter from './fields/alarm/opened-resolved-filter.vue';
import FieldFilters from './fields/common/filters.vue';
import FieldInfoPopup from './fields/alarm/info-popup.vue';
import FieldTextEditor from './fields/common/text-editor.vue';
import FieldSwitcher from './fields/common/switcher.vue';
import FieldFastAckOutput from './fields/alarm/fast-ack-output.vue';
import FieldGridRangeSize from './fields/common/grid-range-size.vue';

/**
 * Component to regroup the alarms list settings fields
 */
export default {
  name: SIDE_BARS.alarmSettings,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CopyWidgetId,
    FieldRowGridSize,
    FieldTitle,
    FieldDefaultSortColumn,
    FieldColumns,
    FieldLiveReporting,
    FieldPeriodicRefresh,
    FieldDefaultElementsPerPage,
    FieldOpenedResolvedFilter,
    FieldFilters,
    FieldInfoPopup,
    FieldTextEditor,
    FieldSwitcher,
    FieldFastAckOutput,
    FieldGridRangeSize,
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
    updateMainFilterUpdatedAt() {
      this.settings.widget.parameters.mainFilterUpdatedAt = Date.now();
    },

    prepareWidgetSettings() {
      const { widget } = this.settings;

      return this.prepareAlarmWidgetSettings(widget);
    },
  },
};
</script>

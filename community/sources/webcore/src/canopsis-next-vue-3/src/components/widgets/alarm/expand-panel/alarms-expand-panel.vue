<template>
  <div class="alarms-expand-panel bg-secondary">
    <v-tabs
      v-model="activeTab"
      bg-color="secondary"
      color="primary"
      align-tabs="center"
    >
      <v-tab
        v-if="hasMoreInfos"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos"
      >
        {{ $t('alarm.tabs.moreInfos') }}
      </v-tab>
      <v-tab :value="$constants.ALARMS_EXPAND_PANEL_TABS.timeLine">
        {{ $t('alarm.tabs.timeLine') }}
      </v-tab>
      <v-tab
        v-if="hasTickets"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared"
      >
        {{ $t('alarm.tabs.ticketsDeclared') }}
      </v-tab>
      <v-tab
        v-if="hasWidgetCharts"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.charts"
      >
        {{ $t('alarm.tabs.charts') }}
      </v-tab>
      <v-tab
        v-if="hasAlarmsListPbehaviorAccess"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior"
      >
        {{ $t('common.pbehavior', 2) }}
      </v-tab>
      <v-tab
        v-if="hasChildren"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren"
      >
        {{ $t('alarm.tabs.alarmsChildren') }}
      </v-tab>
      <v-tab
        v-if="hasServiceDependencies"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.trackSource"
      >
        {{ $t('alarm.tabs.trackSource') }}
      </v-tab>
      <v-tab
        v-if="hasImpactsDependencies"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.impactChain"
      >
        {{ $t('alarm.tabs.impactChain') }}
      </v-tab>
      <v-tab
        v-if="hasEntityGantt"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt"
      >
        {{ $t('alarm.tabs.entityGantt') }}
      </v-tab>
      <v-tab :value="$constants.ALARMS_EXPAND_PANEL_TABS.remediation">
        {{ $t('alarm.tabs.remediation') }}
      </v-tab>
      <v-tab
        v-if="isAvailabilityEnabled"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.availability"
      >
        {{ $t('common.availability') }}
      </v-tab>
      <v-tab
        v-if="hasEntityEnrichments"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityEnrichments"
      >
        {{ $t('alarm.tabs.entityEnrichments') }}
      </v-tab>
    </v-tabs>

    <v-window
      v-model="activeTab"
      :transition="false"
      :reverse-transition="false"
    >
      <v-window-item
        v-if="hasMoreInfos"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <alarms-expand-panel-more-infos
            :alarm="alarm"
            :template="widget.parameters.moreInfoTemplate"
            :template-id="moreInfosTemplateId"
            :selected-tags="selectedTags"
            @select:tag="$emit('select:tag', $event)"
            @remove:tag="$emit('remove:tag', $event)"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item :value="$constants.ALARMS_EXPAND_PANEL_TABS.timeLine">
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <template #prepend>
            <v-progress-linear
              :active="pending"
              :height="3"
              indeterminate
            />
          </template>
          <alarm-timeline
            :steps="timelineSteps"
            :meta="stepsMeta"
            :query="stepsQuery"
            :is-html-enabled="isHtmlEnabled"
            @update:query="updateStepsQuery"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasTickets"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <declared-tickets-list
            :tickets="alarm.v.tickets"
            :parent-alarm-id="parentAlarmId"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasWidgetCharts"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.charts"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <entity-charts
            :charts="widget.parameters.charts"
            :entity="alarm.entity"
            :available-metrics="filteredPerfData"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasAlarmsListPbehaviorAccess"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <pbehaviors-simple-list
            :entity="alarm.entity"
            :removable="hasDeleteAnyPbehaviorAccess"
            :updatable="hasUpdateAnyPbehaviorAccess"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasChildren"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <template #prepend>
            <v-progress-linear
              :active="pending"
              :height="3"
              indeterminate
            />
          </template>
          <alarms-expand-panel-children
            :children="children"
            :alarm="alarm"
            :widget="widget"
            :pending="pending"
            :query="childrenQuery"
            :refresh-alarms-list="fetchList"
            @update:query="updateChildrenQuery"
            @select:tag="$emit('select:tag', $event)"
            @remove:tag="$emit('remove:tag', $event)"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasServiceDependencies"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.trackSource"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <service-dependencies
            :root="dependency"
            :columns="widget.parameters.serviceDependenciesColumns"
            include-root
            show-state-setting
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasImpactsDependencies"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.impactChain"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <service-dependencies
            :root="dependency"
            :columns="widget.parameters.serviceDependenciesColumns"
            include-root
            impact
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasEntityGantt"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <entity-gantt :alarm="alarm" />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item :value="$constants.ALARMS_EXPAND_PANEL_TABS.remediation">
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <alarms-expand-panel-remediation :alarm="alarm" />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="isAvailabilityEnabled"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.availability"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <entity-availability
            :entity="alarm.entity"
            :default-time-range="widget.parameters.availability?.default_time_range"
            :default-show-type="widget.parameters.availability?.default_show_type"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>

      <v-window-item
        v-if="hasEntityEnrichments"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityEnrichments"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <alarms-expand-panel-entity-enrichments :alarm="alarm" />
        </alarms-expand-panel-tab-item-wrapper>
      </v-window-item>
    </v-window>
  </div>
</template>

<script>
import { isEqual, map } from 'lodash';
import { mapActions, mapState } from 'pinia';

import {
  ALARMS_EXPAND_PANEL_TABS,
  ENTITY_TYPES,
  JUNIT_ALARM_CONNECTOR,
  USER_PERMISSIONS,
} from '@/constants';

import { useAlarmStore } from '@/stores/modules/alarm';

import { getAlarmWidgetMoreInfoTemplateId } from '@/helpers/entities/alarm/list';
import {
  convertAlarmDetailsQueryToRequest,
  prepareAlarmDetailsQuery,
} from '@/helpers/entities/alarm/query';
import { convertWidgetChartsToPerfDataQuery } from '@/helpers/entities/metric/query';
import { alarmToServiceDependency } from '@/helpers/entities/service-dependencies/list';
import { getFlexClassesForGridRangeSize } from '@/helpers/entities/shared/grid';
import { setField } from '@/helpers/immutable';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import DeclaredTicketsList from '@/components/other/declare-ticket/declared-tickets-list.vue';
import EntityAvailability from '@/components/other/entity/entity-availability.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';
import EntityCharts from '@/components/widgets/chart/entity-charts.vue';

import EntityGantt from '../entity-gantt/entity-gantt.vue';
import AlarmTimeline from '../timeline/alarm-timeline.vue';

import AlarmsExpandPanelChildren from './alarms-expand-panel-children.vue';
import AlarmsExpandPanelEntityEnrichments from './alarms-expand-panel-entity-enrichments.vue';
import AlarmsExpandPanelMoreInfos from './alarms-expand-panel-more-infos.vue';
import AlarmsExpandPanelRemediation from './alarms-expand-panel-remediation.vue';
import AlarmsExpandPanelTabItemWrapper from './alarms-expand-panel-tab-item-wrapper.vue';

export default {
  components: {
    AlarmTimeline,
    AlarmsExpandPanelChildren,
    AlarmsExpandPanelEntityEnrichments,
    AlarmsExpandPanelMoreInfos,
    AlarmsExpandPanelRemediation,
    AlarmsExpandPanelTabItemWrapper,
    DeclaredTicketsList,
    EntityAvailability,
    EntityCharts,
    EntityGantt,
    PbehaviorsSimpleList,
    ServiceDependencies,
  },

  mixins: [authMixin, entitiesInfoMixin],

  props: {
    alarm: {
      type: Object,
      required: true,
    },
    parentAlarmId: {
      type: String,
      default: '',
    },
    widget: {
      type: Object,
      required: true,
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    search: {
      type: String,
      default: '',
    },
    selectedTags: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['select:tag', 'remove:tag'],

  data() {
    return {
      activeTab: undefined,
      pending: false,
      query: prepareAlarmDetailsQuery(this.alarm, this.widget, this.search),
    };
  },

  computed: {
    ...mapState(useAlarmStore, {
      alarmDetailsById(store) {
        return store.getDetailsByWidgetId(this.widget._id);
      },
    }),

    alarmDetails() {
      return {
        ...(this.alarm.data ?? this.alarm),
        ...(this.alarmDetailsById[this.alarm._id]?.data ?? this.alarmDetailsById[this.alarm._id]),
      };
    },

    children() {
      return this.alarmDetails.children ?? {};
    },

    childrenItems() {
      return Array.isArray(this.children) ? this.children : (this.children.data ?? []);
    },

    childrenCount() {
      if (this.childrenItems.length) {
        return this.childrenItems.length;
      }

      const { opened } = this.widget.parameters;

      if (opened === true) {
        return this.alarm.opened_children ?? this.alarm.children ?? 0;
      }

      if (opened === false) {
        return this.alarm.closed_children ?? this.alarm.children ?? 0;
      }

      return this.alarm.children ?? 0;
    },

    steps() {
      return this.alarmDetails.steps ?? {};
    },

    timelineSteps() {
      return Array.isArray(this.steps) ? this.steps : (this.steps.data ?? []);
    },

    stepsMeta() {
      return Array.isArray(this.steps) ? {} : (this.steps.meta ?? {});
    },

    filteredPerfData() {
      return this.alarmDetails.filtered_perf_data ?? [];
    },

    entity() {
      return {
        ...this.alarm.entity,
        ...this.alarmDetails?.entity,
      };
    },

    dependency() {
      const alarmWithDependenciesCounts = setField(this.alarm, 'entity', entity => ({
        ...entity,
        ...this.entity,
      }));

      return alarmToServiceDependency(alarmWithDependenciesCounts);
    },

    childrenQuery() {
      return this.query.children;
    },

    stepsQuery() {
      return this.query.steps;
    },

    availableTabs() {
      return [
        this.hasMoreInfos && ALARMS_EXPAND_PANEL_TABS.moreInfos,
        ALARMS_EXPAND_PANEL_TABS.timeLine,
        this.hasTickets && ALARMS_EXPAND_PANEL_TABS.ticketsDeclared,
        this.hasWidgetCharts && ALARMS_EXPAND_PANEL_TABS.charts,
        this.hasAlarmsListPbehaviorAccess && ALARMS_EXPAND_PANEL_TABS.pbehavior,
        this.hasChildren && ALARMS_EXPAND_PANEL_TABS.alarmsChildren,
        this.hasServiceDependencies && ALARMS_EXPAND_PANEL_TABS.trackSource,
        this.hasImpactsDependencies && ALARMS_EXPAND_PANEL_TABS.impactChain,
        this.hasEntityGantt && ALARMS_EXPAND_PANEL_TABS.entityGantt,
        ALARMS_EXPAND_PANEL_TABS.remediation,
        this.isAvailabilityEnabled && ALARMS_EXPAND_PANEL_TABS.availability,
        this.hasEntityEnrichments && ALARMS_EXPAND_PANEL_TABS.entityEnrichments,
      ].filter(Boolean);
    },

    cardFlexClass() {
      return getFlexClassesForGridRangeSize(this.widget.parameters.expandGridRangeSize);
    },

    hasMoreInfos() {
      return !!this.widget.parameters.moreInfoTemplate;
    },

    isHtmlEnabled() {
      return !!this.widget.parameters.isHtmlEnabledOnTimeLine;
    },

    hasChildren() {
      return this.childrenCount > 0 && !this.hideChildren;
    },

    hasTickets() {
      return !!this.alarm.v?.tickets?.length;
    },

    hasServiceDependencies() {
      return this.entity.type === ENTITY_TYPES.service;
    },

    hasImpactsDependencies() {
      return (this.entity.impacts_count ?? 0) > 0;
    },

    hasEntityGantt() {
      return (
        this.isProVersion &&
        this.alarm.v?.connector === JUNIT_ALARM_CONNECTOR &&
        [ENTITY_TYPES.component, ENTITY_TYPES.resource].includes(this.entity.type)
      );
    },

    hasWidgetCharts() {
      return !!this.widget.parameters.charts?.length && !!this.filteredPerfData.length;
    },

    hasEntityEnrichments() {
      return this.isProVersion;
    },

    isAvailabilityEnabled() {
      return !!this.widget.parameters.availability?.enabled;
    },

    hasAlarmsListPbehaviorAccess() {
      return (
        this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.pbehaviorAdd) ||
        this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.fastPbehaviorAdd)
      );
    },

    hasUpdateAnyPbehaviorAccess() {
      return this.checkUpdateAccess(USER_PERMISSIONS.technical.exploitation.pbehavior);
    },

    hasDeleteAnyPbehaviorAccess() {
      return this.checkDeleteAccess(USER_PERMISSIONS.technical.exploitation.pbehavior);
    },

    moreInfosTemplateId() {
      return getAlarmWidgetMoreInfoTemplateId(this.widget._id);
    },
  },

  watch: {
    availableTabs: {
      immediate: true,
      handler(availableTabs) {
        if (!availableTabs.includes(this.activeTab)) {
          [this.activeTab] = availableTabs;
        }
      },
    },

    'widget.parameters.opened': {
      handler(opened) {
        this.query = {
          ...this.query,
          opened,
        };
      },
    },

    'widget.parameters.charts': {
      handler(charts) {
        this.query = {
          ...this.query,
          perf_data: convertWidgetChartsToPerfDataQuery(charts),
        };
      },
    },

    'widget.parameters.widgetGroupColumns': {
      handler(columns) {
        this.query = {
          ...this.query,
          search_by: map(columns, 'value'),
        };
      },
    },

    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },

    search: {
      immediate: true,
      handler(search, prevSearch = '') {
        if (search !== prevSearch) {
          this.query = {
            ...this.query,
            search,
          };
        }
      },
    },
  },

  beforeUnmount() {
    this.removeDetailsQuery({ widgetId: this.widget._id, id: this.alarm._id });
  },

  methods: {
    ...mapActions(useAlarmStore, ['fetchAlarmDetails', 'removeDetailsQuery']),

    updateChildrenQuery(children) {
      this.query = {
        ...this.query,
        children,
      };
    },

    updateStepsQuery(steps) {
      this.query = {
        ...this.query,
        steps,
      };
    },

    async fetchList() {
      this.pending = true;

      try {
        await this.fetchAlarmDetails({
          widgetId: this.widget._id,
          id: this.alarm._id,
          query: convertAlarmDetailsQueryToRequest(this.query),
        });
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<template>
  <v-tabs
    v-model="activeTab"
    :key="tabsKey"
    class="expand-panel secondary lighten-2"
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab
      v-if="hasMoreInfos"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos}`"
    >
      {{ $t('alarm.tabs.moreInfos') }}
    </v-tab>
    <v-tab :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.timeLine}`">
      {{ $t('alarm.tabs.timeLine') }}
    </v-tab>
    <v-tab
      v-if="hasWidgetCharts"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.charts}`"
    >
      {{ $t('alarm.tabs.charts') }}
    </v-tab>
    <v-tab
      v-if="hasTickets"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared}`"
    >
      {{ $t('alarm.tabs.ticketsDeclared') }}
    </v-tab>
    <v-tab :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior}`">
      {{ $tc('common.pbehavior', 2) }}
    </v-tab>
    <v-tab
      v-if="hasChildren"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren}`"
    >
      {{ $t('alarm.tabs.alarmsChildren') }}
    </v-tab>
    <v-tab
      v-if="hasServiceDependencies"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.trackSource}`"
    >
      {{ $t('alarm.tabs.trackSource') }}
    </v-tab>
    <v-tab
      v-if="hasImpactsDependencies"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.impactChain}`"
    >
      {{ $t('alarm.tabs.impactChain') }}
    </v-tab>
    <v-tab
      v-if="hasEntityGantt"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt}`"
    >
      {{ $t('alarm.tabs.entityGantt') }}
    </v-tab>
    <v-tab
      v-if="isAvailabilityEnabled"
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.availability}`"
    >
      {{ $t('common.availability') }}
    </v-tab>
    <v-tabs-items
      v-model="activeTab"
      mandatory
    >
      <v-tab-item
        v-if="hasMoreInfos"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <alarms-expand-panel-more-infos
            :alarm="alarm"
            :template="widget.parameters.moreInfoTemplate"
            @select:tag="$emit('select:tag', $event)"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item :value="$constants.ALARMS_EXPAND_PANEL_TABS.timeLine">
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <template #prepend>
            <v-progress-linear
              :active="pending"
              :height="3"
              indeterminate
            />
          </template>
          <alarms-time-line
            :steps="steps"
            :is-html-enabled="isHtmlEnabled"
            @update:page="updateStepsQueryPage"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item
        v-if="hasTickets"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <declared-tickets-list
            :tickets="alarm.v.tickets"
            :parent-alarm-id="parentAlarmId"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item
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
      </v-tab-item>
      <v-tab-item :value="$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior">
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <pbehaviors-simple-list
            :entity="alarm.entity"
            :removable="hasDeleteAnyPbehaviorAccess"
            :updatable="hasUpdateAnyPbehaviorAccess"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item
        v-if="hasChildren"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <alarms-expand-panel-children
            :children="children"
            :alarm="alarm"
            :widget="widget"
            :pending="pending"
            :query.sync="childrenQuery"
            :refresh-alarms-list="fetchList"
          />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item
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
      </v-tab-item>
      <v-tab-item
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
      </v-tab-item>
      <v-tab-item
        v-if="hasEntityGantt"
        :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt"
      >
        <alarms-expand-panel-tab-item-wrapper :card-flex-class="cardFlexClass">
          <entity-gantt :alarm="alarm" />
        </alarms-expand-panel-tab-item-wrapper>
      </v-tab-item>
      <v-tab-item
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
      </v-tab-item>
    </v-tabs-items>
  </v-tabs>
</template>

<script>
import { isEqual, map } from 'lodash';

import { ENTITY_TYPES, JUNIT_ALARM_CONNECTOR } from '@/constants';

import { uid } from '@/helpers/uid';
import { setField } from '@/helpers/immutable';
import { alarmToServiceDependency } from '@/helpers/entities/service-dependencies/list';
import { convertAlarmDetailsQueryToRequest } from '@/helpers/entities/alarm/query';
import { convertWidgetChartsToPerfDataQuery } from '@/helpers/entities/metric/query';
import { getFlexClassesForGridRangeSize } from '@/helpers/entities/shared/grid';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { widgetExpandPanelAlarmDetails } from '@/mixins/widget/expand-panel/alarm/details';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import DeclaredTicketsList from '@/components/other/declare-ticket/declared-tickets-list.vue';
import EntityCharts from '@/components/widgets/chart/entity-charts.vue';
import EntityAvailability from '@/components/other/entity/entity-availability.vue';

import AlarmsTimeLine from '../time-line/alarms-time-line.vue';
import EntityGantt from '../entity-gantt/entity-gantt.vue';

import AlarmsExpandPanelTabItemWrapper from './alarms-expand-panel-tab-item-wrapper.vue';
import AlarmsExpandPanelMoreInfos from './alarms-expand-panel-more-infos.vue';
import AlarmsExpandPanelChildren from './alarms-expand-panel-children.vue';

export default {
  components: {
    EntityAvailability,
    AlarmsExpandPanelTabItemWrapper,
    EntityCharts,
    DeclaredTicketsList,
    PbehaviorsSimpleList,
    ServiceDependencies,
    AlarmsTimeLine,
    EntityGantt,
    AlarmsExpandPanelMoreInfos,
    AlarmsExpandPanelChildren,
  },
  mixins: [
    entitiesInfoMixin,
    widgetExpandPanelAlarmDetails,
    permissionsTechnicalExploitationPbehaviorMixin,
  ],
  props: {
    alarm: {
      type: Object,
      required: true,
    },
    parentAlarmId: {
      type: String,
      required: false,
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
  },
  data() {
    return {
      activeTab: undefined,
      tabsKey: uid(),
    };
  },
  computed: {
    cardFlexClass() {
      return getFlexClassesForGridRangeSize(this.widget.parameters.expandGridRangeSize);
    },

    isHtmlEnabled() {
      return this.widget.parameters.isHtmlEnabledOnTimeLine;
    },

    dependency() {
      const alarmWithDependenciesCounts = setField(this.alarm, 'entity', entity => ({
        ...entity,
        ...this.alarmDetails.entity,
      }));

      return alarmToServiceDependency(alarmWithDependenciesCounts);
    },

    hasMoreInfos() {
      return this.widget.parameters.moreInfoTemplate;
    },

    hasChildren() {
      return this.alarm.children && !this.hideChildren;
    },

    hasTickets() {
      return this.alarm.v.tickets?.length;
    },

    hasServiceDependencies() {
      return this.alarm.entity.type === ENTITY_TYPES.service;
    },

    hasImpactsDependencies() {
      const { impacts_count: impactsCount } = this.alarm.entity;

      return impactsCount > 0;
    },

    hasEntityGantt() {
      /**
       * We have junit feature only on `pro` version of canopsis
       */
      return this.isProVersion
        && this.alarm.v.connector === JUNIT_ALARM_CONNECTOR
        && [ENTITY_TYPES.component, ENTITY_TYPES.resource].includes(this.alarm.entity.type);
    },

    hasWidgetCharts() {
      return this.widget.parameters.charts?.length && this.filteredPerfData.length;
    },

    isAvailabilityEnabled() {
      return this.widget.parameters.availability?.enabled;
    },
  },
  watch: {
    'widget.parameters.moreInfoTemplate': {
      handler() {
        this.refreshTabs();
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
      handler(search) {
        this.query = {
          ...this.query,

          search,
        };
      },
    },
  },
  beforeDestroy() {
    return this.removeAlarmDetailsQuery({ widgetId: this.widget._id, id: this.alarm._id });
  },
  methods: {
    refreshTabs() {
      this.tabsKey = uid();
    },

    fetchList() {
      return this.fetchAlarmDetails({
        widgetId: this.widget._id,
        id: this.alarm._id,
        query: convertAlarmDetailsQueryToRequest(this.query),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
@media (min-width: 0) {
  .xs0 {
    max-width: 0;
    max-height: 0;
    overflow: hidden;
  }
}
</style>

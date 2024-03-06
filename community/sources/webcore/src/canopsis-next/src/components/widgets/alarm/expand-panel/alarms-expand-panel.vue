<template lang="pug">
  v-tabs.expand-panel.secondary.lighten-2(
    v-model="activeTab",
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    v-tab(
      v-if="hasMoreInfos",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos}`",
      :class="moreInfosTabClass"
      ) {{ $t('alarm.tabs.moreInfos') }}
    v-tab(
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.timeLine}`",
      :class="timeLineTabClass"
    ) {{ $t('alarm.tabs.timeLine') }}
    v-tab(
      v-if="hasTickets",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared}`"
    ) {{ $t('alarm.tabs.ticketsDeclared') }}
    v-tab(:href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior}`") {{ $tc('common.pbehavior', 2) }}
    v-tab(
      v-if="hasChildren",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren}`"
    ) {{ $t('alarm.tabs.alarmsChildren') }}
    v-tab(
      v-if="hasServiceDependencies",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.trackSource}`"
    ) {{ $t('alarm.tabs.trackSource') }}
    v-tab(
      v-if="hasImpactsDependencies",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.impactChain}`"
    ) {{ $t('alarm.tabs.impactChain') }}
    v-tab(
      v-if="hasEntityGantt",
      :href="`#${$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt}`"
    ) {{ $t('alarm.tabs.entityGantt') }}

    v-tabs-items(v-model="activeTab")
      v-tab-item(v-if="hasMoreInfos", :value="$constants.ALARMS_EXPAND_PANEL_TABS.moreInfos")
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                alarms-expand-panel-more-infos(
                  :alarm="alarm",
                  :template="widget.parameters.moreInfoTemplate"
                )
      v-tab-item(:value="$constants.ALARMS_EXPAND_PANEL_TABS.timeLine")
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-progress-linear(
                :active="pending",
                :height="3",
                indeterminate
              )
              v-card-text
                alarms-time-line(
                  :steps="steps",
                  :is-html-enabled="isHtmlEnabled",
                  @update:page="updateStepsQueryPage"
                )
      v-tab-item(v-if="hasTickets", :value="$constants.ALARMS_EXPAND_PANEL_TABS.ticketsDeclared")
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                declared-tickets-list(:tickets="alarm.v.tickets", :parent-alarm-id="parentAlarmId")
      v-tab-item(:value="$constants.ALARMS_EXPAND_PANEL_TABS.pbehavior")
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                pbehaviors-simple-list(
                  :entity="alarm.entity",
                  :removable="hasDeleteAnyPbehaviorAccess",
                  :updatable="hasUpdateAnyPbehaviorAccess"
                )
      v-tab-item(v-if="hasChildren", :value="$constants.ALARMS_EXPAND_PANEL_TABS.alarmsChildren")
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                alarms-expand-panel-children(
                  :children="children",
                  :alarm="alarm",
                  :widget="widget",
                  :editing="editing",
                  :pending="pending",
                  :query.sync="childrenQuery",
                  :refresh-alarms-list="fetchList"
                )
      v-tab-item(v-if="hasServiceDependencies", :value="$constants.ALARMS_EXPAND_PANEL_TABS.trackSource")
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                service-dependencies(
                  :root="dependency",
                  :columns="widget.parameters.serviceDependenciesColumns",
                  include-root,
                  openable-root
                )
      v-tab-item(v-if="hasImpactsDependencies", :value="$constants.ALARMS_EXPAND_PANEL_TABS.impactChain")
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                service-dependencies(
                  :root="dependency",
                  :columns="widget.parameters.serviceDependenciesColumns",
                  include-root,
                  impact,
                  openable-root
                )
      v-tab-item(v-if="hasEntityGantt", :value="$constants.ALARMS_EXPAND_PANEL_TABS.entityGantt", lazy)
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                entity-gantt(:alarm="alarm")
</template>

<script>
import { isEqual, map } from 'lodash';

import {
  ENTITY_TYPES,
  GRID_SIZES,
  TOURS,
  JUNIT_ALARM_CONNECTOR,
} from '@/constants';

import uid from '@/helpers/uid';
import { setField } from '@/helpers/immutable';
import { getStepClass } from '@/helpers/tour';
import { alarmToServiceDependency } from '@/helpers/treeview/service-dependencies';
import { convertAlarmDetailsQueryToRequest } from '@/helpers/query';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { widgetExpandPanelAlarmDetails } from '@/mixins/widget/expand-panel/alarm/details';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';
import DeclaredTicketsList from '@/components/other/declare-ticket/declared-tickets-list.vue';

import AlarmsTimeLine from '../time-line/alarms-time-line.vue';
import EntityGantt from '../entity-gantt/entity-gantt.vue';
import AlarmsExpandPanelMoreInfos from './alarms-expand-panel-more-infos.vue';
import AlarmsExpandPanelChildren from './alarms-expand-panel-children.vue';

export default {
  components: {
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
    editing: {
      type: Boolean,
      default: false,
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    isTourEnabled: {
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
    moreInfosTabClass() {
      if (this.isTourEnabled) {
        return getStepClass(TOURS.alarmsExpandPanel, 2);
      }

      return '';
    },

    timeLineTabClass() {
      if (this.isTourEnabled) {
        return getStepClass(TOURS.alarmsExpandPanel, 3);
      }

      return '';
    },

    cardFlexClass() {
      const { expandGridRangeSize: [start, end] = [GRID_SIZES.min, GRID_SIZES.max] } = this.widget.parameters;

      return [
        `offset-xs${start}`,
        `xs${end - start}`,
      ];
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
      return this.widget.parameters.moreInfoTemplate ?? this.isTourEnabled;
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

    'widget.parameters.widgetGroupColumns': {
      handler(columns) {
        this.query = {
          ...this.query,

          search_by: map(columns, 'value'),
        };
      },
    },

    isTourEnabled() {
      this.refreshTabs();
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
  .tab-item-card {
    margin: auto;
  }

  @media (min-width: 0) {
    .xs0 {
      max-width: 0;
      max-height: 0;
      overflow: hidden;
    }
  }
</style>

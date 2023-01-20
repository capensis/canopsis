<template lang="pug">
  v-tabs.expand-panel.secondary.lighten-2(
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    template(v-if="hasMoreInfos")
      v-tab(:class="moreInfosTabClass") {{ $t('alarm.tabs.moreInfos') }}
      v-tab-item
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                alarms-expand-panel-more-infos(
                  :alarm="alarm",
                  :template="widget.parameters.moreInfoTemplate"
                )
    v-tab(:class="timeLineTabClass") {{ $t('alarm.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-progress-linear(
              :active="pending",
              :height="3",
              indeterminate
            )
            v-card-text
              time-line(
                :steps="steps",
                :is-html-enabled="isHtmlEnabled",
                @update:page="updateStepsQueryPage"
              )
    v-tab {{ $tc('common.pbehavior', 2) }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              pbehaviors-simple-list(
                :entity="alarm.entity",
                :removable="hasDeleteAnyPbehaviorAccess",
                :updatable="hasUpdateAnyPbehaviorAccess"
              )
    template(v-if="hasChildren")
      v-tab {{ $t('alarm.tabs.alarmsChildren') }}
      v-tab-item
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
                  :query.sync="childrenQuery"
                )
    template(v-if="hasServiceDependencies")
      v-tab {{ $t('alarm.tabs.trackSource') }}
      v-tab-item
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
    template(v-if="hasImpactsDependencies")
      v-tab {{ $t('alarm.tabs.impactChain') }}
      v-tab-item
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
    template(v-if="hasEntityGantt")
      v-tab {{ $t('alarm.tabs.entityGantt') }}
      v-tab-item(lazy)
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                entity-gantt(:alarm="alarm")
</template>

<script>
import { isEqual } from 'lodash';

import {
  ENTITY_TYPES,
  GRID_SIZES,
  TOURS,
  JUNIT_ALARM_CONNECTOR,
} from '@/constants';

import uid from '@/helpers/uid';
import { getStepClass } from '@/helpers/tour';
import { alarmToServiceDependency } from '@/helpers/treeview/service-dependencies';
import { convertAlarmDetailsQueryToRequest } from '@/helpers/query';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { widgetExpandPanelAlarmDetails } from '@/mixins/widget/expand-panel/alarm/details';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';

import TimeLine from '../time-line/time-line.vue';
import EntityGantt from '../entity-gantt/entity-gantt.vue';
import AlarmsExpandPanelMoreInfos from './alarms-expand-panel-more-infos.vue';
import AlarmsExpandPanelChildren from './alarms-expand-panel-children.vue';

export default {
  components: {
    PbehaviorsSimpleList,
    ServiceDependencies,
    TimeLine,
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
  },
  data() {
    return {
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
      return alarmToServiceDependency(this.alarm);
    },

    hasMoreInfos() {
      return this.widget.parameters.moreInfoTemplate ?? this.isTourEnabled;
    },

    hasChildren() {
      return this.alarm.children && !this.hideChildren;
    },

    hasServiceDependencies() {
      return this.alarm.entity.type === ENTITY_TYPES.service;
    },

    hasImpactsDependencies() {
      const { impact } = this.alarm.entity;

      return this.hasServiceDependencies
        ? impact?.length > 0
        // resource and component types having one basic entity into impact
        : impact?.length > 1;
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

    isTourEnabled() {
      this.refreshTabs();
    },

    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
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

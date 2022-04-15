<template lang="pug">
  v-tabs.expand-panel.secondary.lighten-2(
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    template(v-if="hasMoreInfos")
      v-tab(:class="moreInfosTabClass") {{ $t('alarmList.tabs.moreInfos') }}
      v-tab-item
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
    v-tab(:class="timeLineTabClass") {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              time-line(
                :alarm="alarm",
                :widget="widget",
                :is-html-enabled="isHtmlEnabled",
                :hide-groups="hideGroups"
              )
    v-tab {{ $tc('common.pbehavior', 2) }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              pbehaviors-simple-list(
                :entity="alarm.entity",
                :deletable="hasDeleteAnyPbehaviorAccess",
                :editable="hasUpdateAnyPbehaviorAccess"
              )
    template(v-if="hasCauses")
      v-tab {{ $t('alarmList.tabs.alarmsCauses') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :default-query-id="causesKey",
                  :tab-id="causesKey",
                  :alarm="alarm",
                  :is-editing-mode="isEditingMode"
                )
    template(v-if="hasConsequences")
      v-tab {{ $t('alarmList.tabs.alarmsConsequences') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :default-query-id="consequencesKey",
                  :tab-id="consequencesKey",
                  :alarm="alarm",
                  :is-editing-mode="isEditingMode"
                )
    template(v-if="hasServiceDependencies")
      v-tab {{ $t('alarmList.tabs.trackSource') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                service-dependencies(
                  :root="dependency",
                  :columns="widget.parameters.serviceDependenciesColumns",
                  include-root
                )
    template(v-if="hasImpactsDependencies")
      v-tab {{ $t('alarmList.tabs.impactChain') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                service-dependencies(
                  :root="dependency",
                  :columns="widget.parameters.serviceDependenciesColumns",
                  include-root,
                  impact
                )
    template(v-if="hasEntityGantt")
      v-tab {{ $t('alarmList.tabs.entityGantt') }}
      v-tab-item(lazy)
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                entity-gantt(:alarm="alarm")
</template>

<script>
import {
  ALARMS_GROUP_PREFIX,
  ENTITY_TYPES,
  GRID_SIZES,
  TOURS,
  JUNIT_ALARM_CONNECTOR,
} from '@/constants';

import uid from '@/helpers/uid';
import { getStepClass } from '@/helpers/tour';
import { serviceToServiceDependency } from '@/helpers/treeview/service-dependencies';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import ServiceDependencies from '@/components/other/service/table/service-dependencies.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/partials/pbehaviors-simple-list.vue';

import TimeLine from '../time-line/time-line.vue';
import MoreInfos from '../more-infos/more-infos.vue';
import GroupAlarmsList from '../group-alarms-list.vue';
import EntityGantt from '../entity-gantt/entity-gantt.vue';

export default {
  components: {
    PbehaviorsSimpleList,
    ServiceDependencies,
    TimeLine,
    MoreInfos,
    GroupAlarmsList,
    EntityGantt,
  },
  mixins: [
    entitiesInfoMixin,
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
    isEditingMode: {
      type: Boolean,
      default: false,
    },
    hideGroups: {
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
    causesKey() {
      return `${ALARMS_GROUP_PREFIX.CAUSES}${this.alarm._id}`;
    },

    consequencesKey() {
      return `${ALARMS_GROUP_PREFIX.CONSEQUENCES}${this.alarm._id}`;
    },

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
      return serviceToServiceDependency(this.alarm.entity, this.alarm);
    },

    hasMoreInfos() {
      return this.widget.parameters.moreInfoTemplate || this.isTourEnabled;
    },

    hasCauses() {
      return this.alarm.causes && !this.hideGroups;
    },

    hasConsequences() {
      return this.alarm.consequences && !this.hideGroups;
    },

    hasServiceDependencies() {
      return this.alarm.entity.type === ENTITY_TYPES.service;
    },

    hasImpactsDependencies() {
      return this.hasServiceDependencies
        ? this.alarm.entity.impact.length > 0
        // resource and component types having one basic entity into impact
        : this.alarm.entity.impact.length > 1;
    },

    hasEntityGantt() {
      /**
       * We have junit feature only on `cat` version of canopsis
       */
      return this.isCatVersion
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

    isTourEnabled() {
      this.refreshTabs();
    },
  },
  methods: {
    refreshTabs() {
      this.tabsKey = uid();
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

<template lang="pug">
  v-tabs.expand-panel.secondary.lighten-2(
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    template(v-if="widget.parameters.moreInfoTemplate || isTourEnabled")
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
                :isHTMLEnabled="isHTMLEnabled",
                :hideGroups="hideGroups"
              )
    template(v-if="alarm.causes && !hideGroups")
      v-tab {{ $t('alarmList.tabs.alarmsCauses') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :defaultQueryId="causesKey",
                  :tabId="causesKey",
                  :alarm="alarm",
                  :isEditingMode="isEditingMode"
                )
    template(v-if="alarm.consequences && !hideGroups")
      v-tab {{ $t('alarmList.tabs.alarmsConsequences') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :defaultQueryId="consequencesKey",
                  :tabId="consequencesKey",
                  :alarm="alarm",
                  :isEditingMode="isEditingMode"
                )
</template>

<script>
import { ALARMS_GROUP_PREFIX, GRID_SIZES, TOURS } from '@/constants';

import uid from '@/helpers/uid';
import { getStepClass } from '@/helpers/tour';

import TimeLine from '../time-line/time-line.vue';
import MoreInfos from '../more-infos/more-infos.vue';
import GroupAlarmsList from '../group-alarms-list.vue';

export default {
  components: {
    GroupAlarmsList,
    TimeLine,
    MoreInfos,
  },
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

    isHTMLEnabled() {
      return this.widget.parameters.isHtmlEnabledOnTimeLine;
    },
  },
  watch: {
    'widget.parameters.moreInfoTemplate': {
      handler() {
        this.refreshTabs();
      },
    },

    isTourEnabled(value, oldValue) {
      if (value !== oldValue) {
        this.refreshTabs();
      }
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

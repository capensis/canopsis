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
              time-line(:alarm="alarm", :isHTMLEnabled="isHTMLEnabled")
                // TODO replace causes after api will finish
    template(v-if="alarm.causes")
      v-tab {{ $t('alarmList.tabs.alarmsCauses') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex
            v-card.tab-item-card
              v-card-text
                // TODO replace causes after api will finish
                group-alarms-list(:widget="widget", :alarms="alarm.causes", :isEditingMode="isEditingMode")
    // TODO replace consequences after api will finish
    template(v-if="alarm.consequences")
      v-tab {{ $t('alarmList.tabs.alarmsConsequences') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex
            v-card.tab-item-card
              v-card-text
                // TODO replace consequences after api will finish
                group-alarms-list(:widget="widget", :alarms="alarm.consequences", :isEditingMode="isEditingMode")
</template>

<script>
import { GRID_SIZES, TOURS } from '@/constants';

import uid from '@/helpers/uid';
import { getStepClass } from '@/helpers/tour';

import TimeLine from '@/components/other/alarm/time-line/time-line.vue';
import MoreInfos from '@/components/other/alarm/more-infos/more-infos.vue';
import GroupAlarmsList from '@/components/other/alarm/group-alarms-list.vue';

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
</style>

<template lang="pug">
  v-tabs.expand-panel(
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    template(v-if="widget.parameters.moreInfoTemplate")
      v-tab {{ $t('alarmList.tabs.moreInfos') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
    v-tab {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              time-line(:alarm="alarm", :isHTMLEnabled="isHTMLEnabled")
    template(v-if="true")
      v-tab {{ $t('alarmList.tabs.alarmsCauses') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :details="alarm.causes",
                  :defaultQueryId="`causes_${alarm._id}`",
                  :tabId="`causes_${alarm._id}`",
                  :alarmId="alarm._id",
                  :isEditingMode="isEditingMode"
                )
    template(v-if="true")
      v-tab {{ $t('alarmList.tabs.alarmsConsequences') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex
            v-card.tab-item-card
              v-card-text
                group-alarms-list(
                  :widget="widget",
                  :details="alarm.consequences",
                  :defaultQueryId="`consequences_${alarm._id}`",
                  :tabId="`causes_${alarm._id}`",
                  :alarmId="alarm._id",
                  :isEditingMode="isEditingMode"
                )
</template>

<script>
import { GRID_SIZES } from '@/constants';

import uid from '@/helpers/uid';

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
  },
  data() {
    return {
      tabsKey: uid(),
    };
  },
  computed: {
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
        this.tabsKey = uid();
      },
    },
  },
};
</script>

<style lang="scss" scoped>
  .tab-item-card {
    margin: auto;
  }
</style>

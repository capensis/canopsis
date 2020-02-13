<template lang="pug">
  v-tabs.expand-panel(
    :key="tabsKey",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    template(v-if="widget.parameters.moreInfoTemplate || isTourEnabled")
      v-tab(:class="{ 'v-tour-step-2': isTourEnabled }") {{ $t('alarmList.tabs.moreInfos') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
    v-tab(:class="{ 'v-tour-step-3': isTourEnabled }") {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              time-line(:alarm="alarm", :isHTMLEnabled="isHTMLEnabled")
</template>

<script>
import { GRID_SIZES } from '@/constants';

import uid from '@/helpers/uid';

import TimeLine from '@/components/other/alarm/time-line/time-line.vue';
import MoreInfos from '@/components/other/alarm/more-infos/more-infos.vue';

export default {
  components: {
    TimeLine,
    MoreInfos,
  },
  props: {
    isHTMLEnabled: {
      type: Boolean,
      default: false,
    },
    alarm: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
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
    cardFlexClass() {
      const { expandGridRangeSize: [start, end] = [GRID_SIZES.min, GRID_SIZES.max] } = this.widget.parameters;

      return [
        `offset-xs${start}`,
        `xs${end - start}`,
      ];
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

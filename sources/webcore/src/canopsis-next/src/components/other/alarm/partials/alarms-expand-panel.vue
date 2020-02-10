<template lang="pug">
  v-tabs.expand-panel(
    ref="tabs",
    color="secondary lighten-1",
    slider-color="primary",
    dark,
    centered
  )
    v-tab {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(:class="cardFlexClass")
          v-card.tab-item-card
            v-card-text
              time-line(:alarm="alarm", :isHTMLEnabled="isHTMLEnabled")
    template(v-if="widget.parameters.moreInfoTemplate")
      v-tab {{ $t('alarmList.tabs.moreInfos') }}
      v-tab-item
        v-layout.pa-3.secondary.lighten-2(row)
          v-flex(:class="cardFlexClass")
            v-card.tab-item-card
              v-card-text
                more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
</template>

<script>
import { GRID_SIZES } from '@/constants';

import TimeLine from '@/components/other/alarm/time-line/time-line.vue';
import MoreInfos from '@/components/other/alarm/more-infos/more-infos.vue';

import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

export default {
  components: {
    TimeLine,
    MoreInfos,
  },
  mixins: [vuetifyTabsMixin],
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
        this.$nextTick(this.callTabsUpdateTabsMethod);
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

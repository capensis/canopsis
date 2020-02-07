<template lang="pug">
  v-tabs.expand-panel(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t('alarmList.tabs.moreInfos') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(xs12)
          v-card.tab-item-card(:width="expandCardWidthPercentage")
            v-card-text
              more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
    v-tab {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
      v-layout.pa-3.secondary.lighten-2(row)
        v-flex(xs12)
          v-card.tab-item-card(:width="expandCardWidthPercentage")
            v-card-text
              time-line(:alarm="alarm", :isHTMLEnabled="isHTMLEnabled")
</template>

<script>
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
  },
  computed: {
    expandCardWidthPercentage() {
      return `${this.widget.parameters.expandCardWidthPercentage || 100}%`;
    },
  },
};
</script>

<style lang="scss" scoped>
  .tab-item-card {
    margin: auto;
  }
</style>

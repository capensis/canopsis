<template lang="pug">
  v-tabs.expand-panel
    v-tab {{ $t('alarmList.tabs.moreInfos') }}
    v-tab-item
      more-infos(:alarm="alarm", :template="widget.parameters.moreInfoTemplate")
    v-tab {{ $t('alarmList.tabs.timeLine') }}
    v-tab-item
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
};
</script>

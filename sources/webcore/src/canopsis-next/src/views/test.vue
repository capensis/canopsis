<template lang="pug">
  stats-text(:template="template")
</template>

<script>
import StatsText from '@/components/other/stats/text/stats-text.vue';

export default {
  components: { StatsText },
  data() {
    return {
      template: `<h1 style="color: red">Alarms created/resolved last 7 days</h1>
        <ul>
          <li>Alarms resolved: {{ stat name="Alarms resolved" }}</li>
          <li>Alarms created: {{ stat name="Alarms created" }}</li>
        </ul>`,
    };
  },
};
</script>

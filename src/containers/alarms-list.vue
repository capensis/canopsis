<template lang="pug">
  div
    alarms-list(
    :widget="widget",
    )
</template>

<script>
import AlarmsList from '@/components/other/alarm/alarms-list.vue';

export default {
  components: {
    AlarmsList,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
};
</script>

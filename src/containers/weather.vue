<template lang="pug">
  weather(
    :widget="widget",
    @openSettings="$emit('openSettings', $event)"
  )
</template>

<script>
import Weather from '@/components/other/service-weather/weather.vue';

export default {
  components: {
    Weather,
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
};
</script>

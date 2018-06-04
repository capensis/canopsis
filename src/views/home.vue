<template lang="pug">
  div#brand Canopsis Next
</template>

<script>
import AlarmsListSettings from '@/components/alarms-list-settings.vue';

export default {
  name: 'home',
  components: {
    AlarmsListSettings,
  },
};
</script>

<style lang="scss" scoped>
  #brand {
    text-align: center;
    position: relative;
    top: 25%;
    max-width: 50%;
    max-height: 5em;
    margin: auto;
    font-weight: bold;
    font-size: 2em;
  }
</style>

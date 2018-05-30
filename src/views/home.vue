<template lang="pug">
  div
    alarms-list-settings
    m-filter-editor
</template>

<script>
import AlarmsListSettings from '@/components/alarms-list-settings.vue';
import MFilterEditor from '@/components/mfilter-editor.vue';

export default {
  name: 'home',
  components: {
    AlarmsListSettings,
    MFilterEditor,
  },
};
</script>

<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      date-time-picker
</template>

<script>
import DateTimePicker from '@/components/forms/date-time-picker.vue';

export default {
  components: {
    DateTimePicker,
  },
  props: {
    title: {
      type: String,
    },
  },
};
</script>


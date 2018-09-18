<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      date-time-picker(:value="dateValue", @input="dateChange")
</template>

<script>
import moment from 'moment';
import DateTimePicker from '@/components/forms/date-time-picker.vue';

export default {
  components: {
    DateTimePicker,
  },
  props: {
    value: {
      type: Number,
    },
    title: {
      type: String,
    },
  },
  computed: {
    dateValue() {
      return new Date(this.value * 1000);
    },
  },
  methods: {
    dateChange(event) {
      this.$emit('input', moment(event).unix());
    },
  },
};
</script>


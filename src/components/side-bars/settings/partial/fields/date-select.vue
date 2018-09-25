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
      if (this.value) {
        return moment(this.value).toDate();
      }
      return moment().toDate();
    },
  },
  methods: {
    dateChange(event) {
      this.$emit('input', moment(event).valueOf());
    },
  },
};
</script>


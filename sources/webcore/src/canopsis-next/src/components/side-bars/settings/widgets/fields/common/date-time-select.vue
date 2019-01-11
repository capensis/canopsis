<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      date-time-picker(name="dateTime", :value="dateObject", @input="updateValue")
</template>

<script>
import moment from 'moment';
import DateTimePicker from '@/components/forms/fields/date-time-picker.vue';

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
    name: {
      type: String,
      default: 'date',
    },
  },
  computed: {
    dateObject() {
      if (this.value) {
        return moment.unix(this.value).toDate();
      }
      return moment().toDate();
    },
  },
  methods: {
    updateValue(value) {
      this.$emit('input', moment(value).unix());
    },
  },
};
</script>


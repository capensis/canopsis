<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-btn(@click="click") Select a date
</template>

<script>
import moment from 'moment';
import DateTimePicker from '@/components/forms/date-time-picker.vue';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  components: {
    DateTimePicker,
  },
  mixins: [modalMixin],
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
        return moment.unix(this.value).toDate();
      }
      return moment().toDate();
    },
  },
  methods: {
    click() {
      this.showModal({
        name: MODALS.dateSelect,
        config: {
          value: this.dateValue,
          action: newDate => this.$emit('input', moment(newDate).unix()),
        },
      });
    },
  },
};
</script>


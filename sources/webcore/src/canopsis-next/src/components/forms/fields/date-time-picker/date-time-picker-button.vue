<template lang="pug">
  v-menu(
    v-model="opened",
    content-class="date-time-picker",
    transition="slide-y-transition",
    max-width="290px",
    :close-on-content-click="false",
    right,
    lazy-with-unmount,
    lazy
  )
    v-btn(
      slot="activator",
      color="secondary",
      icon,
      fab,
      small
    )
      v-icon calendar_today
    date-time-picker(
      :value="value",
      :round-hours="roundHours",
      @input="$listeners.input",
      @close="close"
    )
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

import DateTimePicker from './date-time-picker.vue';

export default {
  components: { DateTimePicker },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Date,
      default: null,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  methods: {
    close() {
      this.opened = false;
    },
  },
};
</script>

<template lang="pug">
  v-menu(
    v-model="isMenuOpen",
    content-class="date-time-picker",
    transition="slide-y-transition",
    max-width="290px",
    :close-on-content-click="false",
    right,
    lazy
  )
    v-btn(
      data-test="dateTimePickerButton",
      slot="activator",
      color="secondary",
      icon,
      fab,
      small
    )
      v-icon calendar_today
    date-time-picker(
      data-test="dateTimePickerCalendar",
      :value="value",
      :opened="isMenuOpen",
      :roundHours="roundHours",
      :useSeconds="useSeconds",
      @input="updateModel($event)"
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
    useSeconds: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isMenuOpen: false,
    };
  },
};
</script>

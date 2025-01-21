<template>
  <v-menu
    v-model="opened"
    :close-on-content-click="false"
    content-class="date-time-picker"
    transition="slide-y-transition"
    max-width="290"
    right
  >
    <template #activator="{ on }">
      <slot :value="value" :on="on" name="activator">
        <v-btn
          icon
          fab
          small
          v-on="on"
        >
          <v-icon>calendar_today</v-icon>
        </v-btn>
      </slot>
    </template>
    <date-time-picker
      v-field="value"
      :label="label"
      :round-hours="roundHours"
      @close="close"
    />
  </v-menu>
</template>

<script>
import DateTimePicker from './date-time-picker.vue';

export default {
  components: { DateTimePicker },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Date,
      default: () => new Date(),
    },
    label: {
      type: String,
      default: '',
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

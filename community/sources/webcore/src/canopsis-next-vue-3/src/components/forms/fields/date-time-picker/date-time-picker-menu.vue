<template>
  <v-menu
    v-model="opened"
    :close-on-content-click="false"
    content-class="date-time-picker"
    transition="slide-y-transition"
    max-width="370"
    location="bottom end"
  >
    <template #activator="{ props: activatorProps }">
      <slot
        :value="value"
        :on="activatorProps"
        :props="activatorProps"
        name="activator"
      >
        <v-btn
          v-bind="activatorProps"
          variant="flat"
          icon
        >
          <v-icon>calendar_today</v-icon>
        </v-btn>
      </slot>
    </template>
    <date-time-picker
      :value="value"
      :label="label"
      :round-hours="roundHours"
      :allowed-dates="allowedDates"
      @close="close"
      @update:value="$emit('update:value', $event)"
    />
  </v-menu>
</template>

<script>
import DateTimePicker from './date-time-picker.vue';

export default {
  emits: ['update:value'],
  components: { DateTimePicker },
  props: {
    value: {
      type: [Date, Number],
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
    allowedDates: {
      type: Function,
      required: false,
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

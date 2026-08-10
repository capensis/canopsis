<template>
  <v-menu
    v-model="opened"
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="date-time-picker"
    transition="slide-y-transition"
    max-width="370"
    location="bottom end"
  >
    <template #activator="{ props: activatorProps }">
      <div v-bind="activatorProps" :aria-required="fieldRequired">
        <c-text-field
          :label="label"
          :name="name"
          :model-value="dateTextValue"
          :append-inner-icon="clearable ? 'close' : ''"
          :disabled="disabled"
          :hide-details="hideDetails"
          :required="fieldRequired"
          :validation-rules="validationRules"
          readonly
          @click:append-inner="clear"
        >
          <template v-if="$slots.append" #append>
            <slot name="append" />
          </template>
          <template v-if="$slots['append-outer']" #append-outer>
            <slot name="append-outer" />
          </template>
        </c-text-field>
      </div>
    </template>
    <date-time-picker
      :value="value"
      :label="label"
      :round-hours="roundHours"
      :allowed-dates="allowedDates"
      :control-variant="controlVariant"
      @close="close"
      @update:value="updateValue"
    />
  </v-menu>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';
import { convertDateToString } from '@/helpers/date/date';

import DateTimePicker from './date-time-picker.vue';
import { createDateTimeValidationRules, isRequiredRule } from './validation-rules';

export default {
  components: { DateTimePicker },
  emits: ['update:value'],
  props: {
    clearable: {
      type: Boolean,
      default: false,
    },
    value: {
      type: [Date, Number],
      default: () => new Date(),
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: '',
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    controlVariant: {
      type: String,
      default: undefined,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    allowedDates: {
      type: Function,
      required: false,
    },
    rules: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    dateTextValue() {
      return convertDateToString(this.value, DATETIME_FORMATS.dateTimePicker);
    },

    fieldRequired() {
      return this.required || isRequiredRule(this.rules);
    },

    validationRules() {
      return createDateTimeValidationRules(this.rules);
    },
  },
  methods: {
    close() {
      this.opened = false;
    },

    updateValue(value) {
      this.$emit('update:value', value);
    },

    clear() {
      this.updateValue(null);
    },
  },
};
</script>

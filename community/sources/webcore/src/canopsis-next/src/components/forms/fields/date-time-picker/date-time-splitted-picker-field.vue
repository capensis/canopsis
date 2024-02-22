<template>
  <v-layout
    class="date-time-splitted-field"
    column
  >
    <v-layout :reverse="reverse">
      <v-flex class="date-time-splitted-field__date">
        <c-date-picker-field
          :value="value | date('vuetifyDatePicker', null)"
          :label="!reverse || fullDay ? label : ''"
          :placeholder="placeholder"
          :error="errors.has(name)"
          :disabled="disabled"
          :min="min"
          :max="max"
          hide-details
          @input="updateDate"
        />
      </v-flex>
      <time-picker-field
        v-if="!fullDay"
        :value="value | date('timePicker', null)"
        :label="reverse ? label : ''"
        :error="errors.has(name)"
        :disabled="disabled"
        class="date-time-splitted-field__time"
        hide-details
        @input="updateTime"
      />
    </v-layout>
    <div class="v-text-field__details mt-2">
      <v-messages
        :value="errors.collect(name)"
        color="error"
      />
    </div>
  </v-layout>
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitDateObject } from '@/helpers/date/date';
import { getDateObjectByTime, getDateObjectByDate } from '@/helpers/date/date-time-picker';

import { formBaseMixin } from '@/mixins/form';

import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

export default {
  $_veeValidate: {
    value() {
      if (!this.value) {
        return this.value;
      }

      return convertDateToStartOfUnitDateObject(this.value, TIME_UNITS.minute);
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: { TimePickerField },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Date,
      default: null,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'date',
    },
    reverse: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: String,
      required: false,
    },
    max: {
      type: String,
      required: false,
    },
  },
  methods: {
    updateTime(time) {
      this.updateModel(getDateObjectByTime(this.value, time));
    },

    updateDate(date) {
      this.updateModel(getDateObjectByDate(this.value, date));
    },
  },
};
</script>

<style lang="scss">
.date-time-splitted-field {
  & > .layout {
    gap: 4px;
  }

  &__date {
    flex: 1;
    min-width: 90px;
  }

  &__time {
    width: 50px;
    max-width: 50px;
  }
}
</style>

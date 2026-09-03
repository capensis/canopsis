<template>
  <v-sheet class="date-time-picker__sheet">
    <div v-if="label" class="v-picker__title bg-primary text-center">
      <span class="text-h5">
        {{ label }}
      </span>
    </div>
    <div class="date-time-picker__body">
      <c-layout class="py-2" align="start" justify="space-around">
        <div class="v-date-time-picker__subtitle-wrapper">
          <span
            :class="{ 'text-medium-emphasis': !localValue }"
            class="v-date-time-picker__subtitle"
          >
            {{ valueString }}
          </span>
        </div>
        <div class="v-date-time-picker__subtitle-wrapper">
          <time-picker-field
            :value="timeString"
            :round-hours="roundHours"
            :allowed-time="allowedTime"
            class="v-date-time-picker__subtitle"
            @update:value="updateTime"
          />
        </div>
      </c-layout>
      <div>
        <v-date-picker
          :model-value="dateValue"
          :allowed-dates="allowedDates"
          :control-variant="controlVariant"
          color="primary"
          first-day-of-week="1"
          hide-header
          @update:model-value="updateDate"
        />
      </div>
    </div>
    <slot name="footer">
      <v-divider />
      <c-layout class="mt-1 py-2" justify="space-around">
        <v-btn type="button" variant="text" @click="$emit('close')">
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn type="button" color="primary" @click="submit">
          {{ $t('common.apply') }}
        </v-btn>
      </c-layout>
    </slot>
  </v-sheet>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToDateObject, convertDateToString } from '@/helpers/date/date';
import {
  getDateObjectByDate,
  getDateObjectByTime,
  normalizeTimeString,
} from '@/helpers/date/date-time-picker';

import { formBaseMixin } from '@/mixins/form';

import TimePickerField from '../time-picker/time-picker-field.vue';

export default {
  components: { TimePickerField },
  mixins: [formBaseMixin],
  emits: ['close', 'update:value'],
  props: {
    value: {
      type: [Date, Number],
      default: () => new Date(),
    },
    label: {
      type: String,
      default: '',
    },
    dateFormat: {
      type: String,
      default: DATETIME_FORMATS.short,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    controlVariant: {
      type: String,
      default: undefined,
    },
    allowedDates: {
      type: Function,
      required: false,
    },
  },
  data() {
    return {
      localValue: this.value ? convertDateToDateObject(this.value) : null,
    };
  },
  computed: {
    valueString() {
      return convertDateToString(this.localValue, this.dateFormat, '--/--/----');
    },

    timeString() {
      return convertDateToString(this.localValue, DATETIME_FORMATS.timePicker, null);
    },

    dateValue() {
      return this.localValue ? convertDateToDateObject(this.localValue) : null;
    },

    allowedTime() {
      if (!this.allowedDates || !this.value) {
        return () => true;
      }

      return (time) => {
        const [hours, minutes] = normalizeTimeString(time).split(':');
        const currentDate = convertDateToDateObject(this.value);

        currentDate.setHours(hours, minutes);

        return this.allowedDates?.(currentDate);
      };
    },
  },
  watch: {
    value() {
      if (this.value !== this.localValue) {
        this.localValue = this.value ? convertDateToDateObject(this.value) : null;
      }
    },
  },
  methods: {
    updateTime(time) {
      this.localValue = getDateObjectByTime(this.localValue, time);
    },

    updateDate(date) {
      const preparedDate =
        date instanceof Date ? convertDateToString(date, DATETIME_FORMATS.datePicker) : date;

      this.localValue = getDateObjectByDate(this.localValue, preparedDate);
    },

    submit() {
      this.updateModel(this.localValue);

      this.$emit('close');
    },
  },
};
</script>

<style lang="scss">
.date-time-picker {
  &__sheet {
    width: 360px;
    max-width: 360px;
    overflow-x: hidden;

    .v-field__input {
      padding: 0;
    }
  }

  .v-date-picker {
    width: 360px;
  }

  .v-date-picker-controls {
    width: 100%;
    padding: 0;
  }

  .v-date-picker-month {
    width: 100%;
  }

  .v-date-picker-month__day {
    height: 36px;
  }

  .v-date-picker-month__day .v-btn.v-date-picker-month__day-btn {
    --v-btn-height: 32px;
    --v-btn-size: 0.875rem;

    width: 32px;
    min-width: 32px;
    height: 32px;
  }

  .v-picker__title {
    min-height: 64px;
    padding: 16px;

    .text-h5 {
      font-size: 1.5rem !important;
      font-weight: 400;
      line-height: 2rem;
      letter-spacing: 0;
    }
  }

  .date-time-picker__body {
    position: relative;
    width: 360px;
    min-height: 352px;
    z-index: inherit;
  }

  .v-date-time-picker__subtitle {
    line-height: 30px;
    font-size: 18px;
    font-weight: 400;

    &-wrapper {
      text-align: center;
    }
  }

  .v-card {
    box-shadow: none;
  }
}
</style>

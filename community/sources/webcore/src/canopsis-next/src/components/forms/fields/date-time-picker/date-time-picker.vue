<template>
  <v-sheet>
    <div
      v-if="label"
      class="v-picker__title primary text-center"
    >
      <span class="text-h5">{{ label }}</span>
    </div>
    <div class="date-time-picker__body">
      <v-layout
        class="py-2"
        align-center
        justify-center
      >
        <v-flex class="v-date-time-picker__subtitle-wrapper">
          <span
            :class="{ 'grey--text darken-1': !localValue }"
            class="v-date-time-picker__subtitle"
          >{{ valueString }}</span>
        </v-flex>
        <v-flex class="v-date-time-picker__subtitle-wrapper">
          <time-picker-field
            :value="timeString"
            :round-hours="roundHours"
            class="v-date-time-picker__subtitle"
            @input="updateTime"
          />
        </v-flex>
      </v-layout>
      <div>
        <v-date-picker
          :locale="$i18n.locale"
          :value="dateString"
          color="primary"
          first-day-of-week="1"
          no-title
          @input="updateDate"
        />
      </div>
    </div>
    <slot name="footer">
      <v-divider />
      <v-layout
        class="mt-1 py-2"
        justify-space-around
      >
        <v-btn
          depressed
          text
          @click="$listeners.close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          @click="submit"
        >
          {{ $t('common.apply') }}
        </v-btn>
      </v-layout>
    </slot>
  </v-sheet>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { getDateObjectByDate, getDateObjectByTime } from '@/helpers/date/date-time-picker';
import { convertDateToDateObject, convertDateToString } from '@/helpers/date/date';

import { formBaseMixin } from '@/mixins/form';

import TimePickerField from '../time-picker/time-picker-field.vue';

export default {
  components: { TimePickerField },
  mixins: [formBaseMixin],
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
  },
  data() {
    return {
      localValue: this.value ? convertDateToDateObject(this.value) : null,
    };
  },
  computed: {
    valueString() {
      return convertDateToString(this.localValue, this.dateFormat, '−−/−−/−−−−');
    },

    timeString() {
      return convertDateToString(this.localValue, DATETIME_FORMATS.timePicker, null);
    },

    dateString() {
      return convertDateToString(this.localValue, DATETIME_FORMATS.datePicker, null);
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
      this.localValue = getDateObjectByDate(this.localValue, date);
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
    .date-time-picker__body {
      position: relative;
      width: 290px;
      height: 352px;
      z-index: inherit;
    }

    .v-date-time-picker__subtitle {
      margin-top: -12px;
      line-height: 30px;
      font-size: 18px;
      font-weight: 400;

      &-wrapper {
        text-align: center;
      }
    }

    .v-menu__content {
      max-width: 100%;
    }

    .v-date-picker-table {
      height: 260px;
    }

    .v-card {
      box-shadow: none;
    }

    .v-date-picker-table--date .v-btn {
      height: 35px;
      width: 35px;
    }
  }
</style>

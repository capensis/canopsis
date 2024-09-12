<template lang="pug">
  v-sheet
    div.v-picker__title.primary.text-xs-center(v-if="label")
      span.headline {{ label }}
    div.date-time-picker__body
      v-layout.py-2(row, align-center, justify-center)
        v-flex.v-date-time-picker__subtitle-wrapper
          span.v-date-time-picker__subtitle(
            :class="{ 'grey--text darken-1': !localValue }"
          ) {{ valueString }}
        v-flex.v-date-time-picker__subtitle-wrapper
          time-picker-field.v-date-time-picker__subtitle(
            :value="timeString",
            :round-hours="roundHours",
            @input="updateTime"
          )
      div
        v-date-picker(
          :locale="$i18n.locale",
          :value="dateString",
          color="primary",
          first-day-of-week="1",
          no-title,
          @input="updateDate"
        )
    slot(name="footer")
      v-divider
      v-layout.mt-1(justify-space-around)
        v-btn(depressed, flat, @click="$listeners.close") {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.apply') }}
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

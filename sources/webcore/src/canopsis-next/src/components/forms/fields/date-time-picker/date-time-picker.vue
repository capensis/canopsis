<template lang="pug">
  div
    div.v-picker__title.primary.text-xs-center
      span.headline {{ label }}
    div.date-time-picker__body
      v-layout.py-2(row, align-center, justify-center)
        v-flex.v-date-time-picker__subtitle-wrapper
          span.v-date-time-picker__subtitle {{ localValue | date(dateFormat, true, null) }}
        v-flex.v-date-time-picker__subtitle-wrapper
          time-picker.v-date-time-picker__subtitle(
            :value="localValue | date('timePicker', true, null)",
            :round-hours="roundHours",
            @input="updateTime"
          )
      div
        v-date-picker(
          :locale="$i18n.locale",
          :value="localValue | date('YYYY-MM-DD', true, null)",
          color="primary",
          no-title,
          @input="updateDate"
        )
    slot(name="footer", @submit="submit")
      v-divider
      v-layout.mt-1(justify-space-around)
        v-btn(depressed, flat, @click="$listeners.close") {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.apply') }}
</template>

<script>
import dateTimePickerMixin from '@/mixins/vuetify/date-time-picker';

import TimePicker from '../time-picker/time-picker.vue';

/**
 * Date time picker component
 *
 * @prop {Date} [value=null] - Date value
 * @prop {Boolean} [roundHours=false] - Deny to change minutes it will be only 0
 *
 * @event value#input
 */
export default {
  components: { TimePicker },
  mixins: [dateTimePickerMixin],
  props: {
    value: {
      type: [Date, Number],
      default: null,
    },
    label: {
      type: String,
      default: 'End date time',
    },
    dateFormat: {
      type: String,
      default: 'short',
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      localValue: new Date(this.value),
    };
  },
  methods: {
    updateTime(time = '00:00:00') {
      const value = this.localValue;
      const newValue = new Date(value ? value.getTime() : null);
      const [hours = 0, minutes = 0, seconds = 0] = time.split(':');

      newValue.setHours(parseInt(hours, 10) || 0, parseInt(minutes, 10) || 0, parseInt(seconds, 10) || 0, 0);

      this.localValue = newValue;
    },

    updateDate(date) {
      const value = this.localValue;
      const newValue = new Date(value ? value.getTime() : null);
      const [year, month, day] = date.split('-');

      newValue.setFullYear(parseInt(year, 10), parseInt(month, 10) - 1, parseInt(day, 10));

      if (!value) {
        newValue.setHours(0, 0, 0, 0);
      } else {
        newValue.setSeconds(0, 0);
      }

      this.localValue = newValue;
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

      .v-picker__body,
      .v-time-picker-clock__item,
      .v-time-picker-clock__item span,
      .v-time-picker-clock__hand {
        z-index: inherit;
      }
    }

    .v-date-time-picker__subtitle {
      line-height: 30px;
      font-size: 18px;
      font-weight: 400;

      &-wrapper {
        text-align: center;
      }
    }

    .v-tabs__container--centered .v-tabs__div,
    .v-tabs__container--fixed-tabs .v-tabs__div,
    .v-tabs__container--icons-and-text .v-tabs__div {
      min-width: 145px;
    }

    .v-menu__content {
      max-width: 100%;
    }

    .v-dropdown-footer, &.v-menu__content, .v-tabs__items {
      background-color: #fff;
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

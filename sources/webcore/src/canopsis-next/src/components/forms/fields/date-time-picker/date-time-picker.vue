<template lang="pug">
  div
    .v-picker__title.primary.text-xs-center
      span.v-date-time-picker-title(:class="{ 'use-seconds': useSeconds }")
        span.v-picker__title__btn(
          @click="showDateTab",
          :class="{ 'v-picker__title__btn--active': isActiveDateTab }"
        ) {{ value | date('datePicker', true, '--/--/----') }}
        span &nbsp;
        span.v-picker__title__btn(
          :class="{ 'v-picker__title__btn--active': isActiveHoursTab }",
          @click="showHoursTabInTimeTab"
        ) {{ value | date('HH', true, '--') }}
        span :
        span.v-picker__title__btn(
          :class="{ 'v-picker__title__btn--active': isActiveMinutesTab }",
          @click="showMinutesTabInTimeTab"
        ) {{ value | date('mm', true, '--') }}
        template(v-if="useSeconds")
          span :
          span.v-picker__title__btn(
            :class="{ 'v-picker__title__btn--active': isActiveSecondsTab }",
            @click="showSecondsTabInTimeTab"
          ) {{ value | date('ss', true, '--') }}
    .date-time-picker__body
      v-fade-transition
        v-date-picker(
          v-if="isActiveDateTab",
          :locale="$i18n.locale",
          :value="value | date('YYYY-MM-DD', true, null)",
          color="primary",
          no-title,
          @input="updateDate",
          @change="showHoursTabInTimeTab"
        )
      v-fade-transition
        v-time-picker(
          v-show="isActiveTimeTab",
          ref="timePicker",
          :value="value | date('timePickerWithSeconds', true, null)",
          :allowed-minutes="allowedMinutes",
          :useSeconds="useSeconds",
          color="primary",
          format="24hr",
          no-title,
          @input="updateTime",
          @change="showDateTab"
        )
    slot(name="footer")
</template>

<script>
import { VUETIFY_ANIMATION_DELAY } from '@/config';

const TABS = {
  date: 'date',
  time: 'time',
};

/**
 * Date time picker component
 *
 * @prop {Date} [value=null] - Date value
 * @prop {Boolean} [roundHours=false] - Deny to change minutes it will be only 0
 * @prop {Boolean} [opened=false] - Is fate time picker opened (need for v-menu)
 *
 * @event value#input
 */
export default {
  props: {
    value: {
      type: Date,
      default: null,
    },
    roundHours: {
      type: Boolean,
      default: false,
    },
    opened: {
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
      activeTab: TABS.date,
    };
  },
  computed: {
    isActiveDateTab() {
      return this.activeTab === TABS.date;
    },

    isActiveTimeTab() {
      return this.activeTab === TABS.time;
    },

    isActiveHoursTab() {
      return this.isActiveTimeTab && this.$refs.timePicker.selectingHour;
    },

    isActiveMinutesTab() {
      return this.isActiveTimeTab && this.$refs.timePicker.selectingMinute;
    },

    isActiveSecondsTab() {
      return this.isActiveTimeTab && this.$refs.timePicker.selectingSecond;
    },

    allowedMinutes() {
      if (this.roundHours) {
        return v => v === 0;
      }

      return null;
    },
  },
  watch: {
    opened(value) {
      if (!value) {
        setTimeout(() => {
          this.showDateTab();
        }, VUETIFY_ANIMATION_DELAY);
      }
    },
  },
  methods: {
    updateTime(time = '00:00:00') {
      const newValue = new Date(this.value ? this.value.getTime() : null);
      const [hours = 0, minutes = 0, seconds = 0] = time.split(':');

      newValue.setHours(parseInt(hours, 10) || 0, parseInt(minutes, 10) || 0, parseInt(seconds, 10) || 0, 0);

      this.$emit('input', newValue);
    },

    updateDate(date) {
      const newValue = new Date(this.value ? this.value.getTime() : null);
      const [year, month, day] = date.split('-');

      newValue.setFullYear(parseInt(year, 10), parseInt(month, 10) - 1, parseInt(day, 10));

      if (!this.value) {
        newValue.setHours(0, 0, 0, 0);
      } else if (this.useSeconds) {
        newValue.setMilliseconds(0);
      } else {
        newValue.setSeconds(0, 0);
      }

      this.$emit('input', newValue);
    },

    showDateTab() {
      this.activeTab = TABS.date;
    },

    showTimeTab() {
      this.activeTab = TABS.time;
    },

    showHoursTabInTimeTab() {
      /**
       * Change to vuetify hour tab
       */
      this.$refs.timePicker.selecting = 1;
      this.showTimeTab();
    },

    showMinutesTabInTimeTab() {
      /**
       * Change to vuetify minute tab
       */
      this.$refs.timePicker.selecting = 2;
      this.showTimeTab();
    },

    showSecondsTabInTimeTab() {
      /**
       * Change to vuetify minute tab
       */
      this.$refs.timePicker.selecting = 3;
      this.showTimeTab();
    },
  },
};
</script>

<style lang="scss">
  .date-time-picker {
    .date-time-picker__body {
      position: relative;
      width: 290px;
      height: 312px;

      .v-picker {
        position: absolute;
        top: 0;
        left: 0;
      }
    }

    .v-date-time-picker-title {
      line-height: 50px;
      font-size: 30px;
      font-weight: 500;

      &.use-seconds {
        font-size: 27px;
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

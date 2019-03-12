<template lang="pug">
div
  .v-picker__title.primary.text-xs-center
    span.v-date-time-picker-title
      span.v-picker__title__btn(
      @click="showDateTab",
      :class="{ active: isActiveDateTab }"
      ) {{ dateTimeObject | date('DD/MM/YYYY', true, '--/--/----') }}
      span &nbsp;
      span.v-picker__title__btn(
      @click="showHourTab",
      :class="{ active: isActiveHourTab }"
      ) {{ dateTimeObject | date('HH', true, '--') }}
      span :
      span.v-picker__title__btn(
      @click="showMinuteTab",
      :class="{ active: isActiveMinuteTab }"
      ) {{ dateTimeObject | date('mm', true, '--') }}
  div.date-time-picker__body
    v-fade-transition
      v-date-picker(
      v-if="isActiveDateTab",
      :locale="$i18n.locale",
      v-model="dateString",
      color="primary",
      no-title,
      @input="updateDateTimeObject",
      @change="showHourTab",
      )
    v-fade-transition
      v-time-picker(
      v-show="isActiveTimeTab",
      ref="timePicker"
      v-model="timeString",
      :allowed-minutes="allowedMinutes"
      color="primary",
      format="24hr",
      no-title,
      @input="updateDateTimeObject",
      @change="showDateTab"
      )
  .text-xs-center.dropdown-footer
    v-btn(@click.prevent="submit", color="primary", depressed) Ok
</template>

<script>
import moment from 'moment';

const TABS = {
  date: 'date',
  time: 'time',
};

export default {
  props: {
    value: Date,
    rules: [String, Object],
    roundHours: Boolean,
  },
  data() {
    const value = this.value ? moment(this.value) : null;

    return {
      activeTab: TABS.date,
      dateTimeObject: value,
      dateString: value ? value.format('YYYY-MM-DD') : null,
      timeString: value ? value.format('HH:mm') : null,
    };
  },
  computed: {
    isActiveDateTab() {
      return this.activeTab === TABS.date;
    },
    isActiveTimeTab() {
      return this.activeTab === TABS.time;
    },
    isActiveHourTab() {
      return this.isActiveTimeTab && this.$refs.timePicker.selectingHour;
    },
    isActiveMinuteTab() {
      return this.isActiveTimeTab && !this.$refs.timePicker.selectingHour;
    },
    allowedMinutes() {
      if (this.roundHours) {
        return v => v === 0;
      }
      return null;
    },
  },
  methods: {
    showDateTab() {
      this.activeTab = TABS.date;
    },

    showTimeTab() {
      this.activeTab = TABS.time;
    },

    showHourTab() {
      this.$refs.timePicker.selectingHour = true;
      this.showTimeTab();
    },

    showMinuteTab() {
      this.$refs.timePicker.selectingHour = false;
      this.showTimeTab();
    },

    updateDateTimeObject() {
      if (!this.timeString) {
        this.timeString = '00:00';
      }

      this.dateTimeObject = moment(`${this.dateString} ${this.timeString}`, 'YYYY-MM-DD HH:mm');
    },

    submit() {
      this.$emit('submit', this.dateTimeObject.toDate());
    },
  },
};
</script>

<style lang="scss">
  .date-time-picker {
    .date-time-picker__body {
      position: relative;
      width: 290px;
      height: 300px;

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

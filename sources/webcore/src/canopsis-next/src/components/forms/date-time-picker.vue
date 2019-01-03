<template lang="pug">
  v-menu(
  ref="menu",
  v-model="opened",
  content-class="date-time-picker",
  transition="slide-y-transition",
  max-width="290px",
  :close-on-content-click="false",
  right,
  lazy
  )
    div(slot="activator")
      v-text-field(
      readonly,
      :label="label",
      :error-messages="name ? errors.collect(name) : []",
      :value="dateTimeObject | date('DD/MM/YYYY HH:mm', true)",
      v-validate="rules",
      :data-vv-name="name",
      data-vv-validate-on="none",
      :append-icon="clearable ? 'close' : ''",
      @click:append="clear"
      )
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
        @input="updateDateTimeObject",
        @change="showHourTab"
        no-title,
        )
      v-fade-transition
        v-time-picker(
        v-show="isActiveTimeTab",
        ref="timePicker"
        v-model="timeString",
        @input="updateDateTimeObject",
        @change="showDateTab",
        format="24hr",
        :allowed-minutes="allowedMinutes"
        no-title,
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

/**
 * Date time picker component
 *
 * @prop {Boolean} [clearable] - if it is true then input field will be have cross button with clear event on click
 * @prop {Date} [value] - v-model
 * @prop {string} [label] - label of the input field
 * @prop {string} [name] - name property in the validation object
 * @prop {string} [rules] - validation rules in vee-validate format
 * @prop {string} [format='DD/MM/YYYY HH:mm'] - date format for display
 *
 * @event value#input
 * @type Date - new date value
 */
export default {
  inject: ['$validator'],
  props: {
    clearable: Boolean,
    value: Date,
    label: String,
    name: String,
    rules: [String, Object],
    roundHours: Boolean,
  },
  data() {
    const value = this.value ? moment(this.value) : null;

    return {
      opened: false,
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
  watch: {
    opened(value) {
      if (!value) {
        setTimeout(() => {
          this.activeTab = 'date';
        }, this.$config.VUETIFY_ANIMATION_DELAY);
      }
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

      this.$emit('input', this.dateTimeObject.toDate());

      this.validate();
    },

    clear() {
      this.dateTimeObject = null;
      this.dateString = '';
      this.timeString = '';

      this.$emit('input', this.dateTimeObject);

      this.validate();
    },

    submit() {
      this.validate();

      this.$refs.menu.save();
    },

    validate() {
      if (this.name && this.rules) {
        this.$nextTick(async () => {
          await this.$validator.validate(this.name);
        });
      }
    },
  },
};
</script>

<style lang="scss">
  .date-time-picker {
    .date-time-picker__body {
      position: relative;
      width: 290px;
      height: 290px;

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
      height: 246px;
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

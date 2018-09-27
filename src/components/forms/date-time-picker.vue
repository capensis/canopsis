<template lang="pug">
  v-menu(
  content-class="date-time-picker",
  transition="slide-y-transition",
  v-model="opened",
  ref="menu",
  :close-on-content-click="false",
  left,
  right,
  lazy
  max-width="290px"
  )
    div(slot="activator")
      v-text-field(
      readonly,
      :label="label",
      :error-messages="name ? errors.collect(name) : []",
      :value="dateTimeString",
      v-validate="rules",
      :data-vv-name="name",
      data-vv-validate-on="none",
      :append-icon="clearable ? 'close' : ''",
      @click:append="clear"
      )
    .v-picker__title.primary.text-xs-center
      span.subheading {{ dateTimeString || '--/--/---- --:--' }}
    v-tabs(v-model="activeTab", centered, grow)
      v-tab(href="#date")
        v-icon date_range
      v-tab(href="#time")
        v-icon access_time
      v-tab-item(id="date")
        v-date-picker(
        @input="updateDateTimeObject",
        :locale="$i18n.locale",
        v-model="dateString",
        no-title,
        )
      v-tab-item(id="time")
        v-time-picker(
        @input="updateDateTimeObject",
        v-model="timeString",
        format="24hr"
        no-title,
        )
    .text-xs-center.dropdown-footer
      v-btn(@click.prevent="submit", color="primary", depressed) Ok
</template>

<script>
import moment from 'moment';

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
    format: {
      type: String,
      default: 'DD/MM/YYYY HH:mm',
    },
  },
  data() {
    const value = this.value ? moment(this.value) : null;

    return {
      opened: false,
      activeTab: 'date',
      dateTimeObject: value,
      dateString: value ? value.format('YYYY-MM-DD') : '',
      timeString: value ? value.format('HH:mm') : '',
    };
  },
  computed: {
    dateTimeString() {
      return this.dateTimeObject ? this.dateTimeObject.format(this.format) : this.dateTimeObject;
    },
  },
  watch: {
    opened(value) {
      if (!value) {
        setTimeout(() => {
          this.activeTab = 'date';
        }, 300);
      }
    },
  },
  methods: {
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
  }
</style>

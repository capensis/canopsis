<template lang="pug">
  v-menu(
  left
  right
  transition="slide-y-transition"
  ref="menu",
  v-model="opened",
  :close-on-content-click="false",
  )
    v-text-field(
    readonly,
    slot="activator",
    :label="$t('modals.addChangeStateEvent.output')",
    :error-messages="errors.collect('date')",
    v-model="dateTimeString",
    v-validate="'required'",
    data-vv-name="date"
    )
    .picker__title.primary.text-xs-center
      span.subheading {{ dateTimeString }}
    v-tabs(v-model="activeTab", centered, hide-slider)
      v-tab(href="#date")
        v-icon date_range
      v-tab(href="#time")
        v-icon access_time
      v-tab-item(id="date")
        v-date-picker(
        @input="updateDateTimeObject",
        :locale="$i18n.locale",
        v-model="dateString",
        no-title
        )
      v-tab-item(id="time")
        v-time-picker(
        @input="updateDateTimeObject",
        v-model="timeString",
        format="24hr"
        no-title
        )
    .text-xs-center.dropdown-footer
      v-btn(depressed, color="primary", @click.prevent="submit") Ok
</template>

<script>
import moment from 'moment';

export default {
  props: {
    value: {
      type: Object,
      default() {
        return moment();
      },
    },
    format: {
      type: String,
      default: 'DD/MM/YYYY HH:mm',
    },
  },
  data() {
    return {
      opened: false,
      activeTab: 'date',
      dateTimeObject: this.value,
      dateString: this.value.format('YYYY-MM-DD'),
      timeString: this.value.format('HH:mm'),
    };
  },
  computed: {
    dateTimeString() {
      return this.dateTimeObject.format(this.format);
    },
  },
  methods: {
    updateDateTimeObject() {
      this.dateTimeObject = moment(`${this.dateString} ${this.timeString}`);

      this.$emit('input', this.dateTimeObject);
    },
    submit() {
      this.$refs.menu.save('asd');
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
};
</script>

<template lang="pug">
  div
    v-layout(row)
      v-switch(v-model="showRRule", :label="$t('modals.createPbehavior.fields.rRuleQuestion')")
    template(v-if="showRRule")
      v-layout(row)
        v-tabs.r-rule-tabs(v-model="activeTab", centered)
          v-tab(href="#simple") {{ $t('rRule.tabs.simple') }}
          v-tab(href="#advanced") {{ $t('rRule.tabs.advanced') }}
          v-tab-item(id="simple")
            div
              div
                v-select(
                :label="$t('rRule.fields.freq')",
                v-model="form.rRuleOptions.freq",
                :items="selectItems.frequencies",
                @input="changeRRuleOption"
                )
              div
                v-select(
                :label="$t('rRule.fields.byweekday')",
                v-model="form.rRuleOptions.byweekday",
                :items="selectItems.weekDays",
                @input="changeRRuleOption",
                multiple,
                chips
                )
              div
                v-text-field(
                type="number",
                :label="$t('rRule.fields.count')",
                v-model="form.rRuleOptions.count",
                @input="changeRRuleOption"
                )
              div
                v-text-field(
                type="number",
                :label="$t('rRule.fields.interval')",
                v-model="form.rRuleOptions.interval",
                @input="changeRRuleOption"
                )
          v-tab-item(id="advanced")
            div
              div
                v-select(
                :label="$t('rRule.fields.wkst')",
                v-model="form.rRuleOptions.wkst",
                :items="selectItems.weekDays",
                @input="changeRRuleOption"
                )
              div
                v-select(
                :label="$t('rRule.fields.bymonth')",
                v-model="form.rRuleOptions.bymonth",
                :items="selectItems.months",
                @input="changeRRuleOption",
                multiple,
                chips
                )
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.bysetpos.label')",
                    v-model="form.advancedRRuleOptions.bysetpos",
                    :error-messages="errors.collect('bysetpos')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="bysetpos",
                    persistent-hint,
                    )
                  span {{ $t('rRule.fields.bysetpos.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.bymonthday.label')",
                    v-model="form.advancedRRuleOptions.bymonthday",
                    :error-messages="errors.collect('bymonthday')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="bymonthday",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.bymonthday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.byyearday.label')",
                    v-model="form.advancedRRuleOptions.byyearday",
                    :error-messages="errors.collect('byyearday')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="byyearday",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.byyearday.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.byweekno.label')",
                    v-model="form.advancedRRuleOptions.byweekno",
                    :error-messages="errors.collect('byweekno')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="byweekno",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.byweekno.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.byhour.label')",
                    v-model="form.advancedRRuleOptions.byhour",
                    :error-messages="errors.collect('byhour')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="byhour",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.byhour.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.byminute.label')",
                    v-model="form.advancedRRuleOptions.byminute",
                    :error-messages="errors.collect('byminute')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="byminute",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.byminute.tooltip') }}
              div
                v-tooltip(left, max-width="250")
                  div(slot="activator")
                    v-text-field(
                    :label="$t('rRule.fields.bysecond.label')",
                    v-model="form.advancedRRuleOptions.bysecond",
                    :error-messages="errors.collect('bysecond')",
                    :hint="$t('rRule.advancedHint')",
                    @input="changeRRuleAdvancedOption",
                    v-validate="{ regex: advancedFieldRegex }",
                    data-vv-name="bysecond",
                    persistent-hint
                    )
                  span {{ $t('rRule.fields.bysecond.tooltip') }}
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.stringLabel') }}
        v-flex(xs10)
          p {{ rRuleString }}
      v-layout(row)
        v-flex(xs2)
          strong {{ $t('rRule.textLabel') }}
        v-flex(xs10)
          p {{ rRuleText }}
      v-layout(row)
        v-alert(:value="errors.has('rRule')", type="error")
          span {{ errors.first('rRule') }}
</template>

<script>
import RRule from 'rrule';
import { mapValues, pickBy } from 'lodash';

import DateTimePicker from '@/components/forms/fields/date-time-picker.vue';

/**
 * RRule form component
 *
 * @event rrule#input
 * @type {Object|null} - RRule object
 */
export default {
  inject: ['$validator'],
  components: { DateTimePicker },
  data() {
    const rRuleOptions = {
      freq: RRule.DAILY,
    };

    return {
      activeTab: 'simple',
      showRRule: false,
      rRuleObject: new RRule(rRuleOptions),
      advancedFieldRegex:
        /^(-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))(,-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))*,?$/,
      form: {
        rRuleOptions,
        advancedRRuleOptions: {},
      },
      selectItems: {
        frequencies: [
          { text: 'Secondly', value: RRule.SECONDLY },
          { text: 'Minutely', value: RRule.MINUTELY },
          { text: 'Hourly', value: RRule.HOURLY },
          { text: 'Daily', value: RRule.DAILY },
          { text: 'Weekly', value: RRule.WEEKLY },
          { text: 'Monthly', value: RRule.MONTHLY },
          { text: 'Yearly', value: RRule.YEARLY },
        ],
        weekDays: [
          { text: 'Monday', value: RRule.MO.weekday },
          { text: 'Tuesday', value: RRule.TU.weekday },
          { text: 'Wednesday', value: RRule.WE.weekday },
          { text: 'Thursday', value: RRule.TH.weekday },
          { text: 'Friday', value: RRule.FR.weekday },
          { text: 'Saturday', value: RRule.SA.weekday },
          { text: 'Sunday', value: RRule.SU.weekday },
        ],
        months: [
          { text: 'January', value: 1 },
          { text: 'February', value: 2 },
          { text: 'March', value: 3 },
          { text: 'April', value: 4 },
          { text: 'May', value: 5 },
          { text: 'June', value: 6 },
          { text: 'July', value: 7 },
          { text: 'August', value: 8 },
          { text: 'September', value: 9 },
          { text: 'October', value: 10 },
          { text: 'November', value: 11 },
          { text: 'December', value: 12 },
        ],
      },
    };
  },
  computed: {
    rRuleString() {
      return this.rRuleObject.toString();
    },

    rRuleText() {
      return this.rRuleObject.toText();
    },
  },
  watch: {
    showRRule(value) {
      if (!value) {
        this.errors.remove('rRule');
        this.$emit('input', null);
      } else {
        this.changeRRuleOption();
      }
    },
  },
  methods: {
    changeRRuleAdvancedOption() {
      this.form.rRule = {
        ...this.form.rRuleOptions,
        ...mapValues(this.form.advancedRRuleOptions, options => options.split(',').filter(v => v)),
      };

      this.changeRRuleOption();
    },
    /**
     * For each changes in the form we call this function.
     * If RRule isn't valid then add error message to invisible RRule field
     * Else remove errors and $emit changes
     */
    changeRRuleOption() {
      try {
        this.form.rRule = {
          ...this.form.rRule,
          ...this.form.rRuleOptions,
        };

        const correctOptions = pickBy(this.form.rRule, v => v !== '');

        this.rRuleObject = new RRule(correctOptions);

        if (!this.errors.has('rRule') && !this.rRuleObject.isFullyConvertibleToText()) {
          this.errors.add('rRule', this.$t('rRule.errors.main'));
        } else {
          this.errors.remove('rRule');

          if (this.showRRule) {
            this.$emit('input', this.rRuleString);
          }
        }
      } catch (err) {
        console.warn(err);
      }
    },
  },
};
</script>

<style scoped>
  .r-rule-tabs {
    width: 100%;
  }
  p {
    -ms-word-break: break-all;
    word-break: break-all;
  }
</style>

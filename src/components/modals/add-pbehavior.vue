<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700", lazy)
    v-form(@submit.prevent="submit")
      v-card
        v-card-title
          span.headline {{ $t('modals.addPbehavior.title') }}
        v-card-text
          v-layout(row)
            v-text-field(
            :label="$t('modals.addPbehavior.fields.name')",
            :error-messages="errors.collect('name')",
            v-model="form.name",
            v-validate="'required'",
            data-vv-name="name"
            )
          v-layout(row)
            date-time-picker(
            :label="$t('modals.addPbehavior.fields.start')",
            v-model="form.tstart",
            name="tstart",
            rules="required",
            @input="changeTstart"
            )
          v-layout(row)
            date-time-picker(
            :label="$t('modals.addPbehavior.fields.end')",
            v-model="form.tend",
            name="tend",
            rules="required"
            )
          v-layout(row)
            v-switch(v-model="showRRule", :label="$t('modals.addPbehavior.title')")
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
                      v-model="form.rRule.freq",
                      :items="selectItems.frequencies",
                      @input="changeRRuleOption"
                      )
                    div
                      date-time-picker(
                      :label="$t('rRule.fields.until')",
                      v-model="form.rRule.until",
                      clearable,
                      @input="changeRRuleOption"
                      )
                    div
                      v-select(
                      :label="$t('rRule.fields.byweekday')",
                      v-model="form.rRule.byweekday",
                      :items="selectItems.weekDays",
                      @input="changeRRuleOption"
                      multiple,
                      chips
                      )
                    div
                      v-text-field(
                      :label="$t('rRule.fields.count')",
                      type="number",
                      v-model="form.rRule.count",
                      @input="changeRRuleOption"
                      )
                    div
                      v-text-field(
                      :label="$t('rRule.fields.interval')",
                      type="number",
                      v-model="form.rRule.interval",
                      @input="changeRRuleOption"
                      )
                v-tab-item(id="advanced")
                  div
                    div
                      v-select(
                      :label="$t('rRule.fields.wkst')",
                      v-model="form.rRule.wkst",
                      :items="selectItems.weekDays",
                      @input="changeRRuleOption"
                      )
                    div
                      v-select(
                      :label="$t('rRule.fields.bymonth')",
                      v-model="form.rRule.bymonth",
                      :items="selectItems.months",
                      @input="changeRRuleOption"
                      multiple,
                      chips
                      )
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.bysetpos",
                          :label="$t('rRule.fields.bysetpos.label')",
                          :error-messages="errors.collect('bysetpos')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="bysetpos",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint,
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.bysetpos.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.bymonthday",
                          :label="$t('rRule.fields.bymonthday.label')",
                          :error-messages="errors.collect('bymonthday')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="bymonthday",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.bymonthday.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.byyearday",
                          :label="$t('rRule.fields.byyearday.label')",
                          :error-messages="errors.collect('byyearday')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="byyearday",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.byyearday.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.byweekno",
                          :label="$t('rRule.fields.byweekno.label')",
                          :error-messages="errors.collect('byweekno')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="byweekno",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.byweekno.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.byhour",
                          :label="$t('rRule.fields.byhour.label')",
                          :error-messages="errors.collect('byhour')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="byhour",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.byhour.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.byminute",
                          :label="$t('rRule.fields.byminute.label')",
                          :error-messages="errors.collect('byminute')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="byminute",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.byminute.tooltip') }}
                    div
                      v-tooltip(left, max-width="250")
                        div(slot="activator")
                          v-text-field(
                          v-model="form.advancedRRule.bysecond",
                          :label="$t('rRule.fields.bysecond.label')",
                          :error-messages="errors.collect('bysecond')",
                          v-validate="{ regex: advancedFieldRegexp }",
                          data-vv-name="bysecond",
                          :hint="$t('rRule.advancedHint')",
                          persistent-hint
                          @input="changeRRuleAdvancedOption"
                          )
                        span {{ $t('rRule.fields.bysecond.tooltip') }}
            v-layout(row)
              v-flex(xs2)
                strong {{ $t('rRule.stringLabel') }}
              v-flex(xs10)
                p {{ rRuleObject.toString() }}
            v-layout(row)
              v-flex(xs2)
                strong {{ $t('rRule.textLabel') }}
              v-flex(xs10)
                p {{ rRuleObject.toText() }}
          v-layout(row)
            v-select(
            label="Reason",
            v-model="form.reason",
            :items="selectItems.reasons"
            )
          v-layout(row)
            v-select(
            label="Type",
            v-model="form.type_",
            :items="selectItems.types"
            )
        v-card-actions
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import RRule from 'rrule';
import mapValues from 'lodash/mapValues';

import DateTimePicker from '@/components/forms/date-time-picker.vue';

import ModalMixin from './modal-mixin';

export default {
  name: 'add-pbehavior',
  mixins: [ModalMixin],
  components: { DateTimePicker },
  data() {
    const now = new Date();
    const reasons = ['Problème Habilitation', 'Problème Robot', 'Problème Scénario', 'Autre'];
    const types = ['Pause', 'Maintenance', 'Hors plage horaire de surveillance'];
    const frequencies = [
      { text: 'Secondly', value: RRule.SECONDLY },
      { text: 'Minutely', value: RRule.MINUTELY },
      { text: 'Hourly', value: RRule.HOURLY },
      { text: 'Daily', value: RRule.DAILY },
      { text: 'Weekly', value: RRule.WEEKLY },
      { text: 'Monthly', value: RRule.MONTHLY },
      { text: 'Yearly', value: RRule.YEARLY },
    ];
    const rRuleOptions = {
      dtstart: now,
      freq: frequencies[3].value, // Daily
    };

    return {
      activeTab: 'simple',
      showRRule: false,
      rRuleObject: new RRule(rRuleOptions),
      advancedFieldRegexp:
        /^(-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))(,-?(?:[1-9]|[1-9]\d|[12]\d\d|3[0-6][0-6]?))*,?$/,
      form: {
        name: '',
        tstart: now,
        tend: now,
        type_: types[0],
        reason: reasons[0],
        rRule: rRuleOptions,
        advancedRRule: {},
      },
      selectItems: {
        reasons,
        types,
        frequencies,
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
  methods: {
    changeTstart(value) {
      this.form.rRule.dtstart = value;

      this.changeRRuleOption();
    },
    changeRRuleAdvancedOption() {
      this.form.rRule = {
        ...this.form.rRule,
        ...mapValues(this.form.advancedRRule, options => options.split(',').filter(v => v)),
      };

      this.changeRRuleOption();
    },
    changeRRuleOption() {
      try {
        this.rRuleObject = new RRule(this.form.rRule);
      } catch (err) {
        console.warn(err);
      }
    },
    async submit() {
      const isValid = await this.$validator.validateAll();

      console.log(isValid);
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

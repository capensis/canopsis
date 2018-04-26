<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700", lazy)
    v-form(@submit.prevent="submit")
      v-card
        v-card-title
          span.headline {{ $t('modals.addChangeStateEvent.title') }}
        v-card-text
          v-layout(row)
            v-text-field(
            label="Name",
            :error-messages="errors.collect('output')",
            v-model="form.name",
            v-validate="'required'",
            data-vv-name="output"
            )
          v-layout(row)
            date-time-picker(
            v-model="form.tstart",
            label="Start",
            name="tstart",
            rules="required"
            )
          v-layout(row)
            date-time-picker(
            v-model="form.tend",
            label="End",
            name="tend",
            rules="required"
            )
          v-layout(row)
            v-tabs.r-rule-tabs(v-model="activeTab", centered)
              v-tab(href="#simple") Simple
              v-tab(href="#advanced") Advanced
              v-tab(href="#text-input") Text input
              v-tab-item(id="simple")
                div
                  div
                    v-select(
                    label="Frequency",
                    v-model="form.rRule.freq",
                    :items="selectItems.frequencies"
                    )
                  div
                    date-time-picker(label="Until", v-model="form.rRule.until", clearable)
                  div
                    v-select(
                    label="By week day",
                    v-model="form.rRule.byweekday",
                    :items="selectItems.weekDays",
                    multiple,
                    chips
                    )
                  div
                    v-text-field(type="number", label="Repeat", v-model="form.rRule.count")
                  div
                    v-text-field(
                    type="number",
                    label="Interval",
                    v-model="form.rRule.interval",
                    @input="changeRRuleProp"
                    )
              v-tab-item(id="advanced")
                div
                  div
                    v-select(
                    label="Week start",
                    v-model="form.rRule.wkst",
                    :items="selectItems.weekDays",
                    )
                  div
                    v-select(
                    label="Select",
                    v-model="form.rRule.bymonth",
                    :items="selectItems.months",
                    multiple,
                    chips
                    )
                  div
                    v-tooltip(right)
                      div(slot="activator")
                        v-text-field(type="number", label="By set position", v-model="form.rRule.bysetpos")
                      span Something
                  div
                    v-text-field(type="number", label="By month day", v-model="form.rRule.bymonthday")
                  div
                    v-text-field(type="number", label="By year day", v-model="form.rRule.byyearday")
                  div
                    v-text-field(type="number", label="By week n°", v-model="form.rRule.byweekno")
                  div
                    v-text-field(type="number", label="By hour", v-model="form.rRule.byhour")
                  div
                    v-text-field(type="number", label="By minute", v-model="form.rRule.byminute")
                  div
                    v-text-field(type="number", label="By second", v-model="form.rRule.bysecond")
              v-tab-item(id="text-input")
                div
                  v-text-field(
                  label="Rrule",
                  :value="form.rRuleObject.toString()",
                  @input="changeRRuleString"
                  )
          v-layout(row)
            v-flex(xs2)
              strong Rrule
            v-flex(xs10) {{ form.rRuleObject.toString() }}
          v-layout(row)
            v-flex(xs2)
              strong Summary
            v-flex(xs10) {{ form.rRuleObject.toText() }}
          v-layout(row)
            v-select(label="Reason")
          v-layout(row)
            v-select(label="Type")
        v-card-actions
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import RRule from 'rrule';

import DateTimePicker from '@/components/forms/date-time-picker.vue';

import ModalMixin from './modal-mixin';

export default {
  name: 'add-pbehavior',
  mixins: [ModalMixin],
  components: { DateTimePicker },
  data() {
    const rRule = new RRule({
      until: null,
      freq: RRule.DAILY,
      wkst: RRule.MO.weekday,
      byweekday: [RRule.MO.weekday, RRule.TU.weekday, RRule.WE.weekday],
      bymonth: [1],
      interval: 1,
    });

    return {
      activeTab: 'simple',
      form: {
        name: '',
        tstart: new Date(),
        tend: new Date(),
        rRule: { ...rRule.options },
        rRuleObject: rRule,
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
  methods: {
    changeRRuleProp() {
      try {
        this.form.rRuleObject = new RRule(this.form.rRule);
      } catch (err) {
        console.warn(err);
      }
    },
    changeRRuleString(value) {
      try {
        this.form.rRuleObject = RRule.fromString(value);
        this.form.rRule = { ...this.form.rRuleObject.options };
      } catch (err) {
        // TODO: add error
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
</style>

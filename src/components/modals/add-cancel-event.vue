<template lang="pug">
  v-form(@submit.prevent="submit")
    v-card
      v-card-title
        span.headline {{ $t('modals.addCancelEvent.title') }}
      v-card-text
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
            :label="$t('modals.addCancelEvent.output')",
            :error-messages="errors.collect('output')",
            v-model="form.output",
            v-validate="'required'",
            data-vv-name="output"
            )
      v-card-actions
        v-btn(type="submit", :disabled="errors.any()", color="primary") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import Vue from 'vue';
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';
import AlarmMixin from '../mixins/alarm';


const component = Vue.extend({
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  name: 'add-cancel-event',
  mixins: [
    AlarmMixin,
  ],

  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.cancelAlarmConfirmation({
            comment: this.form.output,
            alarmData: {
              connector: 'toto',
              connector_name: 'toto',
              component: 'localhost',
              state: 0,
              state_type: 1,
              resource: 'res99',
            },
          });
          this.form.output = '';
        } catch (e) {
          console.log(e);
        }
      }
    },
  },
});

export default component;
</script>

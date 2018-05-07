<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700")
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
          v-btn(type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';
import ModalMixin from './modal-mixin';

export default {
  name: 'add-cancel-event',
  mixins: [
    ModalMixin,
  ],
  components: {
    AlarmGeneralTable,
  },
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        console.log(this.$store);
        const response = await this.$store.dispatch('alarm/cancelConfirmation', {
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
        console.log(response);
      }
    },
  },
};
</script>

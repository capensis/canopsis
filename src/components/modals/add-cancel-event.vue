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
import { createNamespacedHelpers } from 'vuex';
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';

const { mapActions } = createNamespacedHelpers('event');

export default {
  name: 'add-cancel-event',
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
  $_veeValidate: {
    validator: 'new',
  },

  methods: {
    ...mapActions([
      'cancelAck',
    ]),
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.cancelAck({
          resource: 'res0',
          id: '652d34d0-4eda-11e8-841e-0242ac12000a',
          customAttributes: {
            output: this.form.output,
            cancel: 1,
          },
        });
      // todo hide modal action
      }
    },
  },
};
</script>

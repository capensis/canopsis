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

const { mapActions } = createNamespacedHelpers('alarmEvents');

export default {
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

  methods: {
    ...mapActions([
      'cancel',
    ]),
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.cancel({
            comment: this.form.output,
            resource: 'res99',
            id: 'ac4f92ea-4eda-11e8-841e-0242ac12000a',
          });
          this.form.output = '';
        } catch (e) {
          console.log(e);
        }
      }
    },
  },
};
</script>

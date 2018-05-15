<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.addAckEvent.title') }}
    v-card-text
      v-container
        v-layout(row align-center)
          v-flex.text-xs-center
            alarm-general-table
        v-layout(row)
          v-divider.my-3
        v-layout(row)
          v-text-field(
          :label="$t('modals.addAckEvent.ticket')",
          :error-messages="errors.collect('ticket')",
          v-model="form.ticket",
          v-validate="rules",
          data-vv-name="ticket"
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.addAckEvent.output')",
          :error-messages="errors.collect('output')",
          v-model="form.output",
          v-validate="rules",
          data-vv-name="output",
          multi-line
          )
        v-layout(row)
          v-checkbox(:label="$t('modals.addAckEvent.ackResources')", v-model="form.ack_resources")
    v-card-actions
      v-btn(@click.prevent="submit", color="primary") {{ $t('common.actions.acknowledge') }}
      v-btn(
      @click.prevent="submitWithAdditions",
      color="warning"
      ) {{ $t('common.actions.acknowledgeAndReport') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';

const { mapActions } = createNamespacedHelpers('event');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    AlarmGeneralTable,
  },
  data() {
    return {
      showValidationErrors: false,
      form: {
        ticket: '',
        output: '',
        ack_resources: false,
      },
    };
  },
  computed: {
    rules() {
      return this.showValidationErrors ? 'required' : '';
    },
  },
  methods: {
    ...mapActions([
      'ack',
      'declare',
    ]),
    async submit() {
      this.showValidationErrors = false;
      this.errors.clear();
      const requestData = {
        id: '652d34d0-4eda-11e8-841e-0242ac12000a',
        resource: 'res0',
        customAttributes: {},
      };

      if (this.form.ticket) {
        requestData.customAttributes.ticket = this.form.ticket;
      }

      if (this.form.output) {
        requestData.customAttributes.output = this.form.output;
      }

      this.ack(requestData);
    },
    async submitWithAdditions() {
      this.showValidationErrors = true;

      this.$nextTick(async () => {
        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          await this.submit();

          await this.declare({
            id: '652d34d0-4eda-11e8-841e-0242ac12000a',
            resource: 'res0',
            customAttributes: {
              output: 'declare ticket',
            },
          });
        //  todo hide modal action
        }
      });
    },
  },
};
</script>

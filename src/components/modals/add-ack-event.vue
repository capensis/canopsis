<template lang="pug">
  v-dialog(:value="opened", @input="hideModal", max-width="700")
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
        v-btn(@click.prevent="hideModal") {{ $t('common.actions.close') }}
        v-btn(@click.prevent="submit") {{ $t('common.actions.acknowledge') }}
        v-btn(@click.prevent="submitWithAdditions") {{ $t('common.actions.acknowledgeAndReport') }}
</template>

<script>
import AlarmGeneralTable from '@/components/tables/alarm-general.vue';

import ModalHOC from './modal-hoc';

export default ModalHOC({
  name: 'add-ack-event',
  components: {
    AlarmGeneralTable,
  },
  data() {
    return {
      form: {
        ticket: '',
        output: '',
        ack_resources: false,
      },
      showValidationErrors: false,
      table: {
        headers: [
          {
            text: 'author',
            sortable: false,
          },
          {
            text: 'connector',
            sortable: false,
          },
          {
            text: 'component',
            sortable: false,
          },
          {
            text: 'resource',
            sortable: false,
          },
        ],
        items: [
          {
            name: 'something',
          },
        ],
      },
    };
  },
  computed: {
    rules() {
      return this.showValidationErrors ? 'required' : '';
    },
  },
  methods: {
    async submit() {
      this.showValidationErrors = false;
      this.errors.clear();
      console.log('ADD_ACK');
    },
    async submitWithAdditions() {
      this.showValidationErrors = true;

      this.$nextTick(async () => {
        const formIsValid = await this.$validator.validateAll();

        if (formIsValid) {
          await this.submit();

          // send request to add report
        }
      });
    },
  },
});
</script>

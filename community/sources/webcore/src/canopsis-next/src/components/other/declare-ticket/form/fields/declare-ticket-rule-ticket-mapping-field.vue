<template lang="pug">
  v-layout(column)
    c-information-block(
      :title="$t('declareTicket.ticketUrlAndId')",
      :help-text="$t('declareTicket.ticketUrlAndIdHelpText')",
      help-icon="help",
      help-icon-color="grey darken-1"
    )
      v-layout
        v-flex(xs6)
          c-enabled-field(v-model="form.empty_response", :label="$t('declareTicket.emptyResponse')")
        v-flex(xs6)
          c-enabled-field(
            v-model="form.is_regexp",
            :disabled="form.empty_response",
            :label="$t('declareTicket.isRegexp')"
          )
      v-layout.mr-5
        v-flex.mr-3(xs6)
          v-text-field(
            v-field="form.ticket_id",
            v-validate="'required'",
            :label="$t('declareTicket.ticketID')",
            :error-messages="errors.collect('ticket_id')",
            name="ticket_id"
          )
            template(#append="")
              c-help-icon(
                :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketID') })",
                icon="help",
                color="grey darken-1",
                left
              )
        v-flex(xs6)
          v-text-field(
            v-field="form.ticket_url",
            :label="$t('declareTicket.ticketURL')"
          )
            template(#append="")
              c-help-icon(
                :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketURL') })",
                icon="help",
                color="grey darken-1",
                left
              )
      c-information-block(:title="$t('declareTicket.customFields')")
        v-flex(v-if="!form.mapping.length", xs12)
          v-alert(:value="true", type="info") {{ $t('declareTicket.emptyFields') }}
        c-text-pairs-field(
          v-field="form.mapping",
          :text-label="$t('declareTicket.alarmFieldName')",
          :value-label="$t('declareTicket.responseField')",
          :name="name",
          text-required,
          value-required
        )
          template(#append-response="{ item }")
            c-help-icon(
              v-if="item.text",
              :text="$t('declareTicket.responseFieldHelpText', { field: item.text })",
              icon="help",
              color="grey darken-1",
              left
            )
</template>

<script>
import { formMixin } from '@/mixins/form';

import RequestForm from '@/components/forms/request/request-form.vue';

export default {
  inject: ['$validator'],
  components: { RequestForm },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'declare_ticket',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

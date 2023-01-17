<template lang="pug">
  v-layout(column)
    c-information-block(
      :title="$t('declareTicket.ticketUrlAndId')",
      :help-text="$t('declareTicket.ticketUrlAndIdHelpText')",
      help-icon="help",
      help-icon-color="grey darken-1"
    )
      c-alert(v-if="isDeclareTicketExist", type="info") {{ $t('declareTicket.webhookTicketDeclarationExist') }}
      v-layout
        v-flex(xs6)
          c-enabled-field(v-field="value.enabled", :disabled="isDeclareTicketExist")
        v-flex(v-if="value.enabled", xs6)
          c-enabled-field(v-field="value.is_regexp", :label="$t('declareTicket.isRegexp')")
      template(v-if="value.enabled")
        v-layout.mr-5
          v-flex.mr-3(xs6)
            declare-ticket-rule-ticket-id-field(
              v-field="value.ticket_id",
              :disabled="disabled",
              :name="ticketIdFieldName"
            )
          v-flex(xs6)
            v-text-field(v-field="value.ticket_url", :label="$t('declareTicket.ticketURL')")
              template(#append="")
                c-help-icon(
                  :text="$t('declareTicket.responseFieldHelpText', { field: $t('declareTicket.ticketURL') })",
                  icon="help",
                  color="grey darken-1",
                  left
                )
        c-information-block(:title="$t('declareTicket.customFields')")
          c-alert(v-if="!value.mapping.length", type="info") {{ $t('declareTicket.emptyFields') }}
          c-text-pairs-field(
            v-field="value.mapping",
            :text-label="$t('declareTicket.alarmFieldName')",
            :value-label="$t('declareTicket.responseField')",
            :name="name",
            text-required,
            value-required
          )
            template(#append-value="{ item }")
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

import DeclareTicketRuleTicketIdField from './declare-ticket-rule-ticket-id-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleTicketIdField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
    isTicketIdExist: {
      type: Boolean,
      default: false,
    },
    isDeclareTicketExist: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ticketIdFieldName() {
      return `${this.name}.ticket_id`;
    },
  },
};
</script>

<template>
  <v-layout column>
    <request-form
      v-field="webhook.request"
      :name="`${name}.request`"
      :headers-variables="payloadVariables"
      :payload-variables="payloadVariables"
      :url-variables="payloadVariables"
    />
    <declare-ticket-rule-ticket-mapping-field v-field="webhook.declare_ticket" />
    <v-layout justify-end>
      <v-btn
        :loading="checking"
        color="orange"
        dark
        @click="validateTemplateVariables"
      >
        {{ $t('declareTicket.checkSyntax') }}
      </v-btn>
    </v-layout>
  </v-layout>
</template>

<script>
import flatten from 'flat';

import { requestTemplateVariablesErrorsToForm } from '@/helpers/entities/shared/request/form';

import { formMixin, validationErrorsMixinCreator } from '@/mixins/form';
import { entitiesTemplateValidatorMixin } from '@/mixins/entities/template-validator';
import { payloadVariablesMixin } from '@/mixins/payload/variables';

import RequestForm from '@/components/forms/request/request-form.vue';
import DeclareTicketRuleTicketMappingField from '@/components/other/declare-ticket/form/fields/declare-ticket-rule-ticket-mapping-field.vue';

export default {
  inject: ['$validator'],
  components: { DeclareTicketRuleTicketMappingField, RequestForm },
  mixins: [
    formMixin,
    payloadVariablesMixin,
    entitiesTemplateValidatorMixin,
    validationErrorsMixinCreator(),
  ],
  model: {
    prop: 'webhook',
    event: 'input',
  },
  props: {
    webhook: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      checking: false,
    };
  },
  computed: {
    payloadVariables() {
      const variables = [
        ...this.alarmPayloadVariables,
      ];

      if (this.hasPrevious) {
        variables.push(...this.payloadVariablesFromPreviousStep);
      }

      variables.push(...this.additionalDataVariables);

      return variables;
    },
  },
  methods: {
    async validateRequestTemplates(request) {
      const [url, payload, ...headers] = await this.validateScenariosVariables({
        data: [
          { text: request.url },
          { text: request.payload },
          ...request.headers.map(({ value }) => ({ text: value })),
        ],
      });

      return {
        url,
        payload,
        headers,
      };
    },

    async validateFormTemplates(webhook) {
      return {
        request: await this.validateRequestTemplates(webhook.request),
      };
    },

    scenarioRequestErrorsToForm({ request }) {
      const flattenErrors = flatten({
        request: requestTemplateVariablesErrorsToForm(request, this.webhook.request),
      });

      return Object.entries(flattenErrors).reduce((acc, [key, value]) => {
        acc[`${this.name}.${key}`] = value;

        return acc;
      }, {});
    },

    async validateTemplateVariables() {
      this.checking = true;

      try {
        const errors = await this.validateFormTemplates(this.webhook);

        const wasSet = this.setFormErrors(this.scenarioRequestErrorsToForm(errors));

        if (!wasSet) {
          this.$popups.success({ text: this.$t('declareTicket.syntaxIsValid') });
        }
      } catch (err) {
        console.error(err);
      } finally {
        this.checking = false;
      }
    },
  },
};
</script>

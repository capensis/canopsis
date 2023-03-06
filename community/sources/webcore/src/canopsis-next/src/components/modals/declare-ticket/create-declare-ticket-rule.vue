<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        declare-ticket-rule-form(v-model="form")
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.orange.white--text(
          :loading="checking",
          @click="validateTemplateVariables"
        ) {{ $t('declareTicket.checkSyntax') }}
        v-btn.primary(
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  declareTicketRuleToForm,
  formToDeclareTicketRule,
  declareTicketRuleErrorsToForm,
  declareTicketRuleTemplateVariablesErrorsToForm,
} from '@/helpers/forms/declare-ticket-rule';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesTemplateValidatorMixin } from '@/mixins/entities/template-validator';

import DeclareTicketRuleForm from '@/components/other/declare-ticket/form/declare-ticket-rule-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createDeclareTicketRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { DeclareTicketRuleForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesTemplateValidatorMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: declareTicketRuleToForm(this.modal.config.declareTicketRule),
      checking: false,
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createDeclareTicketRule.create.title');
    },
  },
  methods: {
    validateRequestHeadersTemplates(headers) {
      return Promise.all(
        headers.map(({ value }) => this.validateDeclareTicketRulesVariables({ data: { text: value } })),
      );
    },

    async validateRequestTemplates(request) {
      const [url, payload, headers] = await Promise.all([
        this.validateDeclareTicketRulesVariables({ data: { text: request.url } }),
        this.validateDeclareTicketRulesVariables({ data: { text: request.payload } }),
        this.validateRequestHeadersTemplates(request.headers),
      ]);

      return {
        url,
        payload,
        headers,
      };
    },

    validateWebhooksTemplates(webhooks) {
      return Promise.all(webhooks.map(async ({ request }) => ({
        request: await this.validateRequestTemplates(request),
      })));
    },

    async validateFormTemplates(form) {
      return {
        webhooks: await this.validateWebhooksTemplates(form.webhooks),
      };
    },

    async validateTemplateVariables() {
      this.checking = true;

      try {
        const errors = await this.validateFormTemplates(this.form);

        const wasSet = this.setFormErrors(declareTicketRuleTemplateVariablesErrorsToForm(errors, this.form));

        if (!wasSet) {
          this.$popups.success({ text: this.$t('declareTicket.syntaxIsValid') });
        }
      } catch (err) {
        console.error(err);
      } finally {
        this.checking = false;
      }
    },

    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToDeclareTicketRule(this.form));
          }

          this.$modals.hide();
        } catch (err) {
          if (err.error) {
            this.$popups.error({ text: err.error });
          } else {
            this.setFormErrors(declareTicketRuleErrorsToForm(err, this.form));
          }
        }
      }
    },
  },
};
</script>

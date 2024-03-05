<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <declare-ticket-rule-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="checking"
          color="orange"
          dark
          @click="validateTemplateVariables"
        >
          {{ $t('declareTicket.checkSyntax') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import {
  declareTicketRuleToForm,
  formToDeclareTicketRule,
  declareTicketRuleErrorsToForm,
  declareTicketRuleTemplateVariablesErrorsToForm,
} from '@/helpers/entities/declare-ticket/rule/form';

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
    async validateRequestTemplates(request) {
      const [url, payload, ...headers] = await this.validateDeclareTicketRulesVariables({
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

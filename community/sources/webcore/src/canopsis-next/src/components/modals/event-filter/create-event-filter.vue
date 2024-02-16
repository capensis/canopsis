<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <event-filter-form
          v-model="form"
          :is-disabled-id-field="config.isDisabledIdField"
        />
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
          v-if="hasVariablesFields"
          :loading="checking"
          color="orange"
          dark
          @click="validateTemplateVariables"
        >
          {{ $t('declareTicket.checkSyntax') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
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
  eventFilterRuleTemplateVariablesErrorsToForm,
  eventFilterToForm,
  formToEventFilter,
} from '@/helpers/entities/event-filter/rule/form';
import {
  isChangeEntityEventFilterRuleType,
  isEnrichmentEventFilterRuleType,
} from '@/helpers/entities/event-filter/rule/entity';
import { isApiExternalDataType } from '@/helpers/entities/shared/external-data/entity';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesTemplateValidatorMixin } from '@/mixins/entities/template-validator';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: { EventFilterForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesTemplateValidatorMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      checking: false,
      form: eventFilterToForm(this.modal.config.rule, this.$system.timezone),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createEventFilter.create.title');
    },

    isEnrichment() {
      return isEnrichmentEventFilterRuleType(this.form.type);
    },

    isChangeEntity() {
      return isChangeEntityEventFilterRuleType(this.form.type);
    },

    hasVariablesFields() {
      return this.isEnrichment || isChangeEntityEventFilterRuleType(this.form.type);
    },
  },
  watch: {
    'form.type': {
      handler() {
        this.errors.clear();
      },
    },
  },
  methods: {
    async validateRequestTemplates(request) {
      const [url, payload] = await this.validateEventFilterRulesVariables({
        data: [
          { text: request.url },
          { text: request.payload },
        ],
      });

      return {
        url,
        payload,
      };
    },

    async validateChangeEntityTemplates(config) {
      const [component, connector, connectorName, resource] = await this.validateEventFilterRulesVariables({
        data: [
          { text: config.component },
          { text: config.connector },
          { text: config.connector_name },
          { text: config.resource },
        ],
      });

      return {
        component,
        connector,
        connector_name: connectorName,
        resource,
      };
    },

    async validateArrayWithValueTemplates(array) {
      const errors = await this.validateEventFilterRulesVariables({
        data: array.map(({ value }) => ({ text: value })),
      });

      return errors.map(value => ({ value }));
    },

    validateExternalDataTemplates(externalData) {
      return Promise.all(externalData.map(async ({ type, conditions, request }) => {
        if (isApiExternalDataType(type)) {
          return {
            request: await this.validateRequestTemplates(request),
          };
        }

        return {
          conditions: await this.validateArrayWithValueTemplates(conditions),
        };
      }));
    },

    async validateConfigTemplates(config) {
      const [actions, fields] = await Promise.all([
        this.isEnrichment ? this.validateArrayWithValueTemplates(config.actions) : undefined,
        this.isChangeEntity ? this.validateChangeEntityTemplates(config) : undefined,
      ]);

      return {
        actions,
        ...fields,
      };
    },

    async validateFormTemplates(form) {
      const [externalData, config] = await Promise.all([
        this.validateExternalDataTemplates(form.external_data),
        this.validateConfigTemplates(form.config),
      ]);

      return {
        external_data: externalData,
        config,
      };
    },

    async validateTemplateVariables() {
      this.checking = true;

      try {
        const errors = await this.validateFormTemplates(this.form);

        const wasSet = this.setFormErrors(eventFilterRuleTemplateVariablesErrorsToForm(errors, this.form));

        if (!wasSet) {
          this.$popups.success({ text: this.$t('eventFilter.syntaxIsValid') });
        }
      } catch (err) {
        console.error(err);
      } finally {
        this.checking = false;
      }
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToEventFilter(this.form, this.$system.timezone));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

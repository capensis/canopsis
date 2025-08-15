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
          @click="close"
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
          :disabled="isDisabled"
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
import { computed, inject, ref, watch } from 'vue';

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

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';
import { useTemplateValidator } from '@/hooks/store/modules/template-validator';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilter,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EventFilterForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const popups = usePopups();
    const system = inject('$system');

    const form = ref(eventFilterToForm(config.value.rule, system.timezone));
    const checking = ref(false);

    const { validator, setFormErrors } = useValidationFormErrors(form);

    const { validateEventFilterRulesVariables } = useTemplateValidator();

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToEventFilter(form.value, system.timezone));
        }

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    const title = computed(() => config.value.title ?? t('modals.createEventFilter.create.title'));

    const isEnrichment = computed(() => isEnrichmentEventFilterRuleType(form.value.type));

    const isChangeEntity = computed(() => isChangeEntityEventFilterRuleType(form.value.type));

    const hasVariablesFields = computed(() => (
      isEnrichment.value || isChangeEntityEventFilterRuleType(form.value.type)
    ));

    /**
     * Validates request templates for API external data
     *
     * @param {Object} request - Request object containing url and payload
     * @param {string} request.url - Request URL template
     * @param {string} request.payload - Request payload template
     * @returns {Promise<Object>} Object containing validation errors for url and payload
     */
    const validateRequestTemplates = async (request) => {
      const [url, payload] = await validateEventFilterRulesVariables({
        data: [
          { text: request.url },
          { text: request.payload },
        ],
      });

      return {
        url,
        payload,
      };
    };

    /**
     * Validates change entity templates for entity modification
     *
     * @param {Object} configData - Configuration data containing entity fields
     * @param {string} configData.component - Component template
     * @param {string} configData.connector - Connector template
     * @param {string} configData.connector_name - Connector name template
     * @param {string} configData.resource - Resource template
     * @returns {Promise<Object>} Object containing validation errors for entity fields
     */
    const validateChangeEntityTemplates = async (configData) => {
      const [component, connector, connectorName, resource] = await validateEventFilterRulesVariables({
        data: [
          { text: configData.component },
          { text: configData.connector },
          { text: configData.connector_name },
          { text: configData.resource },
        ],
      });

      return {
        component,
        connector,
        connector_name: connectorName,
        resource,
      };
    };

    /**
     * Validates an array of objects containing value templates
     *
     * @param {Array<Object>} array - Array of objects with value property
     * @param {string} array[].value - Template value to validate
     * @returns {Promise<Array<Object>>} Array of validation errors for each value
     */
    const validateArrayWithValueTemplates = async (array) => {
      const errors = await validateEventFilterRulesVariables({
        data: array.map(({ value }) => ({ text: value })),
      });

      return errors.map(value => ({ value }));
    };

    /**
     * Validates external data templates based on their type
     *
     * @param {Array<Object>} externalData - Array of external data configurations
     * @param {string} externalData[].type - Type of external data (API or other)
     * @param {Array<Object>} [externalData[].conditions] - Conditions array for non-API types
     * @param {Object} [externalData[].request] - Request configuration for API types
     * @returns {Promise<Array<Object>>} Array of validation errors for each external data item
     */
    const validateExternalDataTemplates = externalData => Promise.all(
      externalData.map(async ({ type, conditions, request }) => {
        if (isApiExternalDataType(type)) {
          return {
            request: await validateRequestTemplates(request),
          };
        }

        return {
          conditions: await validateArrayWithValueTemplates(conditions),
        };
      }),
    );

    /**
     * Validates configuration templates based on event filter rule type
     *
     * @param {Object} configData - Configuration data object
     * @param {Array<Object>} [configData.actions] - Actions array for enrichment rules
     * @param {string} [configData.component] - Component field for change entity rules
     * @param {string} [configData.connector] - Connector field for change entity rules
     * @param {string} [configData.connector_name] - Connector name for change entity rules
     * @param {string} [configData.resource] - Resource field for change entity rules
     * @returns {Promise<Object>} Object containing validation errors for configuration fields
     */
    const validateConfigTemplates = async (configData) => {
      const [actions, fields] = await Promise.all([
        isEnrichment.value ? validateArrayWithValueTemplates(configData.actions) : undefined,
        isChangeEntity.value ? validateChangeEntityTemplates(configData) : undefined,
      ]);

      return {
        actions,
        ...fields,
      };
    };

    /**
     * Validates all form templates including external data and configuration
     *
     * @param {Object} formData - Complete form data object
     * @param {Array<Object>} formData.external_data - External data configurations
     * @param {Object} formData.config - Configuration object
     * @returns {Promise<Object>} Object containing validation errors for all form fields
     */
    const validateFormTemplates = async (formData) => {
      const [externalData, configData] = await Promise.all([
        validateExternalDataTemplates(formData.external_data),
        validateConfigTemplates(formData.config),
      ]);

      return {
        external_data: externalData,
        config: configData,
      };
    };

    /**
     * Validates template variables in the form and displays results to user
     * This is the main validation function that triggers syntax checking for all templates
     * in the form. Sets form errors if validation fails or shows success message.
     *
     * @returns {Promise<void>}
     */
    const validateTemplateVariables = async () => {
      checking.value = true;

      try {
        const errors = await validateFormTemplates(form.value);

        const wasSet = setFormErrors(eventFilterRuleTemplateVariablesErrorsToForm(errors, form.value));

        if (!wasSet) {
          popups.success({ text: t('eventFilter.syntaxIsValid') });
        }
      } catch (err) {
        console.error(err);
      } finally {
        checking.value = false;
      }
    };

    watch(() => form.value.type, () => {
      if (validator?.errors) {
        validator.errors.clear();
      }
    });

    return {
      config,

      form,
      checking,

      isDisabled,
      submitting,

      title,
      isEnrichment,
      isChangeEntity,
      hasVariablesFields,

      validateTemplateVariables,
      submit,
      close,
    };
  },
};
</script>

<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        {{ title }}
      </template>

      <template #text>
        <c-progress-overlay :pending="pending" />
        <service-form
          :form="values"
          :prepare-state-setting-form="prepareStateSettingForm"
          :template-vars="templateVars"
          @update:form="updateForm"
        />
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted } from 'vue';

import { ENTITY_TYPES, MODALS } from '@/constants';

import { formToService, serviceToForm } from '@/helpers/entities/service/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useTemplateVarsList } from '@/hooks/vars/template';

import ServiceForm from '@/components/other/service/form/service-form.vue';

export default {
  name: MODALS.createService,

  components: {
    ServiceForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const item = computed(() => config.value.item);
    const title = computed(() => config.value.title);
    const form = useForm({
      initialValues: serviceToForm(item.value),
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateForm = value => (values.value = value);
    const submitLabel = computed(() => t(item.value?._id ? 'common.save' : 'common.create'));

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToService(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    const { vars: templateVars, pending, fetchList } = useTemplateVarsList({
      type: ENTITY_TYPES.service,
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    /**
     * Prepares a service object for state setting form by transforming it to the proper format
     *
     * @param {Object} service - The service object to prepare
     * @param {string} service._id - The service identifier
     * @returns {Object} The prepared service object with form data, entity type, and ID
     */
    const prepareStateSettingForm = service => ({
      ...formToService(service),
      type: ENTITY_TYPES.service,
      _id: service._id,
    });

    onMounted(fetchList);

    return {
      t,
      title,
      values,
      templateVars,
      pending,
      submitting,
      submitLabel,
      closeModal,
      prepareStateSettingForm,
      updateForm,
      submit,
    };
  },
};
</script>

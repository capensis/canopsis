<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="pa-0" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <template-testing-test-variables-wrapper
          :form="values"
          :rule-id="remediationJobId"
          :type="type"
          @update:form="updateValues"
        >
          <template #default="{ templateVars }">
            <div class="pa-3">
              <remediation-job-form
                :with-payload="withPayload"
                :with-query="withQuery"
                :template-vars="templateVars"
              />
            </div>
          </template>
        </template-testing-test-variables-wrapper>
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
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
import { computed, watch } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import {
  formToRemediationJob,
  remediationJobToForm,
} from '@/helpers/entities/remediation/job/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useInfo } from '@/hooks/store/modules/info';

import RemediationJobForm from '@/components/other/remediation/jobs/form/remediation-job-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

export default {
  components: {
    RemediationJobForm,
    TemplateTestingTestVariablesWrapper,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.job;

    const { t } = useI18n();
    const modals = useModals();
    const { remediationJobConfigTypes } = useInfo();
    const config = computed(() => props.modal.config ?? {});
    const remediationJob = computed(() => config.value.remediationJob);
    const isNew = computed(() => !remediationJob.value);
    const title = computed(
      () => config.value.title ?? t('modals.createRemediationJob.create.title'),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: remediationJobToForm(remediationJob.value),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const remediationJobId = computed(() => remediationJob.value?._id);
    const remediationJobConfigType = computed(() => remediationJobConfigTypes.value.find(({ name }) => name === values.value.config.type),
    );

    const withPayload = computed(() => remediationJobConfigType.value?.with_body);
    const withQuery = computed(() => remediationJobConfigType.value?.with_query);

    const closeModal = () => modals.hide({ id: props.modal.id });
    const updateValues = (value) => {
      values.value = value;
    };

    const submit = handleSubmit(async () => {
      const data = await config.value.action?.(formToRemediationJob(values.value));

      closeModal();

      return data;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    watch(
      remediationJobConfigType,
      (newConfigType) => {
        values.value.configType = newConfigType;
      },
      { immediate: true },
    );

    return {
      t,
      values,
      remediationJobId,
      type,
      title,
      isSubmitting,
      submitLabel,
      submit,
      closeModal,
      updateValues,
      withPayload,
      withQuery,
    };
  },
};
</script>

<template>
  <v-form>
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <template-testing-test-variables-wrapper
          v-model="form"
          :is-new="isNew"
          :type="type"
        >
          <template #default="{ templateVars }">
            <remediation-job-form
              v-model="form"
              :with-payload="withPayload"
              :with-query="withQuery"
              :template-vars="templateVars"
            />
          </template>
        </template-testing-test-variables-wrapper>
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
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
          @click="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { formToRemediationJob, remediationJobToForm } from '@/helpers/entities/remediation/job/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useStore } from '@/hooks/store';

import RemediationJobForm from '@/components/other/remediation/jobs/form/remediation-job-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createRemediationJob,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
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

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();
    const store = useStore();

    const form = ref(remediationJobToForm(config.value.remediationJob));

    const isNew = computed(() => !config.value.remediationJob?._id);
    const title = computed(() => config.value.title ?? t('modals.createRemediationJob.create.title'));

    const remediationJobConfigTypes = computed(() => store.getters['info/remediationJobConfigTypes']);
    const remediationJobConfigType = computed(() => remediationJobConfigTypes.value.find(
      ({ name }) => name === form.value.config.type,
    ));

    const withPayload = computed(() => remediationJobConfigType.value?.with_body);
    const withQuery = computed(() => remediationJobConfigType.value?.with_query);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const data = await config.value.action?.(formToRemediationJob(form.value));

        close();

        return data;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    watch(remediationJobConfigType, newConfigType => form.value.configType = newConfigType, { immediate: true });

    return {
      form,
      config,
      isNew,
      type,
      title,
      isDisabled,
      submitting,
      submit,
      withPayload,
      withQuery,
    };
  },
};
</script>

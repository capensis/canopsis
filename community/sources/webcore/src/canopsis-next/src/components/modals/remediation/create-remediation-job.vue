<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <template-testing-test-variables-wrapper
        v-field="form"
        :is-new="isNew"
        :type="type"
      >
        <template #default="{ templateVars }">
          <v-form>
            <remediation-job-form
              v-model="form"
              :with-payload="withPayload"
              :with-query="withQuery"
              :template-vars="templateVars"
            />
          </v-form>
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
</template>

<script>
import { computed, ref } from 'vue';
import { createNamespacedHelpers } from 'vuex';

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

const { mapGetters: mapInfoGetters } = createNamespacedHelpers('info');

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

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        // Calculate remediationJobConfigType from store
        const remediationJobConfigTypes = store.getters['info/remediationJobConfigTypes'];
        const remediationJobConfigType = remediationJobConfigTypes.find(({ name }) => name === form.value.config.type);

        await config.value.action?.(formToRemediationJob(form.value, remediationJobConfigType));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      isDisabled,
      submitting,
      submit,
    };
  },
  computed: {
    ...mapInfoGetters(['remediationJobConfigTypes']),

    remediationJobConfigType() {
      return this.remediationJobConfigTypes.find(({ name }) => name === this.form.config.type);
    },

    withPayload() {
      return this.remediationJobConfigType?.with_body;
    },

    withQuery() {
      return this.remediationJobConfigType?.with_query;
    },
  },
};
</script>

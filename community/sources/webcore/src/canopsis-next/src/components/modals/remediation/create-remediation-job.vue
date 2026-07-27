<template>
  <v-form>
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <remediation-job-form
          v-model="form"
          :rule-id="remediationJobId"
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
          :disabled="submitting"
          :loading="submitting"
          class="primary"
          type="submit"
          @click="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToRemediationJob, remediationJobToForm } from '@/helpers/entities/remediation/job/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';

import RemediationJobForm from '@/components/other/remediation/jobs/form/remediation-job-form.vue';

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
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(remediationJobToForm(config.value.remediationJob));

    const remediationJobId = computed(() => config.value.remediationJob?._id);
    const title = computed(() => config.value.title ?? t('modals.createRemediationJob.create.title'));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.remediationJob,
      method: async () => {
        const data = await config.value.action?.(formToRemediationJob(form.value));

        close();

        return data;
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      remediationJobId,
      title,
      submitting,
      submitLabel,
      submit,
    };
  },
};
</script>

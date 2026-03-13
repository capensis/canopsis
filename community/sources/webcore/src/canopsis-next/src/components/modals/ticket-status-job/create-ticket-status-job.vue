<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-3" column>
          <v-layout class="gap-3" wrap>
            <v-flex md4 xs12>
              <v-text-field
                :value="ruleTypeDisplay"
                :label="$t('jobs.ruleType')"
                readonly
                outlined
                dense
              />
            </v-flex>
            <v-flex xs12 md4>
              <v-text-field
                :value="form.rule_name"
                :label="$t('jobs.ruleName')"
                readonly
                outlined
                dense
              />
            </v-flex>
            <v-flex xs12 md4>
              <v-text-field
                :value="form.ticket_system_name"
                :label="$t('jobs.ticketSystemName')"
                readonly
                outlined
                dense
              />
            </v-flex>
          </v-layout>
          <v-flex xs12>
            <v-text-field
              v-field="form.ticket_id"
              v-validate="ticketRules"
              :label="$tc('common.ticket')"
              :error-messages="errors.collect('ticket')"
              name="ticket"
              required
              outlined
              dense
            />
          </v-flex>
        </v-layout>
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

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { ticketStatusJobToForm, formToTicketStatusJob } from '@/helpers/entities/ticket-status-job/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';

import ModalWrapper from '../modal-wrapper.vue';

const ticketRules = { required: true };

export default {
  name: MODALS.createTicketStatusJob,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
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

    const job = computed(() => config.value.ticketStatusJob ?? {});
    const form = ref(ticketStatusJobToForm(job.value));

    const ruleTypeDisplay = computed(() => {
      const ruleType = form.value.rule_type;
      if (!ruleType) return '';

      return t(`jobs.ruleTypeValues.${ruleType}`) || ruleType;
    });

    const jobName = computed(() => {
      const j = job.value;
      const fallback = `${j?.ticket_system_name || ''} - ${j?.ticket_id ?? ''}`.trim() || '-';

      return (j?.rule_name ?? fallback);
    });

    const title = computed(
      () => config.value.title ?? t('modals.createTicketStatusJob.edit.title', { jobName: jobName.value }),
    );

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToTicketStatusJob(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    watch(job, (newJob) => {
      form.value = ticketStatusJobToForm(newJob);
    }, { immediate: true });

    return {
      form,
      ruleTypeDisplay,
      ticketRules,
      title,
      isDisabled,
      submitting,
      submit,
      close,
    };
  },
};

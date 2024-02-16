<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <remediation-job-form
          v-model="form"
          :with-payload="withPayload"
          :with-query="withQuery"
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
import { createNamespacedHelpers } from 'vuex';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToRemediationJob, remediationJobToForm } from '@/helpers/entities/remediation/job/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationJobForm from '@/components/other/remediation/jobs/form/remediation-job-form.vue';

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
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: remediationJobToForm(this.modal.config.remediationJob),
    };
  },
  computed: {
    ...mapInfoGetters(['remediationJobConfigTypes']),

    title() {
      return this.config.title ?? this.$t('modals.createRemediationJob.create.title');
    },

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
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToRemediationJob(this.form, this.remediationJobConfigType));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

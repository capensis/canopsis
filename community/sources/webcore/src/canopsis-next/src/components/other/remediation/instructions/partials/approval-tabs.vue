<template>
  <v-tabs
    class="mt-3"
    slider-color="primary"
    fixed-tabs
  >
    <template v-if="updated">
      <v-tab>{{ $t('modals.remediationInstructionApproval.tabs.updated') }}</v-tab>
      <v-tab-item class="pt-3">
        <remediation-instruction-form
          :form="updatedForm"
          disabled-common
          disabled
        />
      </v-tab-item>
    </template>
    <v-tab>{{ $t('modals.remediationInstructionApproval.tabs.original') }}</v-tab>
    <v-tab-item class="pt-3">
      <remediation-instruction-form
        :form="originalForm"
        disabled-common
        disabled
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { remediationInstructionToForm } from '@/helpers/entities/remediation/instruction/form';

import RemediationInstructionForm from '../form/remediation-instruction-form.vue';

export default {
  components: { RemediationInstructionForm },
  props: {
    original: {
      type: Object,
      required: true,
    },
    updated: {
      type: Object,
      required: false,
    },
  },
  computed: {
    originalForm() {
      return remediationInstructionToForm(this.original);
    },
    updatedForm() {
      return remediationInstructionToForm(this.updated);
    },
  },
};
</script>

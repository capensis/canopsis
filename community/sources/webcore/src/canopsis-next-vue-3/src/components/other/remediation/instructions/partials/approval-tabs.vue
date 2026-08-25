<template>
  <div class="mt-3">
    <v-tabs v-model="activeTab" color="primary" fixed-tabs>
      <v-tab v-if="updated" value="updated">
        {{ $t('modals.remediationInstructionApproval.tabs.updated') }}
      </v-tab>
      <v-tab value="original">
        {{ $t('modals.remediationInstructionApproval.tabs.original') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item v-if="updated" value="updated" class="pt-3">
        <remediation-instruction-general-form
          :form="updatedForm"
          disabled-common
          disabled
        />
      </v-window-item>
      <v-window-item value="original" class="pt-3">
        <remediation-instruction-general-form
          :form="originalForm"
          disabled-common
          disabled
        />
      </v-window-item>
    </v-window>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue';

import { remediationInstructionToForm } from '@/helpers/entities/remediation/instruction/form';

import RemediationInstructionGeneralForm from '../form/remediation-instruction-general-form.vue';

const props = defineProps({
  original: {
    type: Object,
    required: true,
  },
  updated: {
    type: Object,
    required: false,
  },
});

const activeTab = ref(props.updated ? 'updated' : 'original');
const originalForm = computed(() => remediationInstructionToForm(props.original));
const updatedForm = computed(() => remediationInstructionToForm(props.updated));
</script>

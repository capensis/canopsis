<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.liveReporting.editLiveReporting') }}</span>
      </template>

      <template #text>
        <h3>{{ t('modals.liveReporting.dateInterval') }}</h3>
        <date-interval-selector
          :value="values"
          :quick-ranges="quickRanges"
          @update:value="updateForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn :loading="submitting" color="primary" type="submit">
          {{ t('common.apply') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed } from 'vue';

import { LIVE_REPORTING_QUICK_RANGES, MODALS } from '@/constants';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import DateIntervalSelector from '@/components/forms/date-interval-selector.vue';

defineOptions({ name: MODALS.editLiveReporting });

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();
const modals = useModals();
const config = computed(() => props.modal.config ?? {});
const form = useForm({
  initialValues: {
    tstart: config.value.tstart ?? '',
    tstop: config.value.tstop ?? '',
    time_field: config.value.time_field ?? '',
  },
});
const { values } = form;
const quickRanges = Object.values(LIVE_REPORTING_QUICK_RANGES);

const closeModal = () => modals.hide({ id: props.modal.id });
const updateForm = value => form.reset(value);
const submit = form.handleSubmit(async () => {
  await config.value.action?.(values.value);
  closeModal();
});

const { submitting } = form;
</script>

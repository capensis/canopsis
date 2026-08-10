<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close minimize>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <c-layout gap="3" column>
          <c-alert type="warning">
            {{ config.warningText }}
          </c-alert>

          <maintenance-form />
        </c-layout>
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :loading="submitting"
          :disabled="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script setup>
import { computed, inject } from 'vue';

import MaintenanceForm from '@/components/other/maintenance/form/maintenance-form.vue';
import { maintenanceToForm } from '@/helpers/entities/maintenance/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const modalContext = inject('$modal');

const modals = useModals();
const { t } = useI18n();

const config = computed(() => props.modal.config ?? {});
const form = useForm({
  initialValues: maintenanceToForm(config.value.maintenance),
});

const title = computed(() => config.value.title ?? t('modals.createMaintenance.setup.title'));

const submitLabel = computed(() =>
  config.value.maintenance ? t('common.submit') : t('modals.createMaintenance.enableMaintenance'),
);

const close = () => modals.hide({ id: modalContext.id });

const { handleSubmit, isSubmitting } = form;

const submit = handleSubmit(async () => {
  await config.value.action?.(form.values.value);

  close();
});

useFormConfirmableCloseModal({ form: form.values, submit, close });

const submitting = computed(() => isSubmitting.value);
</script>

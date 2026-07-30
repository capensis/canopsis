<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <pbehavior-reason-form />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
:disabled="submitting"
:loading="submitting"
color="primary"
type="submit">
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import {
  formToPbehaviorReason,
  pbehaviorReasonToForm,
} from '@/helpers/entities/pbehavior/reason/form';

import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import PbehaviorReasonForm from '@/components/other/pbehavior/reasons/form/pbehavior-reason-form.vue';

export default {
  components: { PbehaviorReasonForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const config = computed(() => props.modal.config ?? {});
    const reason = computed(() => config.value.pbehaviorReason ?? config.value.reason);
    const isNew = computed(() => !reason.value?._id);
    const title = computed(
      () =>
        config.value.title ||
        t(`modals.createPbehaviorReason.${isNew.value ? 'create' : 'edit'}.title`),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: pbehaviorReasonToForm(reason.value),
    });
    const { values, handleSubmit, submitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToPbehaviorReason(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    return {
      t,
      title,
      submitLabel,
      submitting,
      closeModal,
      submit,
    };
  },
};
</script>

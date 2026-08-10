<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <pbehavior-type-form
          :only-color="onlyColor"
          :pending-priority="pendingPriority"
          :sync-original-form="updateOriginalForm"
        />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="submitting"
          :loading="submitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import PbehaviorTypeForm from '@/components/other/pbehavior/types/form/pbehavior-type-form.vue';
import { formToPbehaviorType, pbehaviorTypeToForm } from '@/helpers/entities/pbehavior/type/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

export default {
  components: { PbehaviorTypeForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchNextPbehaviorTypePriority } = usePbehaviorType();
    const config = computed(() => props.modal.config ?? {});
    const pbehaviorType = computed(() => config.value.pbehaviorType ?? config.value.type);
    const isNew = computed(() => !pbehaviorType.value?._id);
    const onlyColor = computed(() => !!pbehaviorType.value?.default);
    const title = computed(
      () =>
        config.value.title ||
        t(`modals.createPbehaviorType.${isNew.value ? 'create' : 'edit'}.title`),
    );
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));
    const pendingPriority = ref(false);

    const form = useForm({
      initialValues: pbehaviorTypeToForm(pbehaviorType.value),
    });
    const { values, handleSubmit, submitting, reset } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToPbehaviorType(values.value));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });

    onMounted(async () => {
      if (!isNew.value) {
        return;
      }

      try {
        pendingPriority.value = true;

        const { data } = await fetchNextPbehaviorTypePriority();

        reset({
          ...values.value,

          priority: data,
        });
        updateOriginalForm();
      } finally {
        pendingPriority.value = false;
      }
    });

    return {
      t,
      title,
      submitLabel,
      onlyColor,
      pendingPriority,
      submitting,
      updateOriginalForm,
      closeModal,
      submit,
    };
  },
};
</script>

<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <role-form :with-template="config.withTemplate" :is-new="isNew" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
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
import { computed } from 'vue';

import { roleToForm, formToRole } from '@/helpers/entities/role/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import RoleForm from '@/components/other/role/form/role-form.vue';

export default {
  components: { RoleForm },

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

    const isNew = computed(() => !config.value.role?._id);
    const title = computed(() => config.value.title || t('modals.createRole.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: roleToForm(config.value.role),
    });
    const { handleSubmit, isSubmitting, values } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToRole(values.value));

      closeModal();
    });

    useFormConfirmableCloseModal({ form: values, submit, close: closeModal });

    return {
      t,
      config,
      isNew,
      title,
      submitLabel,
      isSubmitting,
      closeModal,
      submit,
    };
  },
};
</script>

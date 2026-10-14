<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <user-form :is-new="isNew" :user="config.user" :only-user-prefs="config.onlyUserPrefs" />
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
import { computed, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { userToForm, formToUserRequest, findDefaultViewId } from '@/helpers/entities/user/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useView } from '@/hooks/store/modules/view';

import UserForm from '@/components/other/users/form/user-form.vue';

export default {
  components: { UserForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const { viewGroups, fetchViewGroupsList } = useView();
    const config = computed(() => props.modal.config ?? {});

    const isNew = computed(() => !config.value.user);
    const title = computed(() => config.value.title || t('modals.createUser.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.edit'));

    const form = useForm({
      initialValues: userToForm(config.value.user),
    });
    const { handleSubmit, isSubmitting, values } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      await config.value.action?.(formToUserRequest(values.value));

      closeModal();
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({ form: values, submit, close: closeModal });

    onMounted(async () => {
      if (!isNew.value || values.value.defaultview) {
        return;
      }

      try {
        if (!viewGroups.value?.length) {
          await fetchViewGroupsList({
            params: {
              limit: MAX_LIMIT,
              page: 1,
              with_views: true,
              with_tabs: true,
              with_widgets: true,
              with_flags: true,
              with_private: true,
            },
          });
        }

        const defaultViewId = findDefaultViewId(viewGroups.value);

        if (!defaultViewId) {
          return;
        }

        form.setValue('defaultview', defaultViewId);
        updateOriginalForm();
      } catch (err) {
        console.error(err);
      }
    });

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

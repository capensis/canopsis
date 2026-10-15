<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <user-form
          v-model="form"
          :is-new="isNew"
          :user="config.user"
          :only-user-prefs="config.onlyUserPrefs"
        />
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
          :disabled="submitting"
          :loading="submitting"
          class="primary white--text"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { MODALS } from '@/constants';

import { userToForm, formToUserRequest, findDefaultViewId } from '@/helpers/entities/user/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useViewGroup } from '@/hooks/store/modules/view';

import UserForm from '@/components/other/users/form/user-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createUser,

  $_veeValidate: {
    validator: 'new',
  },

  components: { UserForm, ModalWrapper },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const { groups, fetchAllGroupsListWithWidgets } = useViewGroup();

    const form = ref(userToForm(config.value.user));

    const isNew = computed(() => !config.value.user);
    const title = computed(() => (config.value.title || t('modals.createUser.create.title')));

    const { submit, submitting, submitLabel } = useSubmittableForm({
      form,
      item: config.value.user,
      method: async () => {
        await config.value.action(formToUserRequest(form.value));

        close();
      },
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({ form, submit, close });

    onMounted(async () => {
      if (!isNew.value || form.value.defaultview) {
        return;
      }

      try {
        if (!groups.value?.length) {
          await fetchAllGroupsListWithWidgets();
        }

        const defaultViewId = findDefaultViewId(groups.value);

        if (!defaultViewId) {
          return;
        }

        form.value = {
          ...form.value,
          defaultview: defaultViewId,
        };
        updateOriginalForm();
      } catch (err) {
        console.error(err);
      }
    });

    return {
      form,
      isNew,
      title,
      submitting,
      submitLabel,
      submit,
      close,
      config,
    };
  },
};
</script>

<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper
      text-class="position-relative pa-0"
      close
    >
      <template #title="">
        <span>{{ $t('userInterface.title') }}</span>
      </template>
      <template #text="">
        <c-progress-overlay :pending="pending" />
        <user-interface-form
          v-model="form"
          ref="userForm"
          :disabled="!hasUpdateAccess"
          class="pa-3"
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
          v-if="hasUpdateAccess"
          :disabled="submitting"
          :loading="submitting"
          class="primary white--text"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS, VALIDATION_DELAY } from '@/constants';

import { userInterfaceToForm } from '@/helpers/entities/user-interface/form';
import { notificationsSettingsToForm } from '@/helpers/entities/notification/form';
import { getFileDataUrlContent } from '@/helpers/file/file-select';

import { useCRUDPermissions } from '@/hooks/auth';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useInfo } from '@/hooks/store/modules/info';
import { useStoreModuleHooks } from '@/hooks/store';

import UserInterfaceForm from '@/components/other/user-interface/form/user-interface-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.userInterface,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    UserInterfaceForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const userForm = ref(null);
    const isFormLoading = ref(true);
    const form = ref({
      ...userInterfaceToForm(),
      ...notificationsSettingsToForm(),
    });

    const { t } = useI18n();
    const popups = usePopups();
    const { close } = useInnerModal(props);

    const { hasUpdateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.parameters);

    const {
      appInfo,
      appInfoPending,
      fetchAppInfo,
      updateUserInterface,
      setTitle,
    } = useInfo();

    const { useActions } = useStoreModuleHooks('notificationSettings');
    const {
      fetchNotificationSettings,
      updateNotificationSettings,
    } = useActions({
      fetchNotificationSettings: 'fetchItemWithoutStore',
      updateNotificationSettings: 'update',
    });

    const pending = computed(() => appInfoPending.value || isFormLoading.value);

    const { submitting, submit } = useSubmittableForm({
      form,
      method: async () => {
        const { instruction, ...userInterfaceFormData } = form.value;
        const data = { ...userInterfaceFormData };

        if (form.value.logo) {
          data.logo = await getFileDataUrlContent(form.value.logo);
        }

        await updateUserInterface({ data });
        await updateNotificationSettings({ data: { instruction } });
        await fetchAppInfo();

        await setTitle();

        popups.success({ text: t('success.default') });

        close();
      },
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({ form, submit, close });

    const loadForm = async () => {
      isFormLoading.value = true;

      try {
        const notificationsSettings = await fetchNotificationSettings();

        form.value = {
          ...userInterfaceToForm(appInfo.value),
          ...notificationsSettingsToForm(notificationsSettings),
        };

        updateOriginalForm();
      } finally {
        isFormLoading.value = false;
      }
    };

    onMounted(async () => {
      await fetchAppInfo();
      await loadForm();
    });

    return {
      userForm,
      form,
      pending,
      hasUpdateAccess,
      submitting,
      submit,
      close,
    };
  },
};
</script>

<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="pa-0" close>
      <template #title>
        <span>{{ t('userInterface.title') }}</span>
      </template>

      <template #text>
        <c-progress-overlay :pending="pending" />

        <user-interface-form :disabled="!hasUpdateAccess" class="pa-3" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          v-if="hasUpdateAccess"
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ t('common.save') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { getFileDataUrlContent } from '@/helpers/file/file-select';
import { notificationsSettingsToForm } from '@/helpers/entities/notification/form';
import { userInterfaceToForm } from '@/helpers/entities/user-interface/form';

import { useCRUDPermissions } from '@/hooks/auth';
import { useForm } from '@/hooks/form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useInfo } from '@/hooks/store/modules/info';
import { useNotificationSettings } from '@/hooks/store/modules/notification-settings';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import UserInterfaceForm from '@/components/other/user-interface/form/user-interface-form.vue';

export default {
  components: {
    UserInterfaceForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const isFormLoading = ref(true);

    const { hasUpdateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.parameters);

    const { appInfo, appInfoPending, fetchAppInfo, updateUserInterface, setTitle } = useInfo();
    const { fetchNotificationSettingsWithoutStore, updateNotificationSettings } =
      useNotificationSettings();
    const { fetchTemplateVarsList } = useTemplateVars();

    const buildForm = (appInfoValue, notificationSettingsValue) => ({
      ...userInterfaceToForm(appInfoValue),
      version_description:
        appInfoValue?.version_description || t('userInterface.defaultVersionDescription'),
      ...notificationsSettingsToForm(notificationSettingsValue),
    });

    const form = useForm({
      initialValues: buildForm(),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const pending = computed(() => appInfoPending.value || isFormLoading.value);
    const close = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const { instruction, logo, ...userInterfaceData } = values.value;
      const data = { ...userInterfaceData };

      if (logo) {
        data.logo = await getFileDataUrlContent(logo);
      }

      await updateUserInterface({ data });
      await updateNotificationSettings({ data: { instruction } });
      await fetchAppInfo();
      await setTitle();

      popups.success({ text: t('success.default') });

      close();
    });

    const { updateOriginalForm } = useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    const loadForm = async () => {
      isFormLoading.value = true;

      try {
        const notificationSettings = await fetchNotificationSettingsWithoutStore();

        form.reset(buildForm(appInfo.value, notificationSettings));
        updateOriginalForm();
      } finally {
        isFormLoading.value = false;
      }
    };

    onMounted(async () => {
      await Promise.all([fetchAppInfo(), fetchTemplateVarsList()]);
      await loadForm();
    });

    return {
      t,
      pending,
      hasUpdateAccess,
      isSubmitting,

      submit,
      close,
    };
  },
};
</script>

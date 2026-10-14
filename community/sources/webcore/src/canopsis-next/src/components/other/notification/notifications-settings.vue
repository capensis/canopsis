<template>
  <c-layout class="my-2" justify="center">
    <v-progress-circular v-if="pending" color="primary" indeterminate />

    <v-row v-else justify="center" no-gutters>
      <v-col cols="12" md="10">
        <v-form @submit.prevent="submit">
          <notifications-settings-form />

          <v-divider class="mt-3" />

          <c-layout class="mt-3" justify="end">
            <v-btn
              :disabled="isSubmitting"
              :loading="isSubmitting"
              color="primary"
              type="submit"
            >
              {{ t('common.submit') }}
            </v-btn>
          </c-layout>
        </v-form>
      </v-col>
    </v-row>
  </c-layout>
</template>

<script>
import { onMounted, ref } from 'vue';

import { notificationsSettingsToForm } from '@/helpers/entities/notification/form';

import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useNotificationSettings } from '@/hooks/store/modules/notification-settings';

import NotificationsSettingsForm from './form/notifications-settings-form.vue';

export default {
  components: {
    NotificationsSettingsForm,
  },

  setup() {
    const { t } = useI18n();
    const popups = usePopups();
    const pending = ref(true);

    const { fetchNotificationSettingsWithoutStore, updateNotificationSettings } =
      useNotificationSettings();

    const form = useForm({
      initialValues: notificationsSettingsToForm(),
    });
    const { values, handleSubmit, isSubmitting } = form;

    /**
     * Fetches notification settings and resets the form with server values.
     */
    const fetchNotificationsSettings = async () => {
      pending.value = true;

      try {
        const notificationsSettings = await fetchNotificationSettingsWithoutStore();

        form.reset(notificationsSettingsToForm(notificationsSettings));
      } finally {
        pending.value = false;
      }
    };

    const submit = handleSubmit(async () => {
      await updateNotificationSettings({ data: values.value });

      popups.success({ text: t('success.default') });
    });

    onMounted(fetchNotificationsSettings);

    return {
      t,
      pending,
      isSubmitting,

      submit,
    };
  },
};
</script>

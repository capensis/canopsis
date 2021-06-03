<template lang="pug">
  v-layout.my-2(justify-center)
    v-progress-circular(v-if="!form", indeterminate, color="primary")
    v-flex(v-else, offset-xs1, md10)
      v-form(@submit.prevent="submit")
        notifications-settings-form(v-model="form")
        v-layout.mt-3(row, justify-end)
          v-btn.primary.mr-0(
            :disabled="isDisabled",
            :loading="submitting",
            type="submit"
          ) {{ $t('common.submit') }}
</template>

<script>
import { notificationsSettingsToForm, formToNotificationsSettings } from '@/helpers/forms/notification';

import { entitiesNotificationSettingsMixin } from '@/mixins/entities/notification-settings';
import { validationErrorsMixin } from '@/mixins/form/validation-errors';
import { submittableMixin } from '@/mixins/submittable';

import NotificationsSettingsForm from './form/notifications-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    NotificationsSettingsForm,
  },
  mixins: [
    submittableMixin(),
    entitiesNotificationSettingsMixin,
    validationErrorsMixin(),
  ],
  data() {
    return {
      form: null,
    };
  },
  mounted() {
    this.fetchNotificationsSettings();
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          await this.updateNotificationSettings({ data: formToNotificationsSettings(this.form) });

          this.$popups.success({ text: this.$t('success.default') });
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },

    async fetchNotificationsSettings() {
      const notificationsSettings = await this.fetchNotificationSettingsWithoutStore();

      this.form = notificationsSettingsToForm(notificationsSettings);
    },
  },
};
</script>

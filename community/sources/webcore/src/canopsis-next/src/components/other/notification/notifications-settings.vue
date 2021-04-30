<template lang="pug">
  v-layout.my-2(v-if="!form", justify-center)
    v-progress-circular(indeterminate, color="primary")
  v-form(v-else, @submit.prevent="submit")
    notifications-settings-form(v-model="form")
    v-layout.mt-3(row, justify-end)
      v-btn(
        flat,
        @click="reset"
      ) {{ $t('common.cancel') }}
      v-btn.primary.mr-0(
        type="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import { notificationsSettingsToForm, formToNotificationsSettings } from '@/helpers/forms/notification';

import NotificationsSettingsForm from './form/notifications-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    NotificationsSettingsForm,
  },
  data() {
    return {
      form: null,
    };
  },
  mounted() {
    this.fetchNotificationsSettings();
  },
  methods: {
    reset() {},

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const notificationsSettings = formToNotificationsSettings(this.form);

        console.warn(notificationsSettings);
      }
    },

    async fetchNotificationsSettings() {
      await new Promise(r => setTimeout(r, 5000));

      this.form = notificationsSettingsToForm(this.notificationsSettings);
    },
  },
};
</script>

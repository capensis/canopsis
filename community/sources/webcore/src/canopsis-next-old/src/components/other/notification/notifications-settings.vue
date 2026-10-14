<template>
  <v-layout
    class="my-2"
    justify-center
  >
    <v-progress-circular
      v-if="!form"
      color="primary"
      indeterminate
    />
    <v-flex
      v-else
      offset-xs1
      md10
    >
      <v-form @submit.prevent="submit">
        <notifications-settings-form v-model="form" />
        <v-divider class="mt-3" />
        <v-layout
          class="mt-3"
          justify-end
        >
          <v-btn
            :disabled="isDisabled"
            :loading="submitting"
            class="primary mr-0"
            type="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
        </v-layout>
      </v-form>
    </v-flex>
  </v-layout>
</template>

<script>
import { VALIDATION_DELAY } from '@/constants';

import { notificationsSettingsToForm } from '@/helpers/entities/notification/form';

import { entitiesNotificationSettingsMixin } from '@/mixins/entities/notification-settings';
import { submittableMixinCreator } from '@/mixins/submittable';

import NotificationsSettingsForm from './form/notifications-settings-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    NotificationsSettingsForm,
  },
  mixins: [
    entitiesNotificationSettingsMixin,
    submittableMixinCreator(),
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
    async fetchNotificationsSettings() {
      const notificationsSettings = await this.fetchNotificationSettingsWithoutStore();

      this.form = notificationsSettingsToForm(notificationsSettings);
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.updateNotificationSettings({ data: this.form });

        this.$popups.success({ text: this.$t('success.default') });
      }
    },
  },
};
</script>

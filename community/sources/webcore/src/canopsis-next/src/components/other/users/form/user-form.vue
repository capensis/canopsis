<template>
  <v-layout column>
    <c-id-field
      v-field="form._id"
      :disabled="onlyUserPrefs || !isNew"
    />
    <c-name-field
      v-field="form.name"
      :label="$t('common.username')"
      :disabled="onlyUserPrefs"
      autocomplete="new-password"
      required
    />
    <v-text-field
      v-field="form.firstname"
      :label="$t('user.firstName')"
      :disabled="onlyUserPrefs"
    />
    <v-text-field
      v-field="form.lastname"
      :label="$t('user.lastName')"
      :disabled="onlyUserPrefs"
    />
    <v-text-field
      v-field="form.email"
      v-validate="'required|email'"
      :label="$t('user.email')"
      :disabled="onlyUserPrefs"
      :error-messages="errors.collect('email')"
      name="email"
      autocomplete="new-password"
    />
    <c-password-field
      v-field="form.password"
      :required="isNew"
      autocomplete="new-password"
    />
    <c-role-field
      v-field="form.roles"
      :disabled="onlyUserPrefs"
      :label="$tc('common.role', 2)"
      required
      multiple
      chips
    />
    <c-language-field
      v-field="form.ui_language"
      :label="$t('user.language')"
    />
    <v-select
      class="mt-0"
      v-field="form.ui_groups_navigation_type"
      :label="$t('user.navigationType')"
      :items="groupsNavigationItems"
    />
    <c-theme-field v-field="form.ui_theme" />
    <v-layout
      v-if="!isNew"
      align-center
    >
      <div>{{ $t('common.authKey') }}: {{ user.authkey }}</div>
      <c-copy-btn
        :value="user.authkey"
        :tooltip="$t('common.copyToClipboard')"
        small
        fab
        left
        @success="showCopyAuthKeySuccessPopup"
        @error="showCopyAuthKeyErrorPopup"
      />
    </v-layout>
    <c-enabled-field
      v-field="form.enable"
      :disabled="onlyUserPrefs"
    />
    <view-selector v-field="form.defaultview" />
  </v-layout>
</template>

<script>
import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import ViewSelector from '@/components/forms/fields/view-selector.vue';

export default {
  inject: ['$validator'],
  components: {
    ViewSelector,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    user: {
      type: Object,
      default: () => ({}),
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    onlyUserPrefs: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    groupsNavigationItems() {
      return Object.values(GROUPS_NAVIGATION_TYPES).map(type => ({
        text: this.$t(`user.navigationTypes.${type}`),
        value: type,
      }));
    },
  },
  methods: {
    showCopyAuthKeySuccessPopup() {
      this.$popups.success({ text: this.$t('success.authKeyCopied') });
    },

    showCopyAuthKeyErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>

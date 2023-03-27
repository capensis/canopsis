<template lang="pug">
  v-layout(column)
    c-id-field(v-field="form._id", :disabled="onlyUserPrefs || !isNew")
    c-name-field(
      v-field="form.name",
      :label="$t('common.username')",
      :disabled="onlyUserPrefs",
      browser-autocomplete="new-password",
      required
    )
    v-text-field(
      v-field="form.firstname",
      :label="$t('user.firstName')",
      :disabled="onlyUserPrefs"
    )
    v-text-field(
      v-field="form.lastname",
      :label="$t('user.lastName')",
      :disabled="onlyUserPrefs"
    )
    v-text-field(
      v-field="form.email",
      v-validate="'required|email'",
      :label="$t('user.email')",
      :disabled="onlyUserPrefs",
      :error-messages="errors.collect('email')",
      name="email",
      browser-autocomplete="new-password"
    )
    c-password-field(
      v-field="form.password",
      :required="isNew",
      browser-autocomplete="new-password"
    )
    c-role-field(v-field="form.role", :disabled="onlyUserPrefs", required)
    c-language-field(
      v-field="form.ui_language",
      :label="$t('user.language')"
    )
    v-select.mt-0(
      v-field="form.ui_groups_navigation_type",
      :label="$t('user.navigationType')",
      :items="groupsNavigationItems"
    )
    v-select(
      v-field="form.ui_theme",
      :label="$tc('common.theme')",
      :items="themes"
    )
    v-layout(v-if="!isNew", row, align-center)
      div {{ $t('common.authKey') }}: {{ user.authkey }}
      c-copy-btn(
        :value="user.authkey",
        :tooltip="$t('modals.variablesHelp.copyToClipboard')",
        small,
        fab,
        left,
        @success="showCopyAuthKeySuccessPopup",
        @error="showCopyAuthKeyErrorPopup"
      )
    c-enabled-field(
      v-field="form.enable",
      :disabled="onlyUserPrefs"
    )
    view-selector(v-field="form.defaultview")
</template>

<script>
import { THEMES_NAMES } from '@/config';

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

    themes() {
      return Object.values(THEMES_NAMES).map(name => ({
        text: this.$t(`common.themes.${name}`),
        value: name,
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

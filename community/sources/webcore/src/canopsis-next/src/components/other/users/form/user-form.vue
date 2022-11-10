<template lang="pug">
  v-layout(column)
    c-progress-overlay(:pending="pending")
    c-id-field(v-field="form._id", :disabled="onlyUserPrefs || !isNew")
    c-name-field(
      v-field="form.name",
      :label="$t('common.username')",
      :disabled="onlyUserPrefs",
      browser-autocomplete="new-password"
    )
    v-text-field(
      v-field="form.firstname",
      :label="$t('users.firstName')",
      :disabled="onlyUserPrefs"
    )
    v-text-field(
      v-field="form.lastname",
      :label="$t('users.lastName')",
      :disabled="onlyUserPrefs"
    )
    v-text-field(
      v-field="form.email",
      v-validate="'required|email'",
      :label="$t('users.email')",
      :disabled="onlyUserPrefs",
      :error-messages="errors.collect('email')",
      name="email",
      browser-autocomplete="new-password"
    )
    v-text-field(
      v-field="form.password",
      v-validate="passwordRules",
      :label="$t('common.password')",
      :error-messages="errors.collect('password')",
      type="password",
      name="password",
      browser-autocomplete="new-password"
    )
    v-select(
      v-field="form.role",
      v-validate="'required'",
      :label="$tc('common.role')",
      :items="roles",
      :disabled="onlyUserPrefs",
      :error-messages="errors.collect('role')",
      return-object,
      item-text="_id",
      item-value="_id",
      name="role"
    )
    c-language-field(
      v-field="form.ui_language",
      :label="$t('users.language')"
    )
    v-select.mt-0(
      v-field="form.ui_groups_navigation_type",
      :label="$t('users.navigationType')",
      :items="groupsNavigationItems"
    )
    v-select(
      v-field="form.ui_theme",
      :label="$t('common.theme')",
      :items="themes",
      @input="setTheme"
    )
    v-layout(row, align-center, v-if="!isNew")
      div {{ $t('common.authKey') }}: {{ user.authkey }}
      v-tooltip(left)
        v-btn(
          v-clipboard:copy="user.authkey",
          v-clipboard:success="addAuthKeyCopiedSuccessPopup",
          v-clipboard:error="addAuthKeyCopiedErrorPopup",
          slot="activator",
          small,
          fab,
          icon,
          depressed
        )
          v-icon content_copy
        span {{ $t('modals.variablesHelp.copyToClipboard') }}
    c-enabled-field(
      v-field="form.enable",
      :disabled="onlyUserPrefs"
    )
    view-selector(v-field="form.defaultview")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { GROUPS_NAVIGATION_TYPES, MAX_LIMIT } from '@/constants';

import ViewSelector from '@/components/forms/fields/view-selector.vue';

const { mapActions } = createNamespacedHelpers('role');

export default {
  inject: ['$validator', '$system'],
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
  data() {
    return {
      roles: [],
      pending: true,
    };
  },
  computed: {
    passwordRules() {
      return {
        required: this.isNew,
      };
    },

    groupsNavigationItems() {
      return Object.values(GROUPS_NAVIGATION_TYPES).map(type => ({
        text: this.$t(`users.navigationTypes.${type}`),
        value: type,
      }));
    },

    themes() {
      return this.$system.themes.map(name => ({
        text: this.$t(`common.themes.${name}`),
        value: name,
      }));
    },
  },
  async mounted() {
    const { data: roles } = await this.fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT } });

    this.roles = roles;
    this.pending = false;
  },
  methods: {
    ...mapActions({
      fetchRolesListWithoutStore: 'fetchListWithoutStore',
    }),

    setTheme(value) {
      this.$system.setTheme(value);
    },

    setDarkMode(value) {
      this.$system.setDarkMode(value);
    },

    addAuthKeyCopiedSuccessPopup() {
      this.$popups.success({ text: this.$t('success.authKeyCopied') });
    },

    addAuthKeyCopiedErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>

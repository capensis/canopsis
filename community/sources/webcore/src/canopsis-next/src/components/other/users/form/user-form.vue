<template lang="pug">
  div
    c-progress-overlay(:pending="pending")
    v-layout(row)
      v-text-field(
        v-field="form._id",
        :label="$t('common.id')",
        :disabled="onlyUserPrefs || !isNew"
      )
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('users.fields.username')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('name')",
        name="name",
        browser-autocomplete="new-password",
        data-test="username"
      )
    v-layout(row)
      v-text-field(
        v-field="form.firstname",
        :label="$t('users.fields.firstName')",
        :disabled="onlyUserPrefs",
        data-test="firstName"
      )
    v-layout(row)
      v-text-field(
        v-field="form.lastname",
        :label="$t('users.fields.lastName')",
        :disabled="onlyUserPrefs",
        data-test="lastName"
      )
    v-layout(row)
      v-text-field(
        v-field="form.email",
        v-validate="'required|email'",
        :label="$t('users.fields.email')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('email')",
        name="email",
        browser-autocomplete="new-password",
        data-test="email"
      )
    v-layout(row)
      v-text-field(
        v-field="form.password",
        v-validate="passwordRules",
        :label="$t('users.fields.password')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('password')",
        type="password",
        name="password",
        browser-autocomplete="new-password",
        data-test="password"
      )
    v-layout(data-test="roleLayout", row)
      v-select(
        v-field="form.role",
        v-validate="'required'",
        :label="$t('users.fields.role')",
        :items="roles",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('role')",
        return-object,
        item-text="_id",
        item-value="_id",
        name="role",
        data-test="role"
      )
    v-layout(data-test="languageLayout", row)
      v-select(
        v-field="form.ui_language",
        :label="$t('users.fields.language')",
        :items="languages",
        data-test="language"
      )
    v-layout(data-test="navigationTypeLayout", row)
      v-select.mt-0(
        v-field="form.ui_groups_navigation_type",
        :label="$t('parameters.groupsNavigationType.title')",
        :items="groupsNavigationItems",
        data-test="navigationType"
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
    v-layout(row)
      c-enabled-field(
        v-field="form.enable",
        :disabled="onlyUserPrefs"
      )
    v-layout
      view-selector(v-field="form.defaultview")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { GROUPS_NAVIGATION_TYPES, MAX_LIMIT } from '@/constants';

import ViewSelector from '@/components/forms/fields/view-selector.vue';

const { mapActions } = createNamespacedHelpers('role');

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
  data() {
    return {
      roles: [],
      pending: true,
    };
  },
  computed: {
    languages() {
      return Object.keys(this.$i18n.messages);
    },

    passwordRules() {
      if (this.isNew) {
        return 'required';
      }

      return null;
    },

    groupsNavigationItems() {
      return [
        {
          text: this.$t('parameters.groupsNavigationType.items.sideBar'),
          value: GROUPS_NAVIGATION_TYPES.sideBar,
        },
        {
          text: this.$t('parameters.groupsNavigationType.items.topBar'),
          value: GROUPS_NAVIGATION_TYPES.topBar,
        },
      ];
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

    addAuthKeyCopiedSuccessPopup() {
      this.$popups.success({ text: this.$t('success.authKeyCopied') });
    },

    addAuthKeyCopiedErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>

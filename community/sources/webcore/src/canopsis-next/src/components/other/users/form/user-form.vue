<template lang="pug">
  div
    c-progress-overlay(:pending="pending")
    v-layout(row)
      c-id-field(v-field="form._id", :disabled="onlyUserPrefs || !isNew")
    v-layout(row)
      c-name-field(
        v-field="form.name",
        :label="$t('common.username')",
        :disabled="onlyUserPrefs",
        browser-autocomplete="new-password"
      )
    v-layout(row)
      v-text-field(
        v-field="form.firstname",
        :label="$t('users.firstName')",
        :disabled="onlyUserPrefs",
        data-test="firstName"
      )
    v-layout(row)
      v-text-field(
        v-field="form.lastname",
        :label="$t('users.lastName')",
        :disabled="onlyUserPrefs",
        data-test="lastName"
      )
    v-layout(row)
      v-text-field(
        v-field="form.email",
        v-validate="'required|email'",
        :label="$t('users.email')",
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
        :label="$t('common.password')",
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
        :label="$tc('common.role')",
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
      c-language-field(
        v-field="form.ui_language",
        :label="$t('users.language')"
      )
    v-layout(data-test="navigationTypeLayout", row)
      v-select.mt-0(
        v-field="form.ui_groups_navigation_type",
        :label="$t('users.navigationType')",
        :items="groupsNavigationItems",
        data-test="navigationType"
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

    showCopyAuthKeySuccessPopup() {
      this.$popups.success({ text: this.$t('success.authKeyCopied') });
    },

    showCopyAuthKeyErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>

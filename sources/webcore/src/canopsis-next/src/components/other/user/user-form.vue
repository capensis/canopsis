<template lang="pug">
  v-form
    progress-overlay(:pending="pending")
    v-layout(row)
      v-text-field(
        :value="form._id",
        v-validate="'required'",
        :label="$t('modals.createUser.fields.username')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('username')",
        name="username",
        data-test="username",
        @input="updateField('_id', $event)"
      )
    v-layout(row)
      v-text-field(
        :value="form.firstname",
        :label="$t('modals.createUser.fields.firstName')",
        :disabled="onlyUserPrefs",
        data-test="firstName",
        @input="updateField('firstname', $event)"
      )
    v-layout(row)
      v-text-field(
        :value="form.lastname",
        :label="$t('modals.createUser.fields.lastName')",
        :disabled="onlyUserPrefs",
        data-test="lastName",
        @input="updateField('lastname', $event)"
      )
    v-layout(row)
      v-text-field(
        :value="form.mail",
        v-validate="'required|email'",
        :label="$t('modals.createUser.fields.email')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('email')",
        name="email",
        data-test="email",
        @input="updateField('mail', $event)"
      )
    v-layout(row)
      v-text-field(
        :value="form.password",
        v-validate="passwordRules",
        :label="$t('modals.createUser.fields.password')",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('password')",
        type="password",
        name="password",
        data-test="password",
        @input="updateField('password', $event)"
      )
    v-layout(data-test="roleLayout", row)
      v-select(
        :value="form.role",
        v-validate="'required'",
        :label="$t('modals.createUser.fields.role')",
        :items="roles",
        :disabled="onlyUserPrefs",
        :error-messages="errors.collect('role')",
        item-text="_id",
        item-value="_id",
        name="role",
        data-test="role",
        @input="updateField('role', $event)"
      )
    v-layout(data-test="languageLayout", row)
      v-select(
        :value="form.ui_language",
        :label="$t('modals.createUser.fields.language')",
        :items="languages",
        data-test="language",
        @input="updateField('ui_language', $event)"
      )
    v-layout(data-test="navigationTypeLayout", row)
      v-select.mt-0(
        :value="form.groupsNavigationType",
        :label="$t('parameters.groupsNavigationType.title')",
        :items="groupsNavigationItems",
        data-test="navigationType",
        @input="updateField('groupsNavigationType', $event)"
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
          v-icon file_copy
        span {{ $t('modals.variablesHelp.copyToClipboard') }}
    v-layout(row)
      v-switch(
        :input-value="form.enable",
        :label="$t('modals.createUser.fields.enabled')",
        :disabled="onlyUserPrefs",
        color="primary",
        data-test="enabled",
        @change="updateField('enable', $event)"
      )
    v-layout
      view-selector(:value="form.defaultview", @input="updateField('defaultview', $event)")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { GROUPS_NAVIGATION_TYPES } from '@/constants';

import formMixin from '@/mixins/form';
import popupMixin from '@/mixins/popup';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import ViewSelector from '@/components/forms/fields/view-selector.vue';

const { mapActions } = createNamespacedHelpers('role');

export default {
  inject: ['$validator'],
  components: {
    ProgressOverlay,
    ViewSelector,
  },
  mixins: [formMixin, popupMixin],
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
    const { data: roles } = await this.fetchRolesListWithoutStore({ params: { limit: 0 } });

    this.roles = roles;
    this.pending = false;
  },
  methods: {
    ...mapActions({
      fetchRolesListWithoutStore: 'fetchListWithoutStore',
    }),

    addAuthKeyCopiedSuccessPopup() {
      this.addSuccessPopup({ text: this.$t('success.authKeyCopied') });
    },

    addAuthKeyCopiedErrorPopup() {
      this.addErrorPopup({ text: this.$t('errors.default') });
    },
  },
};
</script>

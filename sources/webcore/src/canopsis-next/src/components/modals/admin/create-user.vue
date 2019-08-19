<template lang="pug">
  v-card(data-test="createUserModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ title }}
    v-card-text
      progress-overlay(:pending="pending")
      v-form
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.username')",
          v-model="form._id",
          data-vv-name="username",
          v-validate="'required'",
          :disabled="onlyUserPrefs",
          :error-messages="errors.collect('username')",
          data-test="username"
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.firstName')",
          v-model="form.firstname",
          :disabled="onlyUserPrefs",
          data-test="firstName"
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.lastName')",
          v-model="form.lastname",
          :disabled="onlyUserPrefs",
          data-test="lastName"
          )
        v-layout(row)
          v-text-field(
          :label="$t('modals.createUser.fields.email')",
          v-model="form.mail",
          data-vv-name="email",
          v-validate="'required|email'",
          :error-messages="errors.collect('email')",
          :disabled="onlyUserPrefs",
          data-test="email"
          )
        v-layout(row)
          v-text-field(
          type="password",
          :label="$t('modals.createUser.fields.password')",
          v-model="form.password",
          data-vv-name="password",
          v-validate="passwordRules",
          :error-messages="errors.collect('password')",
          :disabled="onlyUserPrefs",
          data-test="password"
          )
        v-layout(data-test="roleLayout", row)
          v-select(
          :label="$t('modals.createUser.fields.role')",
          v-model="form.role",
          :items="roles",
          item-text="_id",
          item-value="_id",
          data-vv-name="role",
          v-validate="'required'",
          :disabled="onlyUserPrefs",
          :error-messages="errors.collect('role')",
          data-test="role"
          )
        v-layout(data-test="languageLayout", row)
          v-select(
          data-test="language",
          :label="$t('modals.createUser.fields.language')",
          v-model="form.ui_language",
          :items="languages",
          )
        v-layout(data-test="navigationTypeLayout", row)
          v-select.mt-0(
          data-test="navigationType",
          v-model="form.groupsNavigationType",
          :label="$t('parameters.groupsNavigationType.title')",
          :items="groupsNavigationItems",
          )
        v-layout(row, align-center, v-if="!isNew")
          div {{ $t('common.authKey') }}: {{ config.user.authkey }}
          v-tooltip(left)
            v-btn(
            v-clipboard:copy="config.user.authkey",
            v-clipboard:success="() => addSuccessPopup({ text: $t('success.pathCopied') })",
            v-clipboard:error="() => addErrorPopup({ text: $t('errors.default') })",
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
          data-test="enabled"
          color="primary",
          :label="$t('modals.createUser.fields.enabled')",
          :disabled="onlyUserPrefs",
          v-model="form.enable",
          )
        v-layout
          view-selector(v-model="form.defaultview")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(data-test="submitButton", @click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { pick } from 'lodash';

import { MODALS, GROUPS_NAVIGATION_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesUserMixin from '@/mixins/entities/user';
import entitiesViewMixin from '@/mixins/entities/view/index';
import popupMixin from '@/mixins/popup';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import ViewSelector from './partial/view-selector.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createUser,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ProgressOverlay,
    ViewSelector,
  },
  mixins: [
    authMixin,
    modalInnerMixin,
    entitiesRoleMixin,
    entitiesUserMixin,
    entitiesViewMixin,
    popupMixin,
  ],
  data() {
    return {
      pending: true,
      form: {
        _id: '',
        firstname: '',
        lastname: '',
        mail: '',
        password: '',
        role: null,
        ui_language: 'fr',
        enable: true,
        defaultview: '',
        groupsNavigationType: GROUPS_NAVIGATION_TYPES.sideBar,
      },
    };
  },
  computed: {
    languages() {
      return Object.keys(this.$i18n.messages);
    },

    title() {
      return this.config.title || this.$t('modals.createUser.title');
    },

    passwordRules() {
      if (this.isNew) {
        return 'required';
      }

      return null;
    },
    isNew() {
      return !this.config.user;
    },

    onlyUserPrefs() {
      return this.config.onlyUserPrefs;
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
    await this.fetchRolesList({ params: { limit: 0 } });

    if (!this.isNew) {
      this.form = pick(this.config.user, [
        '_id',
        'firstname',
        'lastname',
        'mail',
        'password',
        'role',
        'ui_language',
        'enable',
        'defaultview',
        'groupsNavigationType',
      ]);

      if (!this.form.groupsNavigationType) {
        this.form.groupsNavigationType = GROUPS_NAVIGATION_TYPES.sideBar;
      }
    }

    this.pending = false;
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.hideModal();
      }
    },
  },
};
</script>

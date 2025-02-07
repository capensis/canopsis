<template>
  <v-layout column>
    <c-id-field
      v-field="form._id"
      :disabled="onlyUserPrefs || !isNew"
      autofocus
    />
    <c-name-field
      v-field="form.name"
      :label="$t('common.username')"
      :disabled="onlyUserPrefs || idpFieldsMap['name']"
      :autofocus="!isNew"
      autocomplete="new-password"
      required
    />
    <v-text-field
      v-field="form.firstname"
      :label="$t('user.firstName')"
      :disabled="onlyUserPrefs || idpFieldsMap['firstname']"
    />
    <v-text-field
      v-field="form.lastname"
      :label="$t('user.lastName')"
      :disabled="onlyUserPrefs || idpFieldsMap['lastname']"
    />
    <v-text-field
      v-field="form.email"
      v-validate="'required|email'"
      :label="$t('user.email')"
      :disabled="onlyUserPrefs || idpFieldsMap['email']"
      :error-messages="errors.collect('email')"
      name="email"
      autocomplete="new-password"
    />
    <c-password-field
      v-if="hasPassword"
      v-field="form.password"
      :required="isNew"
      :autofocus="onlyUserPrefs"
      autocomplete="new-password"
    />
    <c-role-field
      v-field="form.roles"
      :disabled="onlyUserPrefs || idpFieldsMap['roles']"
      :label="$tc('common.role', 2)"
      :is-disabled-items="isDisabledRoleItem"
      required
      multiple
      chips
    />
    <c-language-field
      v-field="form.ui_language"
      :label="$t('user.language')"
    />
    <v-select
      v-field="form.ui_groups_navigation_type"
      :label="$t('user.navigationType')"
      :items="groupsNavigationItems"
      class="mt-0"
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
import { computed } from 'vue';

import { AUTH_SOURCES_WITH_PASSWORD_CHANGING, GROUPS_NAVIGATION_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

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
  setup(props) {
    const { t } = useI18n();
    const popups = usePopups();

    const hasPassword = computed(() => (
      Object.values(AUTH_SOURCES_WITH_PASSWORD_CHANGING).includes(props.user?.source ?? '')
    ));

    const groupsNavigationItems = computed(() => Object.values(GROUPS_NAVIGATION_TYPES).map(type => ({
      text: t(`user.navigationTypes.${type}`),
      value: type,
    })));

    const idpFieldsMap = computed(() => (props.user?.idp_fields ?? []).reduce((acc, field) => {
      acc[field] = true;

      return acc;
    }, {}));

    const isDisabledRoleItem = item => (props.user?.idp_roles ?? []).includes(item._id);
    const showCopyAuthKeySuccessPopup = () => popups.success({ text: t('success.authKeyCopied') });
    const showCopyAuthKeyErrorPopup = () => popups.error({ text: t('errors.default') });

    return {
      hasPassword,
      groupsNavigationItems,
      idpFieldsMap,

      isDisabledRoleItem,
      showCopyAuthKeySuccessPopup,
      showCopyAuthKeyErrorPopup,
    };
  },
};
</script>

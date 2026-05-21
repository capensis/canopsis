<template>
  <v-layout class="gap-2" column>
    <c-enabled-field
      v-field="form.enabled"
      :disabled="onlyUserPrefs || isSelf"
      hide-details
      with-background
    />

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

    <c-form-block>
      <c-form-block-row :label="$t('user.email')">
        <v-text-field
          v-field="form.email"
          v-validate="'required|email'"
          :label="$t('user.email')"
          :disabled="onlyUserPrefs || idpFieldsMap['email']"
          :error-messages="errors.collect('email')"
          name="email"
          autocomplete="new-password"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.password')">
        <c-password-field
          v-field="form.password"
          :required="isNew"
          :autofocus="onlyUserPrefs"
          autocomplete="new-password"
          visibility
        />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.role', 2)">
        <c-role-field
          v-field="form.roles"
          :disabled="onlyUserPrefs || idpFieldsMap['roles']"
          :label="$tc('common.role', 2)"
          :is-disabled-items="isDisabledRoleItem"
          name="roles"
          required
          multiple
          chips
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('user.language')">
        <c-language-field
          v-field="form.ui_language"
          :label="$t('user.language')"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('user.navigationType')">
        <v-select
          v-field="form.ui_groups_navigation_type"
          :label="$t('user.navigationType')"
          :items="groupsNavigationItems"
          :menu-props="menuProps"
          class="mt-0"
        />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.theme', 2)">
        <c-theme-field v-if="hasReadThemeAccess" v-field="form.ui_theme" clearable />
      </c-form-block-row>

      <c-form-block-row
        v-if="!isNew"
        :label="$t('common.authKey')"
        align-center
      >
        <user-auth-key-field :value="user.authkey" />
      </c-form-block-row>

      <c-form-block-row :label="$t('role.defaultView')">
        <view-selector v-field="form.defaultview" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { AUTH_SOURCES_WITH_PASSWORD_CHANGING, GROUPS_NAVIGATION_TYPES, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useAuth, useCRUDPermissions } from '@/hooks/auth';

import ViewSelector from '@/components/forms/fields/view-selector.vue';
import UserAuthKeyField from '@/components/other/users/form/fields/user-auth-key-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ViewSelector,
    UserAuthKeyField,
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
    const menuProps = { offsetY: true };

    const { t } = useI18n();
    const { currentUser } = useAuth();

    const { hasReadAccess: hasReadThemeAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.profile.theme);

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

    const isSelf = computed(() => props.user._id === currentUser.value._id);

    const isDisabledRoleItem = item => (props.user?.idp_roles ?? []).includes(item._id);

    return {
      menuProps,

      hasReadThemeAccess,
      hasPassword,
      groupsNavigationItems,
      idpFieldsMap,
      isSelf,

      isDisabledRoleItem,
    };
  },
};
</script>

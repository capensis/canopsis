<template>
  <c-layout gap="2" column>
    <c-enabled-field
      v-model="enabledField"
      :disabled="onlyUserPrefs || isSelf"
      hide-details
      with-background
    />

    <c-id-field
      v-model="idField"
      :disabled="onlyUserPrefs || !isNew"
      name="_id"
      autofocus
    />

    <c-name-field
      v-model="nameField"
      :label="t('common.username')"
      :disabled="onlyUserPrefs || idpFieldsMap.name"
      :autofocus="!isNew"
      name="name"
      autocomplete="new-password"
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('user.email')">
        <c-text-field
          v-model="emailField"
          :label="t('user.email')"
          :disabled="onlyUserPrefs || idpFieldsMap.email"
          name="email"
          autocomplete="new-password"
          required
          email
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.password')">
        <c-password-field
          v-model="passwordField"
          :required="isNew"
          :autofocus="onlyUserPrefs"
          name="password"
          autocomplete="new-password"
          visibility
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.role', 2)">
        <c-role-field
          v-model="rolesField"
          :disabled="onlyUserPrefs || idpFieldsMap.roles"
          :label="t('common.role', 2)"
          :is-disabled-items="isDisabledRoleItem"
          name="roles"
          required
          multiple
          chips
        />
      </c-form-block-row>

      <c-form-block-row :label="t('user.language')">
        <c-language-field
          v-model="uiLanguageField"
          :label="t('user.language')"
          name="ui_language"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('user.navigationType')">
        <c-select
          v-model="uiGroupsNavigationTypeField"
          :label="t('user.navigationType')"
          :items="groupsNavigationItems"
          :menu-props="menuProps"
          item-title="text"
          item-value="value"
          name="ui_groups_navigation_type"
          class="mt-0"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.theme', 2)">
        <c-theme-field
          v-if="hasReadThemeAccess"
          v-model="uiThemeField"
          name="ui_theme"
          clearable
        />
      </c-form-block-row>

      <c-form-block-row v-if="!isNew" :label="t('common.authKey')" align-center>
        <user-auth-key-field :value="user.authkey" />
      </c-form-block-row>

      <c-form-block-row :label="t('role.defaultView')">
        <view-selector v-model="defaultviewField" name="defaultview" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { GROUPS_NAVIGATION_TYPES, USER_PERMISSIONS } from '@/constants';

import { useAuth, useCRUDPermissions } from '@/hooks/auth';
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ViewSelector from '@/components/forms/fields/view-selector.vue';
import UserAuthKeyField from '@/components/other/users/form/fields/user-auth-key-field.vue';

export default {
  components: {
    ViewSelector,
    UserAuthKeyField,
  },

  props: {
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
    const { hasReadAccess: hasReadThemeAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.profile.theme,
    );

    const {
      enabledField,
      idField,
      nameField,
      emailField,
      passwordField,
      rolesField,
      uiLanguageField,
      uiGroupsNavigationTypeField,
      uiThemeField,
      defaultviewField,
    } = useFormModels({
      enabledField: 'enabled',
      idField: '_id',
      nameField: 'name',
      emailField: 'email',
      passwordField: 'password',
      rolesField: 'roles',
      uiLanguageField: 'ui_language',
      uiGroupsNavigationTypeField: 'ui_groups_navigation_type',
      uiThemeField: 'ui_theme',
      defaultviewField: 'defaultview',
    });

    const groupsNavigationItems = computed(() => Object.values(GROUPS_NAVIGATION_TYPES).map(type => ({
        text: t(`user.navigationTypes.${type}`),
        value: type,
      })),
    );

    const idpFieldsMap = computed(() => (props.user?.idp_fields ?? []).reduce((acc, field) => {
        acc[field] = true;

        return acc;
      }, {}),
    );

    const isSelf = computed(() => props.user._id === currentUser.value?._id);

    const isDisabledRoleItem = item => (props.user?.idp_roles ?? []).includes(item._id);

    return {
      t,
      menuProps,
      enabledField,
      idField,
      nameField,
      emailField,
      passwordField,
      rolesField,
      uiLanguageField,
      uiGroupsNavigationTypeField,
      uiThemeField,
      defaultviewField,
      groupsNavigationItems,
      hasReadThemeAccess,
      idpFieldsMap,
      isSelf,
      isDisabledRoleItem,
    };
  },
};
</script>

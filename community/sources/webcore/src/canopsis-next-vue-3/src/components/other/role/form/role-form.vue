<template>
  <c-layout gap="2" column>
    <role-template-field v-if="withTemplate" v-model="permissionsField" />

    <c-name-field
      v-model="nameField"
      :disabled="!isNew"
      name="name"
      required
      autofocus
    />

    <c-form-block>
      <c-form-block-row :label="t('common.type')">
        <role-type-field v-model="typeField" :disabled="!isNew" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.description')" indented>
        <c-description-field
          v-model="descriptionField"
          name="description"
        />
      </c-form-block-row>

      <c-form-block-row v-if="isUiType" :label="t('common.theme', 2)">
        <c-theme-field v-model="uiThemeField" name="ui_theme" clearable />
      </c-form-block-row>

      <c-form-block-row :label="t('role.expirationSettings')" align-center>
        <c-enabled-field
          v-model="intervalsEnabledField"
          :label="t('role.expirationSettings')"
          hide-details
        />
      </c-form-block-row>

      <v-expand-transition>
        <c-layout v-if="intervalsEnabledField" column>
          <c-form-block-row v-if="isUiType" :label="t('role.inactivityInterval')">
            <c-duration-field
              v-model:duration="inactivityIntervalField"
              :label="t('role.inactivityInterval')"
              name="auth_config.inactivity_interval"
              long
            >
              <template #append>
                <c-help-icon :text="t('role.inactivityIntervalHelpText')" icon="help" top />
              </template>
            </c-duration-field>
          </c-form-block-row>

          <c-form-block-row :label="t('role.expirationInterval')" bottom-border>
            <c-duration-field
              v-model:duration="expirationIntervalField"
              :label="t('role.expirationInterval')"
              name="auth_config.expiration_interval"
              long
            >
              <template #append>
                <c-help-icon :text="t('role.expirationIntervalHelpText')" icon="help" top />
              </template>
            </c-duration-field>
          </c-form-block-row>
        </c-layout>
      </v-expand-transition>

      <c-form-block-row :label="t('role.defaultView')">
        <view-selector v-model="defaultviewField" name="defaultview" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { ROLE_TYPES } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import ViewSelector from '@/components/forms/fields/view-selector.vue';

import RoleTemplateField from './fields/role-template-field.vue';
import RoleTypeField from './fields/role-type-field.vue';

export default {
  components: {
    RoleTemplateField,
    RoleTypeField,
    ViewSelector,
  },

  props: {
    withTemplate: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();

    const {
      nameField,
      typeField,
      descriptionField,
      uiThemeField,
      permissionsField,
      intervalsEnabledField,
      inactivityIntervalField,
      expirationIntervalField,
      defaultviewField,
    } = useFormModels({
      nameField: 'name',
      typeField: 'type',
      descriptionField: 'description',
      uiThemeField: 'ui_theme',
      permissionsField: 'permissions',
      intervalsEnabledField: 'auth_config.intervals_enabled',
      inactivityIntervalField: 'auth_config.inactivity_interval',
      expirationIntervalField: 'auth_config.expiration_interval',
      defaultviewField: 'defaultview',
    });

    const isUiType = computed(() => typeField.value === ROLE_TYPES.ui);

    return {
      t,
      nameField,
      typeField,
      descriptionField,
      uiThemeField,
      permissionsField,
      intervalsEnabledField,
      inactivityIntervalField,
      expirationIntervalField,
      defaultviewField,
      isUiType,
    };
  },
};
</script>

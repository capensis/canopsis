<template>
  <v-layout class="gap-3" column>
    <role-template-field
      v-if="withTemplate"
      v-field="form.permissions"
    />
    <c-name-field
      v-field="form.name"
      :disabled="!isNew"
      required
      autofocus
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.type')">
        <role-type-field
          v-field="form.type"
          :disabled="!isNew"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-field="form.description" />
      </c-form-block-row>

      <c-form-block-row v-if="isUiType" :label="$tc('common.theme', 1)">
        <c-theme-field v-field="form.ui_theme" clearable />
      </c-form-block-row>

      <c-form-block-row :label="$t('role.expirationSettings')">
        <c-enabled-field v-field="form.auth_config.intervals_enabled" />
      </c-form-block-row>

      <v-expand-transition>
        <div v-if="form.auth_config.intervals_enabled">
          <c-form-block-row v-if="isUiType" :label="$t('role.inactivityInterval')">
            <c-duration-field v-field="form.auth_config.inactivity_interval" long>
              <template #append="">
                <c-help-icon
                  :text="$t('role.inactivityIntervalHelpText')"
                  icon="help"
                  top
                />
              </template>
            </c-duration-field>
          </c-form-block-row>

          <c-form-block-row :label="$t('role.expirationInterval')" bottom-border>
            <c-duration-field v-field="form.auth_config.expiration_interval" long>
              <template #append="">
                <c-help-icon
                  :text="$t('role.expirationIntervalHelpText')"
                  icon="help"
                  top
                />
              </template>
            </c-duration-field>
          </c-form-block-row>
        </div>
      </v-expand-transition>

      <c-form-block-row :label="$t('role.defaultView')">
        <view-selector v-field="form.defaultview" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { ROLE_TYPES } from '@/constants';

import ViewSelector from '@/components/forms/fields/view-selector.vue';

import RoleTemplateField from './fields/role-template-field.vue';
import RoleTypeField from './fields/role-type-field.vue';

export default {
  inject: ['$validator'],
  components: { ViewSelector, RoleTemplateField, RoleTypeField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const isUiType = computed(() => props.form.type === ROLE_TYPES.ui);

    return {
      isUiType,
    };
  },
};
</script>

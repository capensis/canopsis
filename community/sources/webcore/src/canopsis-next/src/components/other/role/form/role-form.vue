<template>
  <v-layout column>
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
    <role-type-field v-field="form.type" :disabled="!isNew" />
    <v-text-field
      v-field="form.description"
      :label="$t('common.description')"
    />
    <c-theme-field
      v-if="isUiType"
      v-field="form.ui_theme"
      clearable
    />
    <c-information-block :title="$t('role.expirationSettings')">
      <c-enabled-field
        v-field="form.auth_config.intervals_enabled"
        :label="form.auth_config.intervals_enabled ? $t('common.enabled') : $t('common.disabled')"
      />
      <v-expand-transition>
        <v-layout v-if="form.auth_config.intervals_enabled">
          <c-information-block
            :title="$t('role.inactivityInterval')"
            :help-text="$t('role.inactivityIntervalHelpText')"
          >
            <c-duration-field
              v-field="form.auth_config.inactivity_interval"
              long
            />
          </c-information-block>
          <c-information-block
            :title="$t('role.expirationInterval')"
            :help-text="$t('role.expirationIntervalHelpText')"
            class="ml-3"
          >
            <c-duration-field
              v-field="form.auth_config.expiration_interval"
              long
            />
          </c-information-block>
        </v-layout>
      </v-expand-transition>
    </c-information-block>
    <view-selector v-field="form.defaultview" />
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

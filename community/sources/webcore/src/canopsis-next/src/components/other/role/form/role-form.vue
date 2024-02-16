<template>
  <v-layout column>
    <role-template-field
      v-if="withTemplate"
      v-field="form.permissions"
    />
    <c-name-field
      v-field="form.name"
      required
    />
    <v-text-field
      v-field="form.description"
      :label="$t('common.description')"
    />
    <c-information-block :title="$t('role.expirationSettings')">
      <c-enabled-field
        v-field="form.auth_config.intervals_enabled"
        :label="form.auth_config.intervals_enabled ? $t('common.enabled') : $t('common.disabled')"
      />
      <v-expand-transition>
        <v-layout
          v-if="form.auth_config.intervals_enabled"
        >
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
            class="ml-3"
            :title="$t('role.expirationInterval')"
            :help-text="$t('role.expirationIntervalHelpText')"
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
import ViewSelector from '@/components/forms/fields/view-selector.vue';

import RoleTemplateField from './fields/role-template-field.vue';

export default {
  inject: ['$validator'],
  components: { ViewSelector, RoleTemplateField },
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
  },
};
</script>

<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.title"
      :label="$t('common.title')"
      name="title"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="$t('common.priority')">
        <c-priority-field
          v-field="form.priority"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('stateSetting.appliedForEntityType')">
        <c-entity-type-field
          v-field="form.type"
          :label="$t('stateSetting.appliedForEntityType')"
          :types="STATE_SETTING_ENTITY_TYPES"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('stateSetting.computeMethod')" indented>
        <state-setting-method-field v-field="form.method" />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import StateSettingMethodField from './fields/state-setting-method-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingMethodField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    return {
      STATE_SETTING_ENTITY_TYPES,
    };
  },
};
</script>

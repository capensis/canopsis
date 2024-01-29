<template>
  <v-layout column>
    <v-layout
      class="mb-2"
      align-center
    >
      <v-flex xs8>
        <c-name-field
          class="mr-2"
          v-field="form.title"
          :label="$t('common.title')"
          name="title"
          required
        />
      </v-flex>
      <v-flex xs4>
        <c-priority-field
          class="mx-2"
          v-field="form.priority"
          required
        />
      </v-flex>
    </v-layout>
    <v-layout align-center>
      <v-flex xs5>
        <c-enabled-field
          class="ml-2"
          v-field="form.enabled"
        />
      </v-flex>
      <v-flex xs7>
        <c-entity-type-field
          v-field="form.type"
          :label="$t('stateSetting.appliedForEntityType')"
          :types="availableEntityTypes"
          required
        />
      </v-flex>
    </v-layout>
    <state-setting-method-field v-field="form.method" />
  </v-layout>
</template>

<script>
import { STATE_SETTING_ENTITY_TYPES } from '@/constants';

import { formValidationHeaderMixin } from '@/mixins/form';

import StateSettingMethodField from '../fields/state-setting-method-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingMethodField },
  mixins: [formValidationHeaderMixin],
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
  computed: {
    availableEntityTypes() {
      return [...STATE_SETTING_ENTITY_TYPES];
    },
  },
};
</script>

<style lang="scss">
.state-setting-form {
  background-color: transparent !important;
}
</style>

<template>
  <v-layout column>
    <v-layout
      class="mb-2"
      align-center
    >
      <v-flex xs8>
        <c-name-field
          v-field="form.title"
          :label="$t('common.title')"
          class="mr-2"
          name="title"
          required
        />
      </v-flex>
      <v-flex xs4>
        <c-priority-field
          v-field="form.priority"
          class="mx-2"
          required
        />
      </v-flex>
    </v-layout>
    <v-layout align-center>
      <v-flex xs5>
        <c-enabled-field
          v-field="form.enabled"
          class="ml-2"
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

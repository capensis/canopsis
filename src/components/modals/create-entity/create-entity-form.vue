<template lang="pug">
v-card-text
  v-container
    v-layout(row)
      v-text-field(
      :label="$t('common.name')",
      v-model="form.name",
      @input="$emit('update:name', form.name)",
      v-validate="'required'",
      )
    v-layout(row)
      v-text-field(
      :label="$t('common.description')",
      v-model="form.description",
      @input="$emit('update:description', form.description)",
      v-validate="'required'",
      multi-line
      )
    v-layout(row)
      v-switch(:label="$t('common.enabled')", v-model="form.enabled", @change="$emit('update:enabled', form.enabled)")
      v-select(
      :items="types"
      v-model="form.type"
      @input="$emit('update:type', form.type)",
      label="Type"
      single-line
      )
</template>

<script>
import { MODALS } from '@/constants';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  data() {
    return {
      showValidationErrors: true,
      types: [
        this.$t('modals.createEntity.fields.types.connector'),
        this.$t('modals.createEntity.fields.types.component'),
        this.$t('modals.createEntity.fields.types.resource')],
      form: {
        name: '',
        description: '',
        type: '',
        enabled: true,
      },
    };
  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>

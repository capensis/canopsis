<template>
  <v-layout column>
    <v-text-field
      v-field="form.name"
      v-validate="nameRules"
      :label="$t('common.name')"
      :error-messages="errors.collect('name')"
      name
    />
    <v-text-field
      v-field="form.description"
      v-validate="descriptionRules"
      :label="$t('common.description')"
      :error-messages="errors.collect('description')"
      name="description"
    />
    <c-mixed-field
      v-field="form.value"
      :label="$t('common.value')"
      required
    />
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    entityInfo: {
      type: Object,
      default: () => ({}),
    },
    infos: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    infosNames() {
      return this.infos.map(({ name }) => name);
    },

    descriptionRules() {
      return {
        required: true,
      };
    },

    nameRules() {
      return {
        required: true,
        unique: {
          values: this.infosNames,
          initialValue: this.entityInfo?.name,
        },
      };
    },
  },
};
</script>

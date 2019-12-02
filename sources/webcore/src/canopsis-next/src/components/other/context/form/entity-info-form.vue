<template lang="pug">
  div
    v-text-field(
      v-field="form.name",
      v-validate="'required|unique-name'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    v-text-field(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('common.description')",
      :error-messages="errors.collect('description')",
      name="description"
    )
    v-textarea(
      v-field="form.value",
      v-validate="'required'",
      :label="$t('common.value')",
      :error-messages="errors.collect('value')",
      name="value"
    )
</template>

<script>
export default {
  inject: ['$validator'],
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
      type: Object,
      default: () => ({}),
    },
  },
  created() {
    this.createUniqueValidationRule();
  },
  methods: {
    createUniqueValidationRule() {
      this.$validator.extend('unique-name', {
        getMessage: () => this.$t('validator.unique'),
        validate: (value) => {
          if (this.entityInfo && this.entityInfo.name === value) {
            return true;
          }

          return !this.infos[value];
        },
      });
    },
  },
};
</script>

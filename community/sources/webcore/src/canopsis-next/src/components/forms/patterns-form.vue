<template lang="pug">
  v-layout(column)
    v-text-field(
      v-if="name",
      v-field="form.name",
      v-validate="nameRules",
      :label="$t('common.title')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    c-patterns-field(
      v-field="form",
      :entity="entity",
      :alarm="alarm",
      :pbehavior="pbehavior"
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
      required: true,
    },
    name: {
      type: Boolean,
      default: false,
    },
    alarm: {
      type: Boolean,
      default: false,
    },
    entity: {
      type: Boolean,
      default: false,
    },
    pbehavior: {
      type: Boolean,
      default: false,
    },
    totalEntity: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    nameRules() {
      return {
        required: this.name,
      };
    },
  },
};
</script>

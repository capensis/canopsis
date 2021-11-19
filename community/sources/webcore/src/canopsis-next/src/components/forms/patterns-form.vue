<template lang="pug">
  v-layout(column)
    v-text-field(
      v-if="name",
      v-field="form.name",
      v-validate="nameRules",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    c-patterns-field(
      v-field="form",
      :entity="entity",
      :alarm="alarm",
      :event="event",
      :total-entity="totalEntity"
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
    entity: {
      type: Boolean,
      default: false,
    },
    event: {
      type: Boolean,
      default: false,
    },
    alarm: {
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

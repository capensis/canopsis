<template lang="pug">
  div
    v-layout(row)
      c-pbehavior-type-field(
        v-field="form.active_on_pbh",
        :label="$t('remediation.pattern.tabs.pbehaviorTypes.fields.activeOnTypes')",
        :is-item-disabled="isActiveItemDisabled",
        chips,
        multiple
      )
    v-layout(row)
      c-pbehavior-type-field(
        v-field="form.disabled_on_pbh",
        :label="$t('remediation.pattern.tabs.pbehaviorTypes.fields.disabledOnTypes')",
        :is-item-disabled="isDisabledItemDisabled",
        chips,
        multiple
      )
</template>

<script>
export default {
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
  methods: {
    isActiveItemDisabled(item) {
      return this.form.disabled_on_pbh.includes(item._id);
    },

    isDisabledItemDisabled(item) {
      return this.form.active_on_pbh.includes(item._id);
    },
  },
};
</script>

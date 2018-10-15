<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.periodsNumber')}}
    v-container
      v-text-field(
      :value="value",
      @input="$emit('input', parseInt(event, 10))",
      v-validate="'required|integer|min:1'",
      data-vv-name="periodsNumber",
      :error-messages="errors.collect('periodsNumber')",
      )
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Number,
    },
  },
  methods: {
    updatePeriodsNumber(event) {
      this.$emit('input', parseInt(event, 10));
    },
  },
};
</script>


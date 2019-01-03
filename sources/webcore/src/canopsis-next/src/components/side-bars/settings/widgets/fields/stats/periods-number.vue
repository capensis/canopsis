<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.periodsNumber')}}
      .font-italic.caption.ml-1 ({{ $t('common.default') }}: 2)
    v-container
      v-text-field(
      :value="value",
      @input="updatePeriodsNumber",
      v-validate="'required|integer|min_value:2'",
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
      default: 2,
    },
  },
  methods: {
    updatePeriodsNumber(event) {
      if (event) {
        this.$emit('input', parseInt(event, 10));
      }
    },
  },
};
</script>


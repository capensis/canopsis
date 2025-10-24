<template>
  <v-tooltip :disabled="!disabled" top>
    <template #activator="{ on }">
      <span class="test-variables-tab__wrapper" v-on="on">
        <v-tab :disabled="disabled">
          {{ $tc('templateTesting.testVariables') }}
        </v-tab>
      </span>
    </template>
    <span>{{ $t('templateTesting.testVariablesDisabledTooltip') }}</span>
  </v-tooltip>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.test-variables-tab__wrapper {
  display: flex;
  flex: 1 1 auto;
  width: 100%;
  max-width: 360px;
}
</style>

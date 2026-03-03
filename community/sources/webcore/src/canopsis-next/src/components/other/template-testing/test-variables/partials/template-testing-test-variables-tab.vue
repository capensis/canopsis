<template>
  <v-tab :disabled="disabled" class="v-tab--tooltip">
    <v-tooltip :disabled="!disabled" top>
      <template #activator="{ on }">
        <span class="test-variables-tab__wrapper" v-on="on">
          {{ $tc('templateTesting.testVariables') }}
        </span>
      </template>
      <span>{{ $t('templateTesting.testVariablesDisabledTooltip') }}</span>
    </v-tooltip>
  </v-tab>
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

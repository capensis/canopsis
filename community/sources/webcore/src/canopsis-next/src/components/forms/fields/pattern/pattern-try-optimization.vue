<template>
  <c-alert type="info">
    <c-layout justify="space-between" align="center">
      <div>
        <span v-html="$t('pattern.tryOptimizationHint')" />
      </div>
      <v-btn color="info" type="button" @click="tryOptimization">
        <v-icon class="mr-2" color="white">
          wand_shine
        </v-icon>
        {{ $t('pattern.tryOptimization') }}
      </v-btn>
    </c-layout>
  </c-alert>
</template>

<script>
export default {
  setup(props, { emit }) {
    const tryOptimization = () => emit('try:optimization');

    return {
      tryOptimization,
    };
  },
};
</script>

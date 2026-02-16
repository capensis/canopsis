<template>
  <c-alert :value="true" type="info">
    <v-layout justify-space-between>
      <v-flex>
        <span v-html="$t('pattern.tryOptimizationHint')" />
      </v-flex>
      <v-btn color="info" @click="tryOptimization">
        <v-icon class="mr-2" color="white">
          $vuetify.icons.wand_shine
        </v-icon>
        {{ $t('pattern.tryOptimization') }}
      </v-btn>
    </v-layout>
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

<template>
  <div>
    <h6 class="my-2 text-h6 text-center">
      {{ $t('common.type') }}
    </h6>
    <v-layout justify-center>
      <v-radio-group
        class="my-0"
        v-field="impact"
        row
        hide-details
      >
        <v-radio
          :value="false"
          :label="$t('context.treeOfDependencies')"
          color="primary"
        />
        <v-radio
          :value="true"
          :label="$t('context.impactChain')"
          color="primary"
        />
      </v-radio-group>
    </v-layout>
  </div>
</template>

<script>
export default {
  model: {
    prop: 'impact',
    event: 'change',
  },
  props: {
    impact: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

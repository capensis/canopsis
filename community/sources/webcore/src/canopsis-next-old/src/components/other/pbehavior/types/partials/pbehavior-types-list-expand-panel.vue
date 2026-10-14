<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <v-tab>{{ $t('common.description') }}</v-tab>
    <v-tab-item>
      <v-layout
        class="pa-3 secondary lighten-2"
        column
      >
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <pre>{{ pbehaviorType.description }}</pre>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
    <v-tab>{{ $t('common.type') }}</v-tab>
    <v-tab-item>
      <v-layout
        class="pa-3 secondary lighten-2"
        column
      >
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <pre>{{ pbehaviorType.type }}</pre>
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
export default {
  props: {
    pbehaviorType: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

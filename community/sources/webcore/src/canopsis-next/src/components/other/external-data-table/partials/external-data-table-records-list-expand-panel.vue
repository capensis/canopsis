<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <v-layout column>
          <c-information-block-row
            v-for="(field, value) in record"
            :key="field"
            :label="field"
          >
            {{ value }}
          </c-information-block-row>
        </v-layout>
      </v-card-text>
    </v-card>
  </div>
</template>

<script>
export default {
  props: {
    record: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

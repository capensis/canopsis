<template>
  <div class="secondary pa-3">
    <v-card>
      <v-card-text>
        <v-layout column align-center>
          <v-flex>
            <c-information-block-row
              v-for="(value, field) in record"
              :key="field"
              :label="field"
            >
              {{ value }}
            </c-information-block-row>
          </v-flex>
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

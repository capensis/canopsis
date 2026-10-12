<template>
  <div>
    <template v-if="source">
      <c-layout v-for="field in fields" :key="field" align="center">
        <v-chip class="chip-small" color="info">
          <span>{{ field }}</span>
        </v-chip>
        <v-icon>arrow_right_alt</v-icon>
        <v-chip v-show="source[field]" class="chip-small" color="warning">
          <span>{{ source[field] }}</span>
        </v-chip>
      </c-layout>
    </template>
    <template v-else>
      <c-layout v-for="field in fields" :key="field" align="center">
        <v-chip class="chip-small" color="info">
          <span>{{ field }}</span>
        </v-chip>
      </c-layout>
    </template>
  </div>
</template>

<script>
export default {
  props: {
    fields: {
      type: Array,
      required: true,
    },
    source: {
      type: Object,
      default: null,
    },
  },
};
</script>

<style lang="scss" scoped>
.chip-small {
  height: 20px;

  :deep(.v-chip__content) {
    padding: 0 5px;
  }

  span {
    font-size: 0.75em;
    font-weight: 500;
  }
}
</style>

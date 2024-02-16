<template>
  <div>
    <template v-if="source">
      <v-layout
        v-for="field in fields"
        :key="field"
      >
        <v-chip
          class="chip-small va-middle"
          color="info"
          text-color="white"
        >
          <span>{{ field }}</span>
        </v-chip>
        <v-icon class="va-middle">
          arrow_right_alt
        </v-icon>
        <v-chip
          class="chip-small va-middle"
          v-show="source[field]"
          color="warning"
          text-color="white"
        >
          <span>{{ source[field] }}</span>
        </v-chip>
      </v-layout>
    </template>
    <template v-else>
      <v-layout
        v-for="field in fields"
        :key="field"
      >
        <v-chip
          class="chip-small va-middle"
          color="info"
          text-color="white"
        >
          <span>{{ field }}</span>
        </v-chip>
      </v-layout>
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
  .va-middle {
    vertical-align: middle;
  }

  .chip-small {
    height: 20px;

    & ::v-deep .v-chip__content {
      padding: 0 5px;
    }

    span {
      font-size: .75em;
      font-weight: 500;
    }
  }
</style>

<template>
  <v-breadcrumbs
    class="pa-0"
    :items="breadcrumbs"
  >
    <template #item="{ item }">
      <v-btn
        :disabled="item.last"
        :loading="item.last && pending"
        small
        text
        @click="$emit('click', item)"
      >
        <span class="text-none">{{ item.text }}</span>
      </v-btn>
    </template>
  </v-breadcrumbs>
</template>

<script>
export default {
  props: {
    previousMaps: {
      type: Array,
      default: () => [],
    },
    activeMap: {
      type: Object,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    breadcrumbs() {
      const breadcrumbs = this.previousMaps.map(({ name }, index) => ({
        index,
        text: name,
      }));

      if (this.activeMap || this.pending) {
        breadcrumbs.push({
          text: this.activeMap?.name,
          last: true,
        });
      }

      return breadcrumbs;
    },
  },
};
</script>

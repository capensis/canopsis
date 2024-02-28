<template>
  <v-list-item
    :to="link.route"
    :class="link.class"
    active-class=""
    @click="handleClick"
  >
    <v-list-item-avatar class="my-0">
      <v-icon
        :class="{ 'text--secondary': !link.color }"
        :color="link.color"
      >
        {{ link.icon }}
      </v-icon>
    </v-list-item-avatar>
    <v-list-item-title
      :class="{ [`${link.color}--text`]: link.color }"
      class="text-uppercase text-body-2"
    >
      {{ link.title }}
    </v-list-item-title>
  </v-list-item>
</template>

<script>
export default {
  props: {
    link: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    handleClick() {
      this.link.handler?.();
    },
  },
};
</script>

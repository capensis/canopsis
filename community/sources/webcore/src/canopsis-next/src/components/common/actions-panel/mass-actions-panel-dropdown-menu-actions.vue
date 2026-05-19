<template>
  <v-menu
    bottom
    left
    @click.native.stop=""
  >
    <template #activator="{ on }">
      <v-btn
        :small="small"
        icon
        v-on="on"
      >
        <v-icon>more_vert</v-icon>
      </v-btn>
    </template>
    <v-list>
      <v-list-item
        v-for="action in actions"
        :key="action.type"
        :disabled="action.disabled || action.loading"
        @click.stop="action.method"
      >
        <v-list-item-title>
          <span class="mr-4">
            <v-icon
              :color="action.iconColor"
              :disabled="action.disabled"
              class="ma-0 pa-0"
              left
              small
            >
              {{ action.icon }}
            </v-icon>
          </span>
          <span
            :class="action.cssClass"
            class="text-body-1"
          >
            {{ action.title }}
          </span>
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
export default {
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    small: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

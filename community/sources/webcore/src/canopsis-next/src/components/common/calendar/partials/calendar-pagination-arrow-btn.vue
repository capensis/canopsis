<template>
  <v-tooltip bottom>
    <template #activator="{ on }">
      <v-btn
        class="mx-2"
        v-on="on"
        fab
        small
        icon
        depressed
        @click="$emit('click', $event)"
      >
        <v-icon>{{ icon }}</v-icon>
      </v-btn>
    </template>
    <span>{{ label }}</span>
  </v-tooltip>
</template>

<script>
export default {
  props: {
    icon: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
  },
};
</script>

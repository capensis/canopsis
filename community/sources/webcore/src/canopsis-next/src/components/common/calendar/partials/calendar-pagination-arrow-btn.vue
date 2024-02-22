<template>
  <v-tooltip bottom>
    <template #activator="{ on }">
      <v-btn
        class="mx-2"
        fab
        small
        icon
        depressed
        v-on="on"
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

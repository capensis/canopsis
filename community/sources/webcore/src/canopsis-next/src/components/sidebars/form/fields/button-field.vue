<template>
  <v-container
    class="pa-3"
    fluid
  >
    <v-layout
      align-center
      justify-space-between
    >
      <slot name="title" />
      <v-layout justify-end>
        <v-btn
          class="primary"
          v-if="isEmpty && addable"
          small
          @click="$emit('create', $event)"
        >
          <span>{{ $t('common.create') }}</span>
        </v-btn>
        <template v-else>
          <v-btn
            class="primary"
            small
            @click="$emit('edit', $event)"
          >
            <span>{{ $t('common.edit') }}</span>
          </v-btn>
          <v-btn
            class="error"
            v-if="removable"
            small
            @click="$emit('delete', $event)"
          >
            <v-icon>delete</v-icon>
          </v-btn>
        </template>
      </v-layout>
    </v-layout>
  </v-container>
</template>

<script>
export default {
  props: {
    isEmpty: {
      type: Boolean,
      default: true,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

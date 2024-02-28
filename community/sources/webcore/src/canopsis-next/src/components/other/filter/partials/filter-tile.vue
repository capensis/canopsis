<template>
  <v-list-item class="pa-0">
    <v-layout align-center>
      <v-layout align-center>
        <v-icon
          v-if="editable"
          class="draggable ml-0 mr-3 action-drag-handler"
          small
        >
          drag_indicator
        </v-icon>
        <v-list-item-content>{{ filter.title }}</v-list-item-content>
      </v-layout>
      <v-list-item-action v-if="editable">
        <v-layout align-center>
          <c-action-btn
            type="edit"
            @click="$emit('edit')"
          />
          <c-action-btn
            type="delete"
            @click="$emit('delete')"
          />
        </v-layout>
      </v-list-item-action>
    </v-layout>
  </v-list-item>
</template>

<script>
export default {
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

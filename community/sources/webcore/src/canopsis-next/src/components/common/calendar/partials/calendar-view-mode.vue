<template>
  <v-menu>
    <template #activator="{ on }">
      <v-btn
        class="mx-2 my-1"
        v-on="on"
        text
      >
        {{ currentType }}
        <v-icon>arrow_drop_down</v-icon>
      </v-btn>
    </template>
    <v-list>
      <v-list-item
        v-for="type in types"
        :key="type.value"
        @click="$emit('input', type.value)"
      >
        <v-list-item-content>
          <v-list-item-title>{{ type.text }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      required: true,
    },
  },
  computed: {
    currentType() {
      return this.$t(`calendar.${this.value}`);
    },

    types() {
      return ['day', 'week', 'month'].map(type => ({
        text: this.$t(`calendar.${type}`),
        value: type,
      }));
    },
  },
};
</script>

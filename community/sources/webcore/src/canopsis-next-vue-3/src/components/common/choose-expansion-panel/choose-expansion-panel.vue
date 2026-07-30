<template>
  <div class="choose-expansion-panel">
    <v-expansion-panels class="my-1">
      <v-expansion-panel>
        <v-expansion-panel-title color="grey-darken-2">
          <div class="text-white">
            {{ label }}
          </div>

          <template #actions>
            <v-icon color="white" icon="keyboard_arrow_down">
keyboard_arrow_down
</v-icon>
          </template>
        </v-expansion-panel-title>
        <v-expansion-panel-text :class="{ error: errors.length }" class="grey-darken-2 text-white">
          <v-card class="pt-1">
            <v-alert v-if="errors.length" class="pa-2 mx-2" type="error">
              {{ errors.join(' ') }}
            </v-alert>
            <chips-list
              :entities="entities"
              :disabled-entities="disabledEntities"
              :existing-entities="existingEntities"
              :content-key="contentKey"
              :item-key="itemKey"
              :clearable="clearable"
              @remove="$emit('remove', $event)"
              @clear="$emit('clear', $event)"
            />
            <slot />
          </v-card>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import ChipsList from './partials/chips-list.vue';

export default {
  components: { ChipsList },
  emits: ['remove', 'clear'],
  props: {
    label: {
      type: String,
      required: true,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    contentKey: {
      type: String,
      required: false,
    },
    entities: {
      type: Array,
      default: () => [],
    },
    disabledEntities: {
      type: Array,
      default: () => [],
    },
    existingEntities: {
      type: Array,
      default: () => [],
    },
    errors: {
      type: Array,
      default: () => [],
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped lang="scss">
.choose-expansion-panel {
  & :deep(.v-expansion-panel-title .v-icon) {
    color: white !important;
  }
}
</style>

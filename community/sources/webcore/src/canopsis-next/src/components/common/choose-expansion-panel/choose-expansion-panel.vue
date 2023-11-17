<template>
  <div class="choose-expansion-panel">
    <v-expansion-panels class="my-1">
      <v-expansion-panel>
        <v-expansion-panel-header color="grey darken-2">
          <div class="white--text">
            {{ label }}
          </div>

          <template #actions>
            <v-icon color="white">
              keyboard_arrow_down
            </v-icon>
          </template>
        </v-expansion-panel-header>
        <v-expansion-panel-content
          class="grey darken-2 white--text"
          :class="{ error: errors.length }"
        >
          <v-card class="pt-1">
            <v-alert
              class="pa-2 mx-2"
              :value="!!errors.length"
              type="error"
            >
              {{ errors.join(' ') }}
            </v-alert>
            <chips-list
              :entities="entities"
              :disabled-entities="disabledEntities"
              :existing-entities="existingEntities"
              :content-key="contentKey"
              :item-key="itemKey"
              :clearable="clearable"
              @remove="$listeners.remove"
              @clear="$listeners.clear"
            />
            <slot />
          </v-card>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import ChipsList from './partials/chips-list.vue';

export default {
  components: { ChipsList },
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
    & ::v-deep .v-expansion-panel__header .v-icon {
      color: white !important;
    }
  }
</style>

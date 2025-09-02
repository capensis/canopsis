<template>
  <div>
    <v-layout align-center>
      <span>{{ message }}</span>
      <v-btn
        icon
        small
        @click="toggleDescriptionOpened"
      >
        <v-icon>help</v-icon>
      </v-btn>
    </v-layout>
    <v-expand-transition>
      <v-card v-show="opened" class="mt-2 pa-4 event-filter-enrichment-action-form-type-info-card">
        <v-layout column>
          <div
            v-html="description"
            class="pre-wrap"
          />
          <img
            v-if="image"
            :src="image"
            class="my-2"
            alt=""
            @click="showImageViewerModal"
          >
        </v-layout>
      </v-card>
    </v-expand-transition>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { eventFilterActionsTypesImages } from '@/assets';

export default {
  props: {
    type: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    message() {
      return this.$t(`eventFilter.actionsTypes.${this.type}.message`);
    },

    description() {
      return this.$t(`eventFilter.actionsTypes.${this.type}.description`);
    },

    image() {
      const imageName = `${this.$i18n.locale.toUpperCase()}_${this.type}`;

      return eventFilterActionsTypesImages[`./${imageName}.svg`] ?? '';
    },
  },
  methods: {
    toggleDescriptionOpened() {
      this.opened = !this.opened;
    },

    showImageViewerModal() {
      this.$modals.show({
        name: MODALS.imageViewer,
        config: {
          src: this.image,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
img {
  width: 100%;
  cursor: pointer;
}
</style>
<style lang="scss">
.event-filter-enrichment-action-form-type-info-card {
  background-color: var(--v-background-darken1, #dfdfdf) !important;
  border-radius: 10px !important;

  table {
    border-spacing: 0;

    th {
      background-color: var(--v-background-darken2, #b1b1b1);
    }

    td {
      background-color: var(--v-background-base, #FFFFFF);
    }

    td, th {
      border: 1px solid var(--v-background-darken3, #b1b1b1);
    }
  }
}
</style>

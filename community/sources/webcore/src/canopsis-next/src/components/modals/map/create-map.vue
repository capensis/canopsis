<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createMap.title') }}</span>
      </template>
      <template #text="">
        <v-layout column>
          <v-flex
            v-for="type in availableTypes"
            :key="type"
            class="my-1 cursor-pointer"
            @click="selectType(type)"
          >
            <v-card>
              <v-card-title
                class="py-3"
                primary-title
              >
                <v-layout>
                  <div class="text-subtitle-1">
                    {{ $t(`map.types.${type}`) }}
                  </div>
                  <v-spacer />
                  <v-icon class="text--secondary">
                    {{ getIconByType(type) }}
                  </v-icon>
                </v-layout>
              </v-card-title>
            </v-card>
          </v-flex>
        </v-layout>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, MAP_TYPES, MAP_ICON_BY_TYPES, CREATE_MAP_MODAL_NAMES_BY_TYPE } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMap,
  components: {
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
  ],
  computed: {
    availableTypes() {
      return Object.values(MAP_TYPES);
    },
  },
  methods: {
    getIconByType(type) {
      return MAP_ICON_BY_TYPES[type];
    },

    selectType(type) {
      this.$modals.show({
        name: CREATE_MAP_MODAL_NAMES_BY_TYPE[type],
        config: {
          map: { type },
          action: async (map) => {
            await this.config.action(map);

            this.$modals.hide();
          },
        },
      });
    },
  },
};
</script>

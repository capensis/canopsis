<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createMap.title') }}
      template(#text="")
        v-layout(column)
          v-flex.my-1.cursor-pointer(
            v-for="type in availableTypes",
            :key="type",
            @click="selectType(type)"
          )
            v-card
              v-card-title.py-3(primary-title)
                v-layout
                  div.subheading {{ $t(`map.types.${type}`) }}
                  v-spacer
                  v-icon {{ getIconByType(type) }}
</template>

<script>
import { MODALS, MAP_TYPES, MAP_ICON_BY_TYPES } from '@/constants';

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
      const modalName = {
        [MAP_TYPES.geo]: MODALS.createGeoMap,
        [MAP_TYPES.flowchart]: MODALS.createFlowchartMap,
        [MAP_TYPES.mermaid]: MODALS.createMermaidMap,
        [MAP_TYPES.treeOfDependencies]: MODALS.createTreeOfDependenciesMap,
      }[type];

      this.$modals.show({
        name: modalName,
        config: {
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

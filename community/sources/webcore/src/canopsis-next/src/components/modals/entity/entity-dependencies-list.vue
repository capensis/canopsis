<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      entity-dependencies-list-component(
        :widget="widget",
        :entity-id="config.entityId",
        :impact="config.impact"
      )
</template>

<script>
import { MODALS } from '@/constants';

import { generateDefaultContextWidget } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesListComponent from '@/components/other/entity/entity-dependencies-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entityDependenciesList,
  components: {
    ModalWrapper,
    EntityDependenciesListComponent,
  },
  mixins: [modalInnerMixin],
  computed: {
    title() {
      return this.config.title || this.$t('modals.createEntity.create.title');
    },

    widget() {
      return generateDefaultContextWidget();
    },
  },
};
</script>

<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <entity-upstream-network-graph
        :entity="config.entity"
        :pending="pending"
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import EntityUpstreamNetworkGraph from '@/components/other/entity/entity-upstream-network-graph.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entityUpstream,
  components: {
    ModalWrapper,
    EntityUpstreamNetworkGraph,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config } = useInnerModal(props);

    const pending = ref(false);

    const entity = computed(() => config.value.entity ?? {});
    const hasUpstream = computed(() => entity.value.upstream !== undefined);
    const title = computed(() => (hasUpstream.value
      ? t('modals.entityUpstream.topLevelEntities')
      : t('modals.entityUpstream.entities')));

    return {
      pending,
      config,
      entity,
      hasUpstream,
      title,
    };
  },
};
</script>
